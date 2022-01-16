package server

import (
	"cocoon/pkg/model/rpc"
	"context"
	"go.uber.org/zap"
)

type CocoonHandler struct {
	logger     *zap.Logger
	rpcService *RpcServer
}

func NewCocoonHandler(logger *zap.Logger, rpcService *RpcServer) *CocoonHandler {
	return &CocoonHandler{
		logger:     logger,
		rpcService: rpcService,
	}
}

func (c *CocoonHandler) Upload(ctx context.Context, args *rpc.UploadReq, resp *rpc.UploadResp) error {
	c.logger.Debug("Receive upload",
		zap.String("app", args.Appname),
		zap.String("src", args.Packet.Source),
		zap.String("direction", args.Packet.Direction.String()),
		zap.String("dest", args.Packet.Destination),
		zap.Uint64("seq", args.Packet.Seq),
		zap.Int("size", len(args.Packet.Payload)))
	c.rpcService.tcpManager.Accept(args.Packet)
	return nil
}
