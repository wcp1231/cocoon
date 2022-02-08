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
	requestC   chan *common.GenericMessage
	responseC  chan *common.GenericMessage
}

type dubboHeader struct {
	SerialID       byte
	Type           PackageType
	ID             int64
	BodyLen        int
	ResponseStatus byte
}

type dubboRequest struct {
	header         *dubboHeader
	dubboVersion   string
	target         string
	serviceVersion string
	method         string
	args           map[string]interface{}
	attachments    map[string]interface{}
}

type dubboResponse struct {
	header       *dubboHeader
	dubboVersion string
	exception    string
	respObj      interface{}
	attachments  map[string]interface{}
}

func (d *dubboHeader) isRequest() bool {
	return d.Type&PackageRequest != 0x00
}

func (d *dubboHeader) hasException() bool {
	return d.Type&PackageResponse_Exception != 0x00
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
