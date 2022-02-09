package mongo

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/common"
	"fmt"
	"io"
)

type Dissector struct {
	reqReader  *bufio.Reader
	respReader *bufio.Reader
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
	buf := new(bytes.Buffer)
	br := io.TeeReader(d.reqReader, buf)

	message, err := Parse(br)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		fmt.Println("Mongo request dissect error", err.Error())
		return err
	}

	raw := buf.Bytes()
	message.Raw = &raw

	message.CaptureNow()
	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	buf := new(bytes.Buffer)
	br := io.TeeReader(d.respReader, buf)

	message, err := Parse(br)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		fmt.Println("Mongo response dissect error", err.Error())
		return err
	}

	raw := buf.Bytes()
	message.Raw = &raw
	message.CaptureNow()
	d.responseC <- message
	return nil
}
