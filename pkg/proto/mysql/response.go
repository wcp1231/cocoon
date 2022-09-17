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
	// 非 ResultSet
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
	raw := proto.PackResultSet(resp.resultSet)
	message.Raw = &raw
	return message, nil
}

func (d *Dissector) readCmdQueryResponse() (*queryResponse, error) {

	pkt, err := d.respStream.NextPacket()
	if err != nil {
		return nil, err
	}

	req := d.popRequest()
	resp := &queryResponse{}
	resp.reqCmd = req.Meta["OP_TYPE"]

	resp.ok = &proto.OK{}
	data := pkt.Datas
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
