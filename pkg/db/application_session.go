package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (d *Database) EnsureApplicationAndSession(app, session string) {
	opts := options.UpdateOptions{}
	opts.SetUpsert(true)
	_, err := d.db.Collection(COLLECTION_APPLICATION).
		UpdateOne(d.ctx, bson.M{"name": app}, bson.M{"$set": bson.M{"name": app}}, &opts)
	if err != nil {
		d.logger.Warn("Ensure application failed", zap.Error(err))
	}
	_, err = d.db.Collection(COLLECTION_APPLICATION_SESSION).
		UpdateOne(d.ctx, bson.M{"app": app}, bson.M{"$set": bson.M{"app": app, "session": session}}, &opts)
	if err != nil {
		d.logger.Warn("Ensure application session failed", zap.Error(err))
	}
}
