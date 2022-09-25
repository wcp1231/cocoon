package mongo

import (
	"cocoon/pkg/model/common"
)

const (
	MONGO_HEADER_KEY  = "MONGO_HEADER"
	MONGO_MESSAEG_KEY = "MONGO_MESSAGE"
	MONGO_OP_KEY      = "MONGO_OP_TYPE"
)

type MongoMessage struct {
	common.GenericMessage
}

func NewMongoGenericMessage() *MongoMessage {
	return &MongoMessage{
		common.NewGenericMessage(common.PROTOCOL_MONGO.Name),
	}
}

func (m *MongoMessage) SetMongoHeader(header MessageHeader) {
	m.Payload[MONGO_HEADER_KEY] = header
}
func (m *MongoMessage) GetMongoHeader() MessageHeader {
	return m.Payload[MONGO_HEADER_KEY].(MessageHeader)
}
func (m *MongoMessage) SetMongoMessage(message interface{}) {
	m.Payload[MONGO_MESSAEG_KEY] = message
}
func (m *MongoMessage) GetMongoMessage() interface{} {
	return m.Payload[MONGO_MESSAEG_KEY]
}
func (m *MongoMessage) SetOpType(opType string) {
	m.Payload[MONGO_OP_KEY] = opType
	m.Meta["OP_TYPE"] = opType
}
func (m *MongoMessage) GetOpType() string {
	return m.Payload[MONGO_OP_KEY].(string)
}
func (m *MongoMessage) GetRaw() []byte {
	bs, _ := Dump(m)
	return bs
}
