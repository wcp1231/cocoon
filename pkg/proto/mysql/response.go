package mysql

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/mysql/packet"
	"cocoon/pkg/proto/mysql/proto"
	"errors"
	"fmt"
)

type queryResponse struct {
	ok        *proto.OK
	colNumber int
	error     error
	cols      []*proto.Field
	vals      [][]byte
	raw       []byte
}

func (d *Dissector) readResponse() (*common.GenericMessage, error) {
	fmt.Println("Mysql try read request")
	message := common.NewMysqlGenericMessage()
	resp, err := d.readCmdQueryResponse()
	if err != nil {
		return nil, err
	}
	if resp.error != nil {
		message.Raw = &resp.raw
		// TODO 返回给客户端？
		return message, nil
	}
	if resp.colNumber > 0 {
		columns, err := d.readColumns(resp)
		if err != nil {
			// TODO
		}
		resp.cols = columns

		// TODO Read EOF?
		//if true { //(greeting.Capability & CLIENT_DEPRECATE_EOF) == 0 {
		//	if err = d.readEOF(); err != nil {
		//		return nil, err
		//	}
		//}
	}

	// TODO 如何区分不同类型的 response

	// read all rows
	var vals [][]byte
	rows := proto.NewSimpleRows(d.respStream)
	rows.RowsAffected = resp.ok.AffectedRows
	rows.InsertID = resp.ok.LastInsertID
	rows.Fields = resp.cols
	for rows.Next() {
		row, err := rows.RowValues()
		if err != nil {
			// TODO
			continue
		}
		vals = append(vals, row)
	}
	resp.vals = vals
	resp.raw = append(resp.raw, rows.Raw()...)

	fmt.Printf("Mysql read response. %v\nfields: %v\n%v\n", resp.ok, resp.cols, string(resp.raw))
	message.Raw = &resp.raw
	return message, nil
}

func (d *Dissector) readCmdQueryResponse() (*queryResponse, error) {
	resp := &queryResponse{}
	pkt, err := d.respStream.NextPacket()
	if err != nil {
		return nil, err
	}

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
		resp.error = errors.New("Local.infile.not.implemented")
		return resp, nil
	}
	number, err := proto.ColumnCount(data)
	if err != nil {
		return nil, err
	}
	resp.colNumber = int(number)
	return resp, nil
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

func (d *Dissector) readEOF() error {
	pkt, err := d.respStream.NextPacket()
	if err != nil {
		return err
	}
	data := pkt.Datas
	switch data[0] {
	case proto.EOF_PACKET:
		return nil
	case proto.ERR_PACKET:
		return proto.UnPackERR(data)
	default:
		return errors.New(fmt.Sprintf("unexpected.eof.packet[%+v]", data))
	}
}
