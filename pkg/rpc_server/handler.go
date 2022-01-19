package rpc_server

import (
	"cocoon/pkg/db"
	"cocoon/pkg/model/common"
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

func (c *CocoonHandler) ClientPostStart(ctx context.Context, args *rpc.PostStartReq, resp *rpc.PostStartResp) error {
	c.logger.Debug("Receive client post start",
		zap.String("app", args.Appname),
		zap.String("session", args.Session))
	c.rpcService.database.EnsureApplicationAndSession(args.Appname, args.Session)
	return nil
}

func (c *CocoonHandler) Upload(ctx context.Context, args *rpc.UploadReq, resp *rpc.UploadResp) error {
	c.logger.Debug("Receive upload",
		zap.String("session", args.Session),
		zap.String("src", args.Packet.Source),
		zap.String("direction", args.Packet.Direction.String()),
		zap.String("dest", args.Packet.Destination),
		zap.Uint64("seq", args.Packet.Seq),
		zap.Int("size", len(args.Packet.Payload)))
	dbModel := c.generatePacketDBModel(args.Session, args.Packet)
	c.rpcService.database.AppendTcpPacket(dbModel)
	return nil
}

func (c *CocoonHandler) Analysis(ctx context.Context, args *rpc.AnalysisReq, resp *rpc.AnalysisResp) error {
	c.logger.Debug("Handle analysis request", zap.String("session", args.Session))
	err := c.rpcService.dissectManager.Dissect(args.Session, c.rpcService.database)
	resp.Error = err
	return nil
}

func (c *CocoonHandler) generatePacketDBModel(session string, packet *common.TcpPacket) *db.TcpTraffic {
	return &db.TcpTraffic{
		Session:     session,
		Source:      packet.Source,
		Destination: packet.Destination,
		IsOutgoing:  packet.IsOutgoing,
		Direction:   packet.Direction,
		Seq:         packet.Seq,
		Timestamp:   packet.Timestamp,
		Size:        len(packet.Payload),
		Raw:         packet.Payload,
	}
}
