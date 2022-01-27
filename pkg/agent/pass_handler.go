package agent

import (
	"cocoon/pkg/model/common"
	"go.uber.org/zap"
)

type passHandler struct {
	logger *zap.Logger
	proto *common.Protocol
	inboundConn *conn
	outboundConn *conn
}

func newPassHandler(logger *zap.Logger, proto *common.Protocol, inbound, outbound *conn) *passHandler {
	return &passHandler{
		logger: logger,
		proto: proto,
		inboundConn: inbound,
		outboundConn: outbound,
	}
}

func (p *passHandler) start() {
	p.logger.Info("Conn pass",
		zap.String("src", p.inboundConn.addr),
		zap.String("dst", p.outboundConn.addr),
		zap.String("proto", p.proto.Name))
	go p.pipe(p.inboundConn, p.outboundConn, common.ClientToRemote)
	go p.pipe(p.outboundConn, p.inboundConn, common.RemoteToClient)
}

func (c *passHandler) pipe(src, dst *conn, direction common.Direction) {
	for {
		_, err := src.reader.WriteTo(dst.c)
		if err != nil {
			c.logger.Warn("Conn write error",
				zap.String("src", c.inboundConn.addr),
				zap.String("dst", c.outboundConn.addr),
				zap.String("dir", direction.String()),
				zap.Error(err))
			break
		}
	}
}