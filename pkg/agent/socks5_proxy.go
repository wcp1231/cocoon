package agent

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"strings"
	"sync"
)

type Socks5Proxy struct {
	logger     *zap.Logger
	listenAddr *net.TCPAddr
	server     *Agent
	listener   *net.TCPListener

	ctx        context.Context
	shutdown   context.CancelFunc
	Wg         *sync.WaitGroup
	ClosedChan chan struct{}
}

func NewSocks5Proxy(listenAddr *net.TCPAddr, server *Agent) *Socks5Proxy {
	innerCtx, shutdown := context.WithCancel(server.ctx)
	wg := &sync.WaitGroup{}
	closedChan := make(chan struct{})
	return &Socks5Proxy{
		listenAddr: listenAddr,
		logger:     server.logger,
		server:     server,

		ctx:        innerCtx,
		shutdown:   shutdown,
		Wg:         wg,
		ClosedChan: closedChan,
	}
}

func (s *Socks5Proxy) ProxyStart() error {
	s.logger.Info("Start socks5 proxy mode.",
		zap.String("listen", s.listenAddr.String()))
	lt, err := net.ListenTCP("tcp", s.listenAddr)
	if err != nil {
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

		s.Wg.Add(1)
		go s.process(conn)
	}
}

func (s *Socks5Proxy) process(conn *net.TCPConn) {
	defer s.Wg.Done()
	defer fmt.Println("conn process exit")

	if err := s.auth(conn); err != nil {
		fmt.Printf("socks5 auth error. %v\n", err)
		conn.Close()
		return
	}

	outboundConn, err := s.connectDst(conn)
	if err != nil {
		fmt.Printf("socks5 connect error. %v\n", err)
		conn.Close()
		return
	}

	s.server.HandleConn(s.ctx, conn, outboundConn)
}

func (s *Socks5Proxy) auth(conn net.Conn) error {
	buf := make([]byte, 256)
	n, err := io.ReadFull(conn, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}
	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}

	n, err = io.ReadFull(conn, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}

	// 无认证
	n, err = conn.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write resp err: " + err.Error())
	}
	return nil
}

func (s *Socks5Proxy) connectDst(conn net.Conn) (*net.TCPConn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(conn, buf[:4])
	if n != 4 {
		return nil, errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, fmt.Errorf("invalid ver/cmd. ver=%v, cmd=%v", ver, cmd)
	}

	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(conn, buf[:4])
		if n != 4 {
			return nil, errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 3:
		n, err = io.ReadFull(conn, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])
		n, err := io.ReadFull(conn, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])
	case 4:
		return nil, errors.New("IPv6 not supported")
	default:
		return nil, errors.New("invalid atyp")
	}

	n, err = io.ReadFull(conn, buf[:2])
	if n != 2 {
		return nil, errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])
	dstAddr := fmt.Sprintf("%s:%d", addr, port)
	dstTcpAddr, err := net.ResolveTCPAddr("tcp", dstAddr)
	if err != nil {
		return nil, err
	}
	dstConn, err := net.DialTCP("tcp", nil, dstTcpAddr)
	if err != nil {
		return nil, errors.New("dial dst: " + err.Error())
	}

	n, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dstConn.Close()
		return nil, errors.New("write resp: " + err.Error())
	}
	return dstConn, nil
}

// Shutdown agent.
func (s *Socks5Proxy) Shutdown() {
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
func (s *Socks5Proxy) GracefulShutdown() {
	select {
	case <-s.ctx.Done():
	default:
		if err := s.listener.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			s.logger.WithOptions(zap.AddCaller()).Error("agent listener Close error", zap.Error(err))
		}
	}
}

func (s *Socks5Proxy) CloseWait() {
	s.Wg.Wait()
	<-s.ClosedChan
}
