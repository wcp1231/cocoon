package kafka

import (
	"bytes"
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
)

type Dissector struct {
	resultC chan *api.DissectResult
}

func NewDissector(resultC chan *api.DissectResult) *Dissector {
	return &Dissector{
		resultC: resultC,
	}
}

func (d *Dissector) Name() string {
	return "Kafka Dissector"
}

func (d *Dissector) Dissect(reader api.TcpReader, isRequest bool) error {
	connectionInfo := reader.Connection()
	return d.handle(connectionInfo, reader, isRequest)
}

func (d *Dissector) handle(connectionInfo *common.ConnectionInfo, reader api.TcpReader, isRequest bool) error {
	br := reader.BufferReader()
	req, err := ReadRequest(br)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	WriteRequest(req, buf)
	raw := buf.Bytes()

	//headers, err := d.summary(req)
	//if err != nil {
	//	return err
	//}
	headers := map[string]string{}
	result := &api.DissectResult{}
	result.ConnectionInfo = connectionInfo
	result.IsRequest = isRequest
	result.Protocol = common.PROTOCOL_KAFKA
	result.Payload = &common.GenericMessage{
		Header: headers,
		Raw:    &raw,
	}
	d.resultC <- result
	return nil
}
