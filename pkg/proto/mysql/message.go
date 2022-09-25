package mysql

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/mysql/packet"
	"cocoon/pkg/proto/mysql/proto"
)

const (
	MYSQL_OP_KEY              = "MYSQL_OP_TYPE"
	MYSQL_QUERY_KEY           = "MYSQL_QUERY"
	MYSQL_ERR_PACKET_KEY      = "MYSQL_ERR_PACKET"
	MYSQL_OK_PACKET_KEY       = "MYSQL_OK_PACKET"
	MYSQL_STMT_EXECUTE_KEY    = "MYSQL_STMT_EXECUTE_PACKET"
	MYSQL_STMT_PREPARE_OK_KEY = "MYSQL_STMT_PREPARE_OK_PACKET"
	MYSQL_RESULT_SET_KEY      = "MYSQL_RESULT_SET_PACKET"
)

type MysqlMessage struct {
	common.GenericMessage

	requestMessage *MysqlMessage // response 处理需要原始 request
	requestBytes   []byte        // request 发送的原始数据
}

func NewMysqlGenericMessage() *MysqlMessage {
	return &MysqlMessage{
		GenericMessage: common.NewGenericMessage(common.PROTOCOL_MYSQL.Name),
	}
}

func (m *MysqlMessage) SetRequestMessage(request *MysqlMessage) {
	m.requestMessage = request
}
func (m *MysqlMessage) GetRequestMessage() *MysqlMessage {
	return m.requestMessage
}
func (m *MysqlMessage) SetRequestBytes(requestBytes []byte) {
	m.requestBytes = requestBytes
}

func (m *MysqlMessage) SetOpType(opType string) {
	m.Meta["OP_TYPE"] = opType
	m.Payload[MYSQL_OP_KEY] = opType
}
func (m *MysqlMessage) HasOpType() bool {
	return m.Payload[MYSQL_OP_KEY] != nil
}
func (m *MysqlMessage) GetOpType() string {
	return m.Payload[MYSQL_OP_KEY].(string)
}
func (m *MysqlMessage) SetQuery(query string) {
	m.Payload[MYSQL_QUERY_KEY] = query
}
func (m *MysqlMessage) GetQuery() string {
	return m.Payload[MYSQL_QUERY_KEY].(string)
}
func (m *MysqlMessage) SetError(err *proto.ERR) {
	m.Payload[MYSQL_ERR_PACKET_KEY] = err
}
func (m *MysqlMessage) HasError() bool {
	return m.Payload[MYSQL_ERR_PACKET_KEY] != nil
}
func (m *MysqlMessage) GetError() *proto.ERR {
	return m.Payload[MYSQL_ERR_PACKET_KEY].(*proto.ERR)
}
func (m *MysqlMessage) SetOk(ok *proto.OK) {
	m.Payload[MYSQL_OK_PACKET_KEY] = ok
}
func (m *MysqlMessage) HasOK() bool {
	return m.Payload[MYSQL_OK_PACKET_KEY] != nil
}
func (m *MysqlMessage) GetOk() *proto.OK {
	return m.Payload[MYSQL_OK_PACKET_KEY].(*proto.OK)
}
func (m *MysqlMessage) SetStmtExecute(ok *proto.StmtExecute) {
	m.Payload[MYSQL_STMT_EXECUTE_KEY] = ok
}
func (m *MysqlMessage) SetStmtPrepareOk(ok *proto.StmtPrepareOK) {
	m.Payload[MYSQL_STMT_PREPARE_OK_KEY] = ok
}
func (m *MysqlMessage) HasStmtPrepareOk() bool {
	return m.Payload[MYSQL_STMT_PREPARE_OK_KEY] != nil
}
func (m *MysqlMessage) GetStmtPrepareOk() *proto.StmtPrepareOK {
	return m.Payload[MYSQL_STMT_PREPARE_OK_KEY].(*proto.StmtPrepareOK)
}
func (m *MysqlMessage) SetResultSet(ok *proto.ResultSet) {
	m.Payload[MYSQL_RESULT_SET_KEY] = ok
}
func (m *MysqlMessage) HasResultSet() bool {
	return m.Payload[MYSQL_RESULT_SET_KEY] != nil
}
func (m *MysqlMessage) GetResultSet() *proto.ResultSet {
	return m.Payload[MYSQL_RESULT_SET_KEY].(*proto.ResultSet)
}
func (m *MysqlMessage) GetRaw() []byte {
	// request
	if m.HasOpType() {
		return m.requestBytes
	}

	// response
	if m.HasError() {
		return packet.ToPacketBytes(proto.PackERR(m.GetError()))
	}
	if m.HasStmtPrepareOk() {
		return proto.PackStmtPrepareOk(m.GetStmtPrepareOk())
	}
	if m.HasOK() {
		return packet.ToPacketBytes(proto.PackOK(m.GetOk()))
	}
	if m.HasResultSet() {
		return proto.PackResultSet(m.GetResultSet())
	}
	return nil
}
