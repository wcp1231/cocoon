package redis

import (
	"bufio"
	"cocoon/pkg/model/common"
	"fmt"
	"time"
)

type Dissector struct {
	reqReader  *RESPReader
	respReader *RESPReader
	requestC   chan *common.GenericMessage
	responseC  chan *common.GenericMessage
}

func NewRequestDissector(reqC, respC chan *common.GenericMessage) *Dissector {
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
	message := &common.GenericMessage{}

	object, err := d.reqReader.ReadObject()
	if err != nil {
		// TODO response empty?
		fmt.Printf("Redis read request failed. %v\n", err)
		return err
	}

	request := object.GetRequest()
	message.CaptureTime = time.Now()
	message.Header = map[string]string{}
	message.Header["CMD"] = request.Cmd
	message.Header["KEY"] = request.Key
	body := object.Pretty()
	message.Body = &body
	message.Raw = &request.Raw

	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	message := &common.GenericMessage{}

	object, err := d.respReader.ReadObject()
	if err != nil {
		// TODO response empty?
		fmt.Printf("Redis read response failed. %v\n", err)
		return err
	}

	message.CaptureTime = time.Now()
	body := object.Pretty()
	message.Body = &body
	message.Raw = &object.Raw
	d.responseC <- message
	return nil
}
