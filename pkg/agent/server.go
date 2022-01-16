package agent

import (
	"cocoon/pkg/client"
	"context"
	"github.com/cybozu-go/transocks"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"net"
	"strings"
	"sync"
)

type Server struct {
	logger     *zap.Logger
	appname    string
	listenAddr *net.TCPAddr
	rpcClient  *client.RpcClient
	ctx        context.Context
	shutdown   context.CancelFunc

	Wg         *sync.WaitGroup
	ClosedChan chan struct{}
	listener   *net.TCPListener
}

func NewServer(ctx context.Context, logger *zap.Logger, appname, remote string, listenAddr *net.TCPAddr) *Server {
	innerCtx, shutdown := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	closedChan := make(chan struct{})

	rpcClient := client.NewRpcClient(remote)

	return &Server{
		logger:     logger,
		appname:    appname,
		listenAddr: listenAddr,
		rpcClient:  rpcClient,
		ctx:        innerCtx,
		shutdown:   shutdown,

		Wg:         wg,
		ClosedChan: closedChan,
	}
}

func (s *Server) Start() error {
	lt, err := net.ListenTCP("tcp", s.listenAddr)
	if err != nil {
		s.logger.WithOptions(zap.AddCaller()).Fatal("listenAddr ListenTCP error", zap.Error(err))
		return err
	}
	s.listener = lt
	defer func() {
		if err := s.listener.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			s.logger.WithOptions(zap.AddCaller()).Error("agent listener Close error", zap.Error(err))
		}
		close(s.ClosedChan)
	}()

	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Temporary() {
					continue
				}
				if !strings.Contains(err.Error(), "use of closed network connection") {
					select {
					case <-s.ctx.Done():
						break
					default:
						s.logger.WithOptions(zap.AddCaller()).Fatal("listener Accept TCP error", zap.Error(err))
					}
				}
			}
			return err
		}
		originDst, err := transocks.GetOriginalDST(conn)
		if err != nil {
			s.logger.WithOptions(zap.AddCaller()).Fatal("listener Get Origin DST error", zap.Error(err))
			if err := conn.Close(); err != nil {
				s.logger.WithOptions(zap.AddCaller()).Error("agent conn Close error")
			}
			continue
		}
		s.Wg.Add(1)
		go s.handleConn(conn, originDst)
	}
}

// Shutdown agent.
func (s *Server) Shutdown() {
	select {
	case <-s.ctx.Done():
	default:
		s.shutdown()
		if err := s.listener.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			s.logger.WithOptions(zap.AddCaller()).Error("agent listener Close error", zap.Error(err))
		}
	}
}

// GracefulShutdown agent.
func (s *Server) GracefulShutdown() {
	select {
	case <-s.ctx.Done():
	default:
		if err := s.listener.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			s.logger.WithOptions(zap.AddCaller()).Error("agent listener Close error", zap.Error(err))
		}
	}
}

func (s *Server) handleConn(inboundConn *net.TCPConn, originDst *net.TCPAddr) {
	defer s.Wg.Done()

	outboundConn, err := net.DialTCP("tcp", nil, originDst)
	if err != nil {
		fields := s.fieldsWithErrorAndConn(err, inboundConn, originDst)
		s.logger.WithOptions(zap.AddCaller()).Error("remoteAddr DialTCP error", fields...)
		if err := inboundConn.Close(); err != nil {
			s.logger.WithOptions(zap.AddCaller()).Error("agent conn Close error", fields...)
		}
		return
	}

	p := NewDumpHandler(s, inboundConn, outboundConn)
	p.Start()
}

func (s *Server) fieldsWithErrorAndConn(err error, conn *net.TCPConn, originDst *net.TCPAddr) []zapcore.Field {
	fields := []zapcore.Field{
		zap.Error(err),
		zap.String("client_addr", conn.RemoteAddr().String()),
		zap.String("proxy_listen_addr", conn.LocalAddr().String()),
		zap.String("remote_addr", originDst.String()),
	}
	return fields
}
