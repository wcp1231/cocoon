package redis

import (
	"bufio"
	"cocoon/pkg/model/common"
	"fmt"
)

type Dissector struct {
	reqReader  *RESPReader
	respReader *RESPReader
	requestC   chan common.Message
	responseC  chan common.Message
}

func NewRequestDissector(reqC, respC chan common.Message) *Dissector {
	return &Dissector{
		requestC:  reqC,
		responseC: respC,
	}
}

func (d *Dissector) StartRequestDissect(reader *bufio.Reader) {
	d.reqReader = NewReader(reader)
	for {
		err := d.dissectRequest()
		if err != nil {
			break
		}
	}
}

func (d *Dissector) StartResponseDissect(reader *bufio.Reader) {
	d.respReader = NewReader(reader)
	for {
		err := d.dissectResponse()
		if err != nil {
			break
		}
	}
}

func (d *Dissector) dissectRequest() error {
	message := NewRedisGenericMessage()

	object, err := d.reqReader.ReadObject()
	if err != nil {
		// TODO response empty?
		fmt.Printf("Redis read request failed. %v\n", err)
		return err
	}

	message.SetRequestCmd(object.Pretty())
	reqCmds := object.(*RedisArray)
	message.SetCmd(reqCmds.Items[0].Pretty())
	if reqCmds.Len > 1 {
		message.SetKey(reqCmds.Items[1].Pretty())
	}
	if message.GetCmd() == `"PING"` {
		message.SetHeartbeat()
	}

	raw := object.Raw()
	message.SetRaw(&raw)
	message.CaptureNow()
	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	message := NewRedisGenericMessage()

	object, err := d.respReader.ReadObject()
	if err != nil {
		// TODO response empty?
		fmt.Printf("Redis read response failed. %v\n", err)
		return err
	}

	message.SetResponseObj(object.Pretty())
	raw := object.Raw()
	message.SetRaw(&raw)
	message.CaptureNow()
	d.responseC <- message
	return nil
}
