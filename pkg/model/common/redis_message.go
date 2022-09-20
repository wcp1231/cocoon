package common

const (
	REDIS_CMD_KEY = "REDIS_CMD"
	REDIS_KEY_KEY = "REDIS_KEY"
	REDIS_PAYLOAD = "REDIS_PAYLOAD"
)

type RedisMessage struct {
	GenericMessage
}

func NewRedisGenericMessage() *RedisMessage {
	return &RedisMessage{
		NewGenericMessage(PROTOCOL_REDIS.Name),
	}
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
func (m *RedisMessage) SetRedisPayload(payload []byte) {
	m.Body = &payload
	m.Payload[REDIS_PAYLOAD] = string(payload)
}
func (m *RedisMessage) GetRedisPayload() string {
	return m.Payload[REDIS_PAYLOAD].(string)
}
