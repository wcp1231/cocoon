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
		zap.String("proto", args.Packet.Protocol),
		zap.Uint64("seq", args.Packet.Seq),
		zap.Int("size", len(args.Packet.Payload)))
	dbModel := c.generatePacketDBModel(args.Session, args.Packet)
	c.rpcService.database.AppendTcpPacket(dbModel)
	return nil
}

func (c *CocoonHandler) ConnClose(ctx context.Context, args *rpc.ConnCloseReq, resp *rpc.ConnCloseResp) error {
	c.logger.Debug("Conn close",
		zap.String("src", args.Source),
		zap.String("direction", args.Direction.String()),
		zap.String("dest", args.Destination))
	//dbModel := c.generatePacketDBModel(args.Session, args.Packet)
	//c.rpcService.database.AppendTcpPacket(dbModel)
	return nil
}

func (c *CocoonHandler) Analysis(ctx context.Context, args *rpc.AnalysisReq, resp *rpc.AnalysisResp) error {
	c.logger.Debug("Handle analysis request", zap.String("session", args.Session))
	err := c.rpcService.dissectManager.Dissect(args.Session, c.rpcService.database)
	resp.Error = err
	return nil
}

// RequestOutbound 请求 mock 数据
func (c *CocoonHandler) RequestOutbound(ctx context.Context, args *rpc.OutboundReq, resp *rpc.OutboundResp) error {
	c.logger.Debug("Handle outbound request",
		zap.String("session", args.Session),
		zap.String("proto", args.Proto.String()))
	//body := []byte("HTTP/1.1 200 OK\r\nContent-Length: 7\r\n\r\nMock OK")
	//resp.Body = &body
	resp.OpType = rpc.OP_PASS
	return nil
}

func (c *CocoonHandler) RecordRequestResponse(ctx context.Context, args *rpc.RecordReq, resp *rpc.RecordResp) error {
	c.logger.Debug("Handle record call",
		zap.String("session", args.Session),
		zap.String("proto", args.Proto.String()))
	dbModel := c.generateRecordDBModel(args)
	c.rpcService.database.AppendRecord(dbModel)
	return nil
}

func (c *CocoonHandler) generateRecordDBModel(args *rpc.RecordReq) *db.Record {
	record := &db.Record{
		Session:    args.Session,
		IsOutgoing: args.IsOutgoing,
		Proto:      args.Proto.Name,
		ReqHeader:  args.ReqHeader,
		RespHeader: args.RespHeader,
	}
	if args.ReqBody != nil {
		record.ReqBody = string(*args.ReqBody)
	}
	if args.RespBody != nil {
		record.RespBody = string(*args.RespBody)
	}
	return record
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
