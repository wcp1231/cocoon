package server

import (
	"cocoon/pkg/model/rpc"
	"cocoon/pkg/model/traffic"
	"cocoon/pkg/tcp"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
)

type RpcServer struct {
	server *server.Server
	logger *zap.Logger

	//dissectorService *dissector.DissectService
	tcpManager *tcp.TcpManager
}

func NewRpcServer(logger *zap.Logger) *RpcServer {
	rpcxServer := server.NewServer()
	resultC := make(chan *traffic.StreamItem, 1024)
	//dissectService := dissector.NewDissectService(logger)
	tcpManager := tcp.NewTcpManager(logger, resultC)

	rpcServer := &RpcServer{
		server:     rpcxServer,
		logger:     logger,
		tcpManager: tcpManager,
		//dissectorService: dissectService,
	}

	handler := NewCocoonHandler(logger, rpcServer)
	rpcxServer.RegisterName(rpc.COCOON_SERVER_NAME, handler, "")
	return rpcServer
}

func (r *RpcServer) Start(listen string) error {
	//r.dissectorService.Start()
	return r.server.Serve("tcp", listen)
}
