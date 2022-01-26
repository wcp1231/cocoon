package agent

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/common"
	"cocoon/pkg/model/rpc"
	"cocoon/pkg/proto"
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"strings"
	"time"
)

const maxPacketLen = 0xFFFF

type DumpHandler struct {
	server       *Server
	inboundConn  *net.TCPConn
	outboundConn *net.TCPConn
	inboundAddr  string
	outboundAddr string
	ctx          context.Context
	Close        context.CancelFunc
	seqNum       uint64

	pc    *proto.ProtoClassifier
	proto *common.Protocol
}

func NewDumpHandler(server *Server, inboundConn, outboundConn *net.TCPConn) *DumpHandler {
	innerCtx, close := context.WithCancel(server.ctx)
	pc := proto.NewProtoClassifier()
	return &DumpHandler{
		server:       server,
		inboundConn:  inboundConn,
		outboundConn: outboundConn,
		inboundAddr:  inboundConn.RemoteAddr().String(),
		outboundAddr: outboundConn.RemoteAddr().String(),
		ctx:          innerCtx,
		Close:        close,
		seqNum:       0,

		pc:    pc,
		proto: common.PROTOCOL_UNKNOWN,
	}
}

func (d *DumpHandler) Start() {
	defer func() {
		if err := d.inboundConn.Close(); err != nil {
			d.server.logger.WithOptions(zap.AddCaller()).Error("proxy inbound conn Close error")
		}
		if err := d.outboundConn.Close(); err != nil {
			d.server.logger.WithOptions(zap.AddCaller()).Error("proxy outbound conn Close error")
		}
	}()
	go d.pipe(d.inboundConn, d.outboundConn, common.ClientToRemote)
	go d.pipe(d.outboundConn, d.inboundConn, common.RemoteToClient)
	select {
	case <-d.ctx.Done():
		return
	}
}

func (d *DumpHandler) pipe(srcConn, dstConn *net.TCPConn, direction common.Direction) {
	defer d.Close()

	buff := make([]byte, maxPacketLen)
	var longB []byte
	for {
		n, err := srcConn.Read(buff)
		if err != nil {
			if err.Error() != "EOF" && !strings.Contains(err.Error(), "use of closed network connection") {
				fields := d.fieldsWithErrorAndDirection(err, direction)
				d.server.logger.WithOptions(zap.AddCaller()).Error("strCon Read error", fields...)
			}
			d.handleClose(direction)
			break
		}

		b := buff[:n]

		if d.proto == common.PROTOCOL_UNKNOWN && direction == common.ClientToRemote {
			d.proto = d.pc.Classify(d.outboundAddr, bufio.NewReader(bytes.NewBuffer(b)))
		}

		if d.proto.Dump {
			if n == maxPacketLen && buff[n-1] != 0x00 {
				longB = append(longB, b...)
			} else {
				if len(longB) > 0 {
					longB = append(longB, b...)
					err = d.dump(longB, direction)
					longB = nil
				} else {
					err = d.dump(b, direction)
				}
				if err != nil {
					fields := d.fieldsWithErrorAndDirection(err, direction)
					d.server.logger.WithOptions(zap.AddCaller()).Error("dumper Dump error", fields...)
					break
				}
			}
		}

		if _, err := dstConn.Write(b); err != nil {
			fields := d.fieldsWithErrorAndDirection(err, direction)
			d.server.logger.WithOptions(zap.AddCaller()).Error("destCon Write error", fields...)
			break
		}

		select {
		case <-d.ctx.Done():
			break
		default:
			d.seqNum++
		}
	}
}

func (d *DumpHandler) dump(b []byte, direction common.Direction) error {
	d.server.logger.Info("Dump",
		zap.String("i", d.inboundAddr),
		zap.String("d", direction.String()),
		zap.String("o", d.outboundAddr),
		zap.String("p", d.proto.String()))
	packte := common.TcpPacket{
		Source:      d.inboundAddr,
		Destination: d.outboundAddr,
		IsOutgoing:  true, // TODO 这里暂时都是向外的流量
		Direction:   &direction,
		Seq:         d.seqNum,
		Timestamp:   time.Now(),
		Protocol:    d.proto.Name,
		Payload:     b,
	}
	req := rpc.UploadReq{
		Session: d.server.session,
		Packet:  &packte,
	}
	d.server.rpcClient.Upload(d.ctx, &req)
	return nil
}

func (d *DumpHandler) handleClose(direction common.Direction) {
	req := &rpc.ConnCloseReq{
		Source:      d.inboundAddr,
		Destination: d.outboundAddr,
		Direction:   &direction,
	}
	d.server.rpcClient.ConnClose(d.ctx, req)
}

func (d *DumpHandler) fieldsWithErrorAndDirection(err error, direction common.Direction) []zapcore.Field {
	fields := []zapcore.Field{
		zap.Error(err),
		zap.Uint64("conn_seq_num", d.seqNum),
		zap.String("direction", direction.String()),
	}

	//for _, kv := range d.connMetadata.DumpValues {
	//	fields = append(fields, zap.Any(kv.Key, kv.Value))
	//}

	return fields
}
