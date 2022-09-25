package dubbo

import (
	"bufio"
	"cocoon/pkg/model/common"
	"fmt"
)

type PackageType int

// enum part
const (
	PackageError              = PackageType(0x01)
	PackageRequest            = PackageType(0x02)
	PackageResponse           = PackageType(0x04)
	PackageHeartbeat          = PackageType(0x08)
	PackageRequest_TwoWay     = PackageType(0x10)
	PackageResponse_Exception = PackageType(0x20)
	PackageType_BitSize       = 0x2f
)

type Dissector struct {
	reqReader  *bufio.Reader
	respReader *bufio.Reader
	requestC   chan common.Message
	responseC  chan common.Message
}

type DubboHeader struct {
	SerialID       byte
	Type           PackageType
	ID             int64
	BodyLen        int
	ResponseStatus byte
}

type DubboRequest struct {
	Header         *DubboHeader
	DubboVersion   string
	Target         string
	ServiceVersion string
	Method         string
	Args           map[string]interface{}
	Attachments    map[string]interface{}
}

type DubboResponse struct {
	Header       *DubboHeader
	DubboVersion string
	Exception    string
	RespObj      interface{}
	Attachments  map[string]interface{}
}

func (d *DubboHeader) isRequest() bool {
	return d.Type&PackageRequest != 0x00
}

func (d *DubboHeader) isResponse() bool {
	return d.Type&PackageResponse != 0x00
}

func (d *DubboHeader) isHeartbeat() bool {
	return d.Type&PackageHeartbeat != 0x00
}

func (d *DubboHeader) hasException() bool {
	return d.Type&PackageResponse_Exception != 0x00
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
	fmt.Println("Dubbo request dissect stop")
}

func (d *Dissector) StartResponseDissect(reader *bufio.Reader) {
	d.respReader = reader
	for {
		err := d.dissectResponse()
		if err != nil {
			break
		}
	}
	fmt.Println("Dubbo response dissect stop")
}

func (d *Dissector) dissectRequest() error {
	message, err := ReadPacket(d.reqReader)
	if err != nil {
		// TODO response empty?
		fmt.Printf("Dubbo read request failed. %v\n", err)
		return err
	}
	message.CaptureNow()
	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	message, err := ReadPacket(d.respReader)
	if err != nil {
		// TODO response empty?
		fmt.Printf("Dubbo read response failed. %v\n", err)
		return err
	}
	message.CaptureNow()
	d.responseC <- message
	return nil
}
