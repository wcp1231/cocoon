package mysql

import (
	"bufio"
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/mysql/packet"
	"fmt"
)

type Dissector struct {
	reqReader  *bufio.Reader
	respReader *bufio.Reader
	requestC   chan *common.GenericMessage
	responseC  chan *common.GenericMessage

	reqStream  *packet.Stream
	respStream *packet.Stream
}

func NewRequestDissector(reqC, respC chan *common.GenericMessage) *Dissector {
	return &Dissector{
		requestC:  reqC,
		responseC: respC,
	}
}

func (d *Dissector) Init(inbound, outbound *bufio.Reader) {
	d.reqReader = inbound
	d.respReader = outbound
	d.reqStream = packet.NewStream(inbound, packet.PACKET_BUFFER_SIZE)
	d.respStream = packet.NewStream(outbound, packet.PACKET_BUFFER_SIZE)
}

func (d *Dissector) StartRequestDissect(reader *bufio.Reader) {
	//d.reqReader = reader
	for {
		err := d.dissectRequest()
		if err != nil {
			break
		}
	}
	fmt.Println("Mysql request dissect finish")
	close(d.requestC)
}

func (d *Dissector) StartResponseDissect(reader *bufio.Reader) {
	//d.respReader = reader
	for {
		err := d.dissectResponse()
		if err != nil {
			break
		}
	}
	fmt.Println("Mysql response dissect finish")
	close(d.responseC)
}

func (d *Dissector) ReadPacketFromServer() ([]byte, error) {
	pkt, err := d.respStream.NextPacket()
	if err != nil {
		return nil, err
	}
	return pkt.Raw(), nil
}

func (d *Dissector) ReadPacketFromClient() ([]byte, error) {
	pkt, err := d.reqStream.NextPacket()
	if err != nil {
		return nil, err
	}
	return pkt.Raw(), nil
}

func (d *Dissector) dissectRequest() error {
	//buf := new(bytes.Buffer)
	//br := io.TeeReader(d.reqReader, buf)

	message, err := d.readRequest()
	if err != nil {
		return err
	}

	message.CaptureNow()
	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	message, err := d.readResponse()
	if err != nil {
		return err
	}

	message.CaptureNow()
	d.responseC <- message
	return nil
}
