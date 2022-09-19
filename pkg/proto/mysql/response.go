package mysql

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/mysql/packet"
	"cocoon/pkg/proto/mysql/proto"
	"errors"
	"fmt"
)

type queryResponse struct {
	// OK_PACKET
	ok *proto.OK
	// ERR_PACKET
	error *proto.ERR

	// STMT_PREPARE_OK
	stmtPrepareOk *proto.StmtPrepareOK

	// RESULT_SET
	resultSet *proto.ResultSet
	reqCmd    string
}

func (d *Dissector) popRequest() *common.GenericMessage {
	ele := d.flyingRequests.Front()
	d.flyingRequests.Remove(ele)
	return ele.Value.(*common.GenericMessage)
}

func (d *Dissector) readResponse() (*common.GenericMessage, error) {
	//fmt.Println("Mysql try read response")
	message := common.NewMysqlGenericMessage()
	resp, err := d.readCmdQueryResponse()
	if err != nil {
		return nil, err
	}

	if resp.error != nil {
		if resp.error.InternalError != nil {
			return nil, resp.error.InternalError
		}
		errbytes := packet.ToPacketBytes(proto.PackERR(resp.error))
		message.Raw = &errbytes
		// TODO 返回给客户端？
		return message, nil
	}
	if resp.stmtPrepareOk != nil {
		spbytes := proto.PackStmtPrepareOk(resp.stmtPrepareOk)
		message.Raw = &spbytes
		return message, nil
	}
	// ok_packet
	if resp.resultSet == nil {
		okbytes := packet.ToPacketBytes(proto.PackOK(resp.ok))
		message.Raw = &okbytes
		return message, nil
	}

	for _, row := range resp.resultSet.Rows {
		for _, val := range row {
			fmt.Printf("%v ", val)
		}
		fmt.Println()
	}
	raw := proto.PackResultSet(resp.resultSet, d.capabilities)
	message.Raw = &raw
	return message, nil
}

func (d *Dissector) readCmdQueryResponse() (*queryResponse, error) {

	pkt, err := d.respStream.NextPacket()
	if err != nil {
		return nil, err
	}
	data := pkt.Datas

	req := d.popRequest()
	resp := &queryResponse{}
	resp.reqCmd = req.Meta["OP_TYPE"]

	// COM_STMT_PREPARE_OK 和 OK_PACKET 似乎很像
	if resp.reqCmd == "COM_STMT_PREPARE" {
		stmtPrepareOk, err := proto.UnPackStmtPrepareOk(data, d.respStream)
		if err != nil {
			return nil, err
		}
		resp.stmtPrepareOk = stmtPrepareOk
		return resp, nil
	}

	resp.ok = &proto.OK{}
	switch data[0] {
	case proto.OK_PACKET:
		ok, err := proto.UnPackOK(data)
		if err != nil {
			return nil, err
		}
		resp.ok = ok
		return resp, nil
	case proto.ERR_PACKET:
		resp.error = proto.UnPackERR(data)
		return resp, nil
	case 0xfb:
		resp.error = &proto.ERR{
			InternalError: errors.New("Local.infile.not.implemented"),
		}
		return resp, nil
	}

	// ResultSet
	resultSet, err := proto.UnPackResultSet(resp.reqCmd, d.capabilities, data, d.respStream)
	if err != nil {
		return nil, err
	}
	resp.resultSet = resultSet
	return resp, nil
}
