package mongo

import (
	"bufio"
	"cocoon/pkg/model/common"
	"fmt"
)

type Dissector struct {
	reqReader  *bufio.Reader
	respReader *bufio.Reader
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
	d.reqReader = reader
	for {
		err := d.dissectRequest()
		if err != nil {
			break
		}
	}
	fmt.Println("Mongo request dissect finish")
	close(d.requestC)
}

func (d *Dissector) StartResponseDissect(reader *bufio.Reader) {
	d.respReader = reader
	for {
		err := d.dissectResponse()
		if err != nil {
			break
		}
	}
	fmt.Println("Mongo response dissect finish")
	close(d.responseC)
}

func (d *Dissector) dissectRequest() error {
	message, err := Parse(d.reqReader)
	if err != nil {
		fmt.Println("Mongo request dissect error", err.Error())
		return err
	}

	message.CaptureNow()
	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	message, err := Parse(d.respReader)
	if err != nil {
		fmt.Println("Mongo response dissect error", err.Error())
		return err
	}

	message.CaptureNow()
	d.responseC <- message
	return nil
}
