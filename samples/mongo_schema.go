package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func NewRandomMongoRecord() bson.D {
	result := bson.D{}
	strKyes := int(RandomInt64Range(1, 5))
	for i := 0; i < strKyes; i++ {
		key := RandomAlphaString(5)
		val := RandomAlphaString(10)
		result = append(result, bson.E{Key: key, Value: val})
	}
	numKyes := int(RandomInt64Range(1, 5))
	for i := 0; i < numKyes; i++ {
		key := RandomAlphaString(5)
		val := RandomFloat64Range(0.1, 99.9)
		result = append(result, bson.E{Key: key, Value: val})
	}
	result = append(result, bson.E{Key: "date_field", Value: RandomDate()})
	return result
}
func NewSimpleMongoRecord() bson.D {
	result := bson.D{}
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC)
	result = append(result, bson.E{Key: "date_field", Value: max})
	return result
}
