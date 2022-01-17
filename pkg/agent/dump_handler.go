package agent

import (
	"cocoon/pkg/model"
	"cocoon/pkg/model/rpc"
	"context"
	"fmt"
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
}

func NewDumpHandler(server *Server, inboundConn, outboundConn *net.TCPConn) *DumpHandler {
	innerCtx, close := context.WithCancel(server.ctx)
	return &DumpHandler{
		server:       server,
		inboundConn:  inboundConn,
		outboundConn: outboundConn,
		inboundAddr:  inboundConn.RemoteAddr().String(),
		outboundAddr: outboundConn.RemoteAddr().String(),
		ctx:          innerCtx,
		Close:        close,
		seqNum:       0,
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
	go d.pipe(d.inboundConn, d.outboundConn, model.ClientToRemote)
	go d.pipe(d.outboundConn, d.inboundConn, model.RemoteToClient)
	select {
	case <-d.ctx.Done():
		return
	}
}

func (d *DumpHandler) pipe(srcConn, dstConn *net.TCPConn, direction model.Direction) {
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
			break
		}

		b := buff[:n]
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

func (d *DumpHandler) dump(b []byte, direction model.Direction) error {
	//kvs := []dumper.DumpValue{
	//	dumper.DumpValue{
	//		Key:   "conn_seq_num",
	//		Value: p.seqNum,
	//	},
	//	dumper.DumpValue{
	//		Key:   "direction",
	//		Value: direction.String(),
	//	},
	//	dumper.DumpValue{
	//		Key:   "ts",
	//		Value: time.Now(),
	//	},
	//}

	//return p.agent.dumper.Dump(b, direction, p.connMetadata, kvs)
	fmt.Printf("Seq[%d] %s %s %s Dump:\n", d.seqNum, d.inboundAddr, direction.String(), d.outboundAddr)
	packte := rpc.TcpPacket{
		Source:      d.inboundAddr,
		Destination: d.outboundAddr,
		IsOutgoing:  true, // TODO 这里暂时都是向外的流量
		Direction:   &direction,
		Seq:         d.seqNum,
		Timestamp:   time.Now(),
		Payload:     b,
	}
	req := rpc.UploadReq{
		Session: d.server.session,
		Packet:  &packte,
	}
	d.server.rpcClient.Upload(d.ctx, &req)
	return nil
}

func (d *DumpHandler) fieldsWithErrorAndDirection(err error, direction model.Direction) []zapcore.Field {
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
