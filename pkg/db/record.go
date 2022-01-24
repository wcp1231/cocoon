package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (d *Database) AppendRecord(record *Record) {
	_, err := d.db.Collection(COLLECTION_RECORDS).InsertOne(d.ctx, record)
	if err != nil {
		d.logger.Warn("Append record failed", zap.Error(err))
	}
}

func (d *Database) ReadRecordsBySession(session string) (*mongo.Cursor, error) {
	opts := &options.FindOptions{}
	opts.SetSort(bson.D{{"_id", 1}})
	return d.db.Collection(COLLECTION_RECORDS).Find(d.ctx, bson.M{"session": session}, opts)
}
