package mysql

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/mysql/proto"
	"fmt"
	"regexp"
)

// TODO
var removeSpaces = regexp.MustCompile(`\s+`)

func (d *Dissector) readRequest() (common.Message, error) {
	pkt, err := d.reqStream.NextPacket()
	if err != nil {
		return nil, err
	}
	raw := pkt.Raw()
	data := pkt.Datas
	message := NewMysqlGenericMessage()
	message.Raw = &raw
	message.SetOpType(proto.CommandString(data[0]))
	if len(data) > 1 {
		message.SetQuery(string(removeSpaces.ReplaceAll(data[1:], []byte{' '})))
	}

	switch data[0] {
	case proto.COM_QUIT:
		break
	case proto.COM_INIT_DB:
		break
	case proto.COM_PING:
		break
	case proto.COM_QUERY:
		fmt.Printf("Mysql com query. SQL=%v\n", string(removeSpaces.ReplaceAll(data, []byte{' '})))
	case proto.COM_STMT_PREPARE:
		fmt.Printf("Mysql com stmt perpare. %v\n", string(data))
	case proto.COM_STMT_EXECUTE:
		fmt.Printf("Mysql com stmt execute. %v\n", string(data))
	case proto.COM_STMT_SEND_LONG_DATA:
		fmt.Printf("Mysql com stmt send long data. %v\n", string(data))
	case proto.COM_STMT_RESET:
		fmt.Printf("Mysql com stmt reset. %v\n", string(data))
	case proto.COM_STMT_CLOSE:
		fmt.Printf("Mysql com stmt close. %v\n", string(data))
	default:
		fmt.Printf("Mysql command not implemented. %v\n%v\n", proto.CommandString(data[0]), string(data))
	}
	// STMT_CLOSE has not response
	if data[0] != proto.COM_STMT_CLOSE && data[0] != proto.COM_STMT_SEND_LONG_DATA {
		d.flyingRequests.PushBack(message)
	}

	return message, nil
}
