package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const COCOON_DB = "cocoon"

type Database struct {
	logger *zap.Logger
	ctx    context.Context
	client *mongo.Client
	db     *mongo.Database
}

func NewDatabase(ctx context.Context, logger *zap.Logger, uri string) *Database {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Fatal("Connect mongo failed", zap.Error(err))
		panic(err)
	}
	db := client.Database(COCOON_DB)
	return &Database{
		logger: logger,
		ctx:    ctx,
		client: client,
		db:     db,
	}
}
