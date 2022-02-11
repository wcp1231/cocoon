package mysql

import (
	"cocoon/pkg/model/common"
	"fmt"
)

func (d *Dissector) readRequest() (*common.GenericMessage, error) {
	pkt, err := d.reqStream.NextPacket()
	if err != nil {
		return nil, err
	}
	raw := pkt.Raw()
	data := pkt.Datas
	message := common.NewMysqlGenericMessage()
	message.Raw = &raw
	message.Meta["OP_TYPE"] = CommandString(data[0])

	switch data[0] {
	case COM_QUIT:
		break
	case COM_INIT_DB:
		break
	case COM_PING:
		break
	case COM_QUERY:
		fmt.Printf("Mysql com query. %v\n", string(data))
	case COM_STMT_PREPARE:
		fmt.Printf("Mysql com stmt perpare. %v\n", string(data))
	case COM_STMT_EXECUTE:
		fmt.Printf("Mysql com stmt execute. %v\n", string(data))
	case COM_STMT_RESET:
		fmt.Printf("Mysql com stmt reset. %v\n", string(data))
	case COM_STMT_CLOSE:
		fmt.Printf("Mysql com stmt close. %v\n", string(data))
	default:
		fmt.Printf("Mysql command not implemented. %v\n%v\n", CommandString(data[0]), string(data))
	}
	return message, nil
}
