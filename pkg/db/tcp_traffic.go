package db

import (
	"go.uber.org/zap"
)

func (d *Database) AppendTcpPacket(packet *TcpTraffic) {
	_, err := d.db.Collection(COLLECTION_TCP_TRAFFIC).InsertOne(d.ctx, packet)
	if err != nil {
		d.logger.Warn("Append tcp packet failed", zap.Error(err))
	}
}
