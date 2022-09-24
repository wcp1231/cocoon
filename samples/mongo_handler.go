package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

const (
	MONGO_DATABASE   = "cocoon"
	MONGO_COLLECTION = "samples"
)

func (s *SampleApp) mongoFind(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	collection := s.mongo.Database(MONGO_DATABASE).Collection(MONGO_COLLECTION)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		responseErr(w, err)
		return
	}
	defer cur.Close(ctx)
	var samples []bson.D
	for cur.Next(ctx) {
		var item bson.D
		err = cur.Decode(&item)
		if err != nil {
			responseErr(w, err)
			return
		}
		samples = append(samples, item)
	}
	responseOk(w, samples)
}
func (s *SampleApp) mongoInsert(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	collection := s.mongo.Database(MONGO_DATABASE).Collection(MONGO_COLLECTION)
	record := NewRandomMongoRecord()
	res, err := collection.InsertOne(ctx, record)
	if err != nil {
		responseErr(w, err)
		return
	}
	responseOk(w, struct {
		InsertedID interface{}
	}{
		InsertedID: res.InsertedID,
	})
}
func (s *SampleApp) mongoRemove(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	collection := s.mongo.Database(MONGO_DATABASE).Collection(MONGO_COLLECTION)
	res, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		responseErr(w, err)
		return
	}
	responseOk(w, struct {
		DeletedCount int64
	}{
		DeletedCount: res.DeletedCount,
	})
}
