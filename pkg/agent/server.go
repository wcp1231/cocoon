package agent

import (
	"cocoon/pkg/mock"
	"cocoon/pkg/proto"
	"context"
	"go.uber.org/zap"
	"net"
)

type AgentProxy interface {
	ProxyStart() error
	Shutdown()
	GracefulShutdown()
	CloseWait()
}

type Server struct {
	ctx        context.Context
	logger     *zap.Logger
	appname    string
	session    string
	proxy      AgentProxy
	mockServer *mock.MockService
}

func NewServer(ctx context.Context, logger *zap.Logger, appname, session string) *Server {
	return &Server{
		ctx:     ctx,
		logger:  logger,
		appname: appname,
		session: session,
	}
}

func (s *Server) Init(listenAddr *net.TCPAddr, transparent bool, protocols string) error {
	if transparent {
		s.proxy = NewTransparentProxy(listenAddr, s)
	} else {
		s.proxy = NewSocks5Proxy(listenAddr, s)
	}

	proto.InitPresetClassifier(protocols)

	s.mockServer = mock.NewMockService(s.logger)
	return s.mockServer.InitFromFile()
}

func (s *Server) Start() error {
	return s.proxy.ProxyStart()
}

func (s *Server) HandleConn(ctx context.Context, inboundConn, outboundConn *net.TCPConn) {
	p := NewConnHandler(s, ctx, inboundConn, outboundConn)
	p.Start()
}

// Shutdown agent.
func (s *Server) Shutdown() {
	s.proxy.Shutdown()
	s.proxy.CloseWait()
}

// GracefulShutdown agent.
func (s *Server) GracefulShutdown() {
	s.proxy.GracefulShutdown()
	s.proxy.CloseWait()
}
