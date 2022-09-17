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
	reqCmd    string
	colNumber int
	cols      []*proto.Field
	rows      [][]proto.Value
	raw       []byte
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
		//message.Raw = &resp.raw
		errbytes := packet.ToPacketBytes(proto.PackERR(resp.error))
		message.Raw = &errbytes
		// TODO 返回给客户端？
		return message, nil
	}
	// 非 ResultSet
	if resp.colNumber == 0 {
		okbytes := packet.ToPacketBytes(proto.PackOK(resp.ok))
		//fmt.Printf("ok packet compare. %v\n", bytes.Compare(okbytes, resp.raw))
		//message.Raw = &resp.raw
		message.Raw = &okbytes
		return message, nil
	}

	for _, row := range resp.rows {
		for _, val := range row {
			fmt.Printf("%v ", val)
		}
		fmt.Println()
	}
	message.Raw = &resp.raw
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
	resp.raw = pkt.Raw()
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
	err = d.readResultSet(data, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *Dissector) readResultSet(data []byte, resp *queryResponse) error {
	// ResultSet
	number, err := proto.ColumnCount(data)
	if err != nil {
		return err
	}
	resp.colNumber = int(number)
	if resp.colNumber > 0 {
		columns, err := d.readColumns(resp)
		if err != nil {
			return err
		}
		resp.cols = columns
		if d.capabilities&proto.CLIENT_DEPRECATE_EOF == 0 {
			if err = d.readEOF(resp); err != nil {
				return err
			}
		}
	}

	if resp.reqCmd == "COM_STMT_EXECUTE" {
		err := d.readBinaryResultSet(resp)
		if err != nil {
			return err
		}
	} else if resp.reqCmd == "COM_QUERY" {
		err := d.readTextResultSet(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dissector) readBinaryResultSet(resp *queryResponse) error {
	rowReader := proto.NewBinaryRows(d.respStream)
	rowReader.RowsAffected = resp.ok.AffectedRows
	rowReader.InsertID = resp.ok.LastInsertID
	rowReader.Fields = resp.cols
	var rows [][]proto.Value
	for rowReader.Next() {
		values, err := rowReader.RowValues()
		if err != nil {
			return err
		}
		rows = append(rows, values)
	}
	resp.rows = rows
	resp.raw = append(resp.raw, rowReader.Raw()...)
	return nil
}

func (d *Dissector) readTextResultSet(resp *queryResponse) error {
	rowReader := proto.NewTextRows(d.respStream)
	rowReader.Fields = resp.cols
	var rows [][]proto.Value
	for rowReader.Next() {
		values, err := rowReader.RowValues()
		if err != nil {
			return err
		}
		rows = append(rows, values)

	}
	resp.rows = rows
	resp.raw = append(resp.raw, rowReader.Raw()...)
	return nil
}

func (d *Dissector) readColumns(resp *queryResponse) ([]*proto.Field, error) {
	var err error
	var pkt *packet.Packet
	columns := make([]*proto.Field, 0, resp.colNumber)
	for i := 0; i < resp.colNumber; i++ {
		if pkt, err = d.respStream.NextPacket(); err != nil {
			return nil, err
		}
		column, err := proto.UnpackColumn(pkt.Datas)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
		resp.raw = append(resp.raw, pkt.Raw()...)
	}
	return columns, nil
}

func (d *Dissector) readEOF(resp *queryResponse) error {
	pkt, err := d.respStream.NextPacket()
	if err != nil {
		return err
	}
	data := pkt.Datas
	resp.raw = append(resp.raw, pkt.Raw()...)
	switch data[0] {
	case proto.EOF_PACKET:
		return nil
	case proto.ERR_PACKET:
		return proto.UnPackERR(data).ToError()
	default:
		return errors.New(fmt.Sprintf("unexpected.eof.packet[%+v]", data))
	}
}
