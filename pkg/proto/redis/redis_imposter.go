package redis

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/model/mock"
	"strconv"
)

type RedisRequestMatcher struct {
	id     int32
	config mock.RedisMockConfig
	cmd    mock.FieldMatcher
	key    mock.FieldMatcher
	args   []mock.FieldMatcher

	respType string
	respBody RedisObject
}

func NewRedisRequestMatcherFromConfig(config mock.RedisMockConfig, id int32) *RedisRequestMatcher {
	matcher := &RedisRequestMatcher{
		id:     id,
		config: config,
	}
	if config.Request.Cmd != nil {
		matcher.cmd = mock.NewFieldMatcher(config.Request.Cmd)
	}
	if config.Request.Key != nil {
		matcher.key = mock.NewFieldMatcher(config.Request.Key)
	}

	matcher.respType = config.Response.Type
	matcher.respBody = encodeRedisMockData(config.Response)

	return matcher
}

func encodeRedisMockData(resp mock.RedisResponseObject) RedisObject {
	switch resp.Type {
	case "null":
		return &RedisBulkString{Len: -1}
	case "error":
		return &RedisError{Error: resp.Value}
	case "integer":
		value, _ := strconv.ParseInt(resp.Value, 10, 64)
		return &RedisInteger{Integer: value}
	case "string":
		return &RedisBulkString{
			Len:    int64(len(resp.Value)),
			String: resp.Value,
		}
	case "array":
		array := &RedisArray{}
		array.Len = len(resp.Array)
		array.Items = make([]RedisObject, array.Len)
		for i, item := range resp.Array {
			obj := encodeRedisMockData(item)
			array.Items[i] = obj
		}
		return array
	}
	return nil
}

func (h *RedisRequestMatcher) Match(r common.Message) bool {
	req := r.(*RedisMessage)
	if h.cmd != nil {
		if !h.cmd.Match(req.GetCmd()) {
			return false
		}
	}

	if h.key != nil {
		if !h.key.Match(req.GetKey()) {
			return false
		}
	}

	return true
}

func (h *RedisRequestMatcher) Data() common.Message {
	message := NewRedisGenericMessage()
	message.SetMock()
	message.SetResponseObj(h.respBody.Pretty())
	raw := h.respBody.Raw()
	message.Raw = &raw
	return message
}

func (h *RedisRequestMatcher) ID() int32 {
	return h.id
}

func (h *RedisRequestMatcher) GetConfig() interface{} {
	h.config.Id = h.id
	return h.config
}
