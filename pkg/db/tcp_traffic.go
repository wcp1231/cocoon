package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (d *Database) AppendTcpPacket(packet *TcpTraffic) {
	_, err := d.db.Collection(COLLECTION_TCP_TRAFFIC).InsertOne(d.ctx, packet)
	if err != nil {
		d.logger.Warn("Append tcp packet failed", zap.Error(err))
	}
}

func (d *Database) ReadTcpTrafficBySession(session string) (*mongo.Cursor, error) {
	opts := &options.FindOptions{}
	opts.SetSort(bson.D{{"_id", 1}})
	return d.db.Collection(COLLECTION_TCP_TRAFFIC).Find(d.ctx, bson.M{"session": session}, opts)
}
