package dissector

import (
	"cocoon/pkg/db"
	"cocoon/pkg/model/common"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type DissectWorker struct {
	ctx    context.Context
	logger *zap.Logger
	cursor *mongo.Cursor

	streams map[string]*tcpStream
}

func NewDissectWorker(ctx context.Context, logger *zap.Logger, cursor *mongo.Cursor) *DissectWorker {
	return &DissectWorker{
		ctx:    ctx,
		logger: logger,
		cursor: cursor,

		streams: make(map[string]*tcpStream),
	}
}

func (d *DissectWorker) Start() error {
	var traffic db.TcpTraffic
	for d.cursor.Next(d.ctx) {
		err := d.cursor.Decode(&traffic)
		if err != nil {
			return err
		}
		packet := d.toTcpTrafficPacket(&traffic)
		d.accept(packet)
	}
	return d.wait()
}

func (d *DissectWorker) accept(packet *common.TcpPacket) {
	connectionInfo := &common.ConnectionInfo{
		Source:      packet.Source,
		Destination: packet.Destination,
		IsOutgoing:  packet.IsOutgoing,
	}
	stream := d.getOrCreateStream(connectionInfo)
	stream.Accept(packet)
}

func (d *DissectWorker) wait() error {
	for _, s := range d.streams {
		s.Close()
	}
	return nil
}

func (d *DissectWorker) getOrCreateStream(connectionInfo *common.ConnectionInfo) *tcpStream {
	tcpId := connectionInfo.ID()
	stream, ok := d.streams[tcpId]
	if !ok {
		stream = newTcpStream(d.logger, connectionInfo, nil)
		d.streams[tcpId] = stream
	}
	return stream
}

func (d *DissectWorker) toTcpTrafficPacket(traffic *db.TcpTraffic) *common.TcpPacket {
	return &common.TcpPacket{
		Source:      traffic.Source,
		Destination: traffic.Destination,
		IsOutgoing:  traffic.IsOutgoing,
		Direction:   traffic.Direction,
		Seq:         traffic.Seq,
		Timestamp:   traffic.Timestamp,
		Payload:     traffic.Raw,
	}
}
