package rpc_server

import (
	"cocoon/pkg/db"
	"cocoon/pkg/model/rpc"
	"context"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
)

type RpcServer struct {
	server *server.Server
	logger *zap.Logger

	database *db.Database
	//dissectManager *dissector.DissectManager
	//tcpManager     *tcp.TcpManager
}

func NewRpcServer(logger *zap.Logger, dbUri string) *RpcServer {
	rpcxServer := server.NewServer()
	//resultC := make(chan *traffic.StreamItem, 1024)
	ctx := context.Background()
	database := db.NewDatabase(ctx, logger, dbUri)
	//dissectManager := dissector.NewDissectManager(ctx, logger)
	//tcpManager := tcp.NewTcpManager(logger, resultC)

	rpcServer := &RpcServer{
		server:   rpcxServer,
		logger:   logger,
		database: database,
		//tcpManager:     tcpManager,
		//dissectManager: dissectManager,
	}

	handler := NewCocoonHandler(logger, rpcServer)
	rpcxServer.RegisterName(rpc.COCOON_SERVER_NAME, handler, "")
	return rpcServer
}

func (r *RpcServer) Start(listen string) error {
	return r.server.Serve("tcp", listen)
}
