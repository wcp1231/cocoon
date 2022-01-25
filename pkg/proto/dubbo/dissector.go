package dubbo

import (
	"bufio"
	"cocoon/pkg/model/common"
	"encoding/binary"
	"fmt"
	"time"
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

	d.responseC <- message
	return nil
}

func ReadPacket(reader *bufio.Reader) (*common.GenericMessage, error) {
	header, err := ReadHeader(reader)
	if err != nil {
		return nil, err
	}

	body := make([]byte, header.BodyLen)
	_, err = reader.Read(body)
	if err != nil {
		return nil, err
	}

	headerBytes := EncodeHeader(header)
	raw := make([]byte, len(headerBytes)+len(body))
	copy(raw, headerBytes)
	copy(raw[len(headerBytes):], body)

	message := &common.GenericMessage{}
	message.CaptureTime = time.Now()
	message.Body = &body // FIXME
	message.Raw = &raw
	return message, nil
}

func ReadHeader(reader *bufio.Reader) (*dubboHeader, error) {
	header := &dubboHeader{}
	var err error
	buf, err := reader.Peek(HEADER_LENGTH)
	if err != nil { // this is impossible
		return nil, err
	}
	_, err = reader.Discard(HEADER_LENGTH)
	if err != nil { // this is impossible
		return nil, err
	}

	//// read header
	if buf[0] != MAGIC_HIGH && buf[1] != MAGIC_LOW {
		return nil, ErrIllegalPackage
	}

	// Header{serialization id(5 bit), event, two way, req/response}
	if header.SerialID = buf[2] & SERIAL_MASK; header.SerialID == Zero {
		return nil, fmt.Errorf("serialization ID:%v", header.SerialID)
	}

	flag := buf[2] & FLAG_EVENT
	if flag != Zero {
		header.Type |= PackageHeartbeat
	}
	flag = buf[2] & FLAG_REQUEST
	if flag != Zero {
		header.Type |= PackageRequest
		flag = buf[2] & FLAG_TWOWAY
		if flag != Zero {
			header.Type |= PackageRequest_TwoWay
		}
	} else {
		header.Type |= PackageResponse
		header.ResponseStatus = buf[3]
		if header.ResponseStatus != Response_OK {
			header.Type |= PackageResponse_Exception
		}
	}

	// Header{req id}
	header.ID = int64(binary.BigEndian.Uint64(buf[4:]))

	// Header{body len}
	header.BodyLen = int(binary.BigEndian.Uint32(buf[12:]))
	if header.BodyLen < 0 {
		return nil, ErrIllegalPackage
	}

	return header, err
}

func EncodeHeader(header *dubboHeader) []byte {
	bs := make([]byte, 0)
	switch {
	case header.Type&PackageHeartbeat != 0x00:
		if header.ResponseStatus == Zero {
			bs = append(bs, DubboRequestHeartbeatHeader[:]...)
		} else {
			bs = append(bs, DubboResponseHeartbeatHeader[:]...)
		}
	case header.Type&PackageResponse != 0x00:
		bs = append(bs, DubboResponseHeaderBytes[:]...)
		if header.ResponseStatus != 0 {
			bs[3] = header.ResponseStatus
		}
	case header.Type&PackageRequest_TwoWay != 0x00:
		bs = append(bs, DubboRequestHeaderBytesTwoWay[:]...)
	}
	bs[2] |= header.SerialID & SERIAL_MASK
	binary.BigEndian.PutUint64(bs[4:], uint64(header.ID))
	binary.BigEndian.PutUint32(bs[12:], uint32(header.BodyLen))
	return bs
}
