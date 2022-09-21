package redis

import "cocoon/pkg/model/common"

const (
	REDIS_REQ_CMD_KEY   = "REDIS_REQ_CMD"
	REDIS_RESP_OBJ_KEY  = "REDIS_RESP_OBJ"
	REDIS_HEARTBEAT_KEY = "REDIS_HEARTBEAT"
	REDIS_CMD_KEY       = "REDIS_CMD"
	REDIS_KEY_KEY       = "REDIS_KEY"
)

type RedisMessage struct {
	common.GenericMessage
}

func NewRedisGenericMessage() *RedisMessage {
	return &RedisMessage{
		GenericMessage: common.NewGenericMessage(common.PROTOCOL_REDIS.Name),
	}
}

func (m *RedisMessage) SetRequestCmd(request string) {
	m.Payload[REDIS_REQ_CMD_KEY] = request
}
func (m *RedisMessage) GetRequestCmd() string {
	return m.Payload[REDIS_REQ_CMD_KEY].(string)
}
func (m *RedisMessage) SetResponseObj(obj string) {
	m.Payload[REDIS_RESP_OBJ_KEY] = obj
}
func (m *RedisMessage) GetResponseObj() string {
	return m.Payload[REDIS_RESP_OBJ_KEY].(string)
}
func (m *RedisMessage) SetHeartbeat() {
	m.Payload[REDIS_HEARTBEAT_KEY] = true
}
func (m *RedisMessage) IsHeartbeat() bool {
	hb, exist := m.Payload[REDIS_HEARTBEAT_KEY]
	return exist && hb.(bool)
}

func (m *RedisMessage) SetCmd(cmd string) {
	m.Meta["CMD"] = cmd
	m.Payload[REDIS_CMD_KEY] = cmd
}
func (m *RedisMessage) GetCmd() string {
	return m.Payload[REDIS_CMD_KEY].(string)
}
func (m *RedisMessage) SetKey(key string) {
	m.Meta["KEY"] = key
	m.Payload[REDIS_KEY_KEY] = key
}
func (m *RedisMessage) GetKey() string {
	return m.Payload[REDIS_KEY_KEY].(string)
}
