package mock

import (
	"bytes"
	"cocoon/pkg/model/common"
	"fmt"
)

type RedisRequestMatcher struct {
	id int32
	config redisMockConfig
	cmd FieldMatcher
	key FieldMatcher

	respType string
	respBody []byte
}

func newRedisRequestMatcherFromConfig(config redisMockConfig, id int32) *RedisRequestMatcher {
	matcher := &RedisRequestMatcher{
		id: id,
		config: config,
	}
	if config.Request.Cmd != nil {
		matcher.cmd = newFieldMatcher(config.Request.Cmd)
	}
	if config.Request.Key != nil {
		matcher.key = newFieldMatcher(config.Request.Key)
	}

	matcher.respType = config.Response.Type
	matcher.respBody = encodeRedisMockData(config.Response)

	return matcher
}

func encodeRedisMockData(resp redisResponseMockConfig) []byte {
	buf := new(bytes.Buffer)
	switch resp.Type {
	case "nil":
		buf.WriteString("$-1\r\n")
	case "err":
		buf.WriteString(fmt.Sprintf("-%s\r\n", resp.String))
	case "string":
		buf.WriteString(fmt.Sprintf("$%d\r\n", len(resp.String)))
		buf.WriteString(resp.String)
		buf.WriteString("\r\n")
	case "array":
		buf.WriteString(fmt.Sprintf("*%d\r\n", len(resp.Array)))
		// TODO 支持数组中有 null
		for _, str := range resp.Array {
			buf.WriteString(fmt.Sprintf("$%d\r\n", len(str)))
			buf.WriteString(str)
			buf.WriteString("\r\n")
		}
	}
	return buf.Bytes()
}

func (h *RedisRequestMatcher) Match(req *common.GenericMessage) bool {
	if h.cmd != nil {
		if !h.cmd.Match(req.Header["CMD"]) {
			return false
		}
	}

	if h.key != nil {
		if !h.key.Match(req.Header["KEY"]) {
			return false
		}
	}

	return true
}

func (h *RedisRequestMatcher) Data() *[]byte {
	return &h.respBody
}

func (h *RedisRequestMatcher) ID() int32 {
	return h.id
}

func (h *RedisRequestMatcher) GetConfig() interface{} {
	h.config.Id = h.id
	return h.config
}