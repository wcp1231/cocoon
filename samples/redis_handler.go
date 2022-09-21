package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

func (s *SampleApp) redisString(w http.ResponseWriter, req *http.Request) {
	var ctx = context.Background()
	setCount := int(RandomInt64Range(5, 20))
	for i := 0; i < setCount; i++ {
		key := fmt.Sprintf("key_%d", i)
		val := RandomAlphaString(10)
		s.redis.Set(ctx, key, val, time.Duration(1)*time.Minute)
	}

	var result []interface{}
	getCount := int(RandomInt64Range(15, 30))
	for i := 0; i < getCount; i++ {
		key := fmt.Sprintf("key_%d", i)
		val, err := s.redis.Get(ctx, key).Result()
		result = append(result, struct {
			Val string
			Err error
		}{
			Val: val,
			Err: err,
		})
	}
	responseOk(w, result)
}
func (s *SampleApp) redisZSet(w http.ResponseWriter, req *http.Request) {
	var ctx = context.Background()
	setCount := int(RandomInt64Range(5, 20))
	for i := 0; i < setCount; i++ {
		member := fmt.Sprintf("mem_%d", i)
		score := RandomFloat64Range(0, 10)
		s.redis.ZAdd(ctx, "zset_test", &redis.Z{
			Member: member,
			Score:  score,
		})
	}

	members := s.redis.ZRangeWithScores(ctx, "zset_test", 0, 100).Val()
	responseOk(w, members)
}
