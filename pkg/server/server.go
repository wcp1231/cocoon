package server

import (
	"cocoon/pkg/db"
	"cocoon/pkg/dissector"
	"cocoon/pkg/model/rpc"
	"context"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
)

type CocoonServer struct {
	server *server.Server
	logger *zap.Logger

	database       *db.Database
	dissectManager *dissector.DissectManager
	httpServer *CocoonHttpHandler
	//tcpManager     *tcp.TcpManager
}

func NewCocoonServer(logger *zap.Logger, dbUri string) *CocoonServer {
	rpcxServer := server.NewServer()
	ctx := context.Background()
	database := db.NewDatabase(ctx, logger, dbUri)
	dissectManager := dissector.NewDissectManager(ctx, logger)
	//tcpManager := tcp.NewTcpManager(logger, resultC)


	server := &CocoonServer{
		server:   rpcxServer,
		logger:   logger,
		database: database,
		//tcpManager:     tcpManager,
		dissectManager: dissectManager,
	}

	server.httpServer = NewCocoonHttpHandler(logger, server)
	handler := NewCocoonHandler(logger, server)
	rpcxServer.RegisterName(rpc.COCOON_SERVER_NAME, handler, "")
	return server
}

func (r *CocoonServer) Start(listen string) error {
	return r.server.Serve("tcp", listen)
}

func (r *CocoonServer) StartHttp(listen string) error {
	r.httpServer.srv.Addr = listen
	return r.httpServer.srv.ListenAndServe()
}