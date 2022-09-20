package common

type MysqlMessage struct {
	GenericMessage
}

func NewMysqlGenericMessage() *MysqlMessage {
	return &MysqlMessage{
		NewGenericMessage(PROTOCOL_MYSQL.Name),
	}
}
