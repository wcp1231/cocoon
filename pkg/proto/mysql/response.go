package mysql

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/mysql/proto"
	"errors"
	"fmt"
)

func (d *Dissector) popRequest() *MysqlMessage {
	ele := d.flyingRequests.Front()
	d.flyingRequests.Remove(ele)
	return ele.Value.(*MysqlMessage)
}

func (d *Dissector) readResponse() (common.Message, error) {
	message, err := d.readCmdQueryResponse()
	if err != nil {
		return nil, err
	}

	if message.HasError() {
		return message, nil
	}
	if message.HasStmtPrepareOk() {
		return message, nil
	}
	// ok_packet
	if message.HasOK() {
		return message, nil
	}

	if !message.HasResultSet() {
		return message, nil
	}

	for _, row := range message.GetResultSet().Rows {
		for _, val := range row {
			fmt.Printf("%v ", val)
		}
		fmt.Println()
	}
	return message, nil
}

func (d *Dissector) readCmdQueryResponse() (*MysqlMessage, error) {
	message := NewMysqlGenericMessage()
	pkt, err := d.respStream.NextPacket()
	if err != nil {
		return nil, err
	}
	data := pkt.Datas

	req := d.popRequest()
	message.SetRequestMessage(req)
	reqCmd := req.Meta["OP_TYPE"]

	// COM_STMT_PREPARE_OK 和 OK_PACKET 似乎很像
	if reqCmd == "COM_STMT_PREPARE" {
		stmtPrepareOk, err := proto.UnPackStmtPrepareOk(data, d.respStream)
		d.stmtParamsMap[stmtPrepareOk.StatementId] = stmtPrepareOk.ParamsCount
		if err != nil {
			return nil, err
		}
		message.SetStmtPrepareOk(stmtPrepareOk)
		return message, nil
	}

	switch data[0] {
	case proto.OK_PACKET:
		ok, err := proto.UnPackOK(data)
		if err != nil {
			return nil, err
		}
		message.SetOk(ok)
		return message, nil
	case proto.ERR_PACKET:
		mysqlError := proto.UnPackERR(data)
		message.SetError(mysqlError)
		return message, nil
	case 0xfb:
		return nil, errors.New("Local.infile.not.implemented")
	}

	// ResultSet
	var rowMode proto.RowMode
	if reqCmd == "COM_STMT_EXECUTE" {
		rowMode = proto.BinaryRowMode
	} else {
		rowMode = proto.TextRowMode
	}
	resultSet, err := proto.UnPackResultSet(rowMode, d.capabilities, data, d.respStream)
	if err != nil {
		return nil, err
	}
	message.SetResultSet(resultSet)
	return message, nil
}
