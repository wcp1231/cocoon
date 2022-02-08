package agent

import (
	"cocoon/pkg/mock"
	"cocoon/pkg/proto"
	"cocoon/pkg/record"
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/fs"
	"net"
	"net/http"
	"sync/atomic"
)

type AgentProxy interface {
	ProxyStart() error
	Shutdown()
	GracefulShutdown()
	CloseWait()
}

type Agent struct {
	ctx          context.Context
	logger       *zap.Logger
	appname      string
	session      string
	proxy        AgentProxy
	httpServer   *http.Server
	mockServer   *mock.MockService
	recordServer *record.RecordService

	id int32
}

func NewAgent(ctx context.Context, logger *zap.Logger, appname, session string) *Agent {
	return &Agent{
		ctx:     ctx,
		logger:  logger,
		appname: appname,
		session: session,
	}
}

func (s *Agent) Init(proxyListen, httpListen string, transparent bool, protocols string, statics fs.FS) error {
	err := s.initProxy(proxyListen, transparent)
	if err != nil {
		return err
	}

	proto.InitPresetClassifier(protocols)

	s.mockServer = mock.NewMockService(s.logger)
	s.recordServer = record.NewRecordService(s.logger)
	s.initHttp(httpListen, statics)

	return s.mockServer.InitFromFile()
}

func (s *Agent) initProxy(listen string, transparent bool) error {
	listenAddr, err := net.ResolveTCPAddr("tcp", listen)
	if err != nil {
		s.logger.Fatal("error", zap.Error(err))
		return err
	}

	if transparent {
		s.proxy = NewTransparentProxy(listenAddr, s)
	} else {
		s.proxy = NewSocks5Proxy(listenAddr, s)
	}
	return nil
}

func (s *Agent) initHttp(listen string, statics fs.FS) {
	router := mux.NewRouter()

	router.HandleFunc("/api/mocks/", s.mockServer.ListMocks).Methods("GET")
	router.HandleFunc("/api/mocks/", s.mockServer.AddMocks).Methods("POST")
	router.HandleFunc("/api/mocks/{id}", s.mockServer.EditMocks).Methods("POST")
	router.HandleFunc("/api/mocks/{id}", s.mockServer.DeleteMocks).Methods("DELETE")

	router.HandleFunc("/api/ws", s.recordServer.ServeWs)

	fileServer := http.FileServer(http.FS(statics))
	router.PathPrefix("/").Handler(fileServer)

	s.httpServer = &http.Server{
		Addr:    listen,
		Handler: router,
	}
}

func (s *Agent) Start() {
	go s.recordServer.Start()
	go func() {
		s.logger.Info("Listen http at",
			zap.String("listen", s.httpServer.Addr))
		err := s.httpServer.ListenAndServe()
		if err != nil {
			s.logger.Fatal("Listen http failed", zap.Error(err))
		}
	}()
	go func() {
		err := s.proxy.ProxyStart()
		if err != nil {
			s.logger.Fatal("Listen proxy failed", zap.Error(err))
		}
	}()
}

func (s *Agent) nextId() int32 {
	return atomic.AddInt32(&s.id, 1)
}

func (s *Agent) HandleConn(ctx context.Context, inboundConn, outboundConn *net.TCPConn) {
	p := NewConnHandler(s, ctx, inboundConn, outboundConn)
	p.Start()
}

// Shutdown agent.
func (s *Agent) Shutdown() {
	s.proxy.Shutdown()
	s.proxy.CloseWait()
}

// GracefulShutdown agent.
func (s *Agent) GracefulShutdown() {
	s.proxy.GracefulShutdown()
	s.proxy.CloseWait()
}
