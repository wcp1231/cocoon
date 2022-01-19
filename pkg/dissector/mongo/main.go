package mongo

import (
	"bytes"
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"encoding/json"
	"fmt"
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
	return "Mongo Dissector"
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

	headers, err := d.summary(req)
	if err != nil {
		return err
	}
	result := &api.DissectResult{}
	result.ConnectionInfo = connectionInfo
	result.IsRequest = isRequest
	result.Protocol = common.PROTOCOL_MONGO
	result.Payload = &common.GenericMessage{
		Header: headers,
		Raw:    &raw,
	}
	d.resultC <- result
	return nil
}

func (d *Dissector) summary(req RequestMsg) (map[string]string, error) {
	result := map[string]string{}
	switch r := req.(type) {
	case *Query:
		q, err := json.Marshal(r.Query)
		if err != nil {
			return nil, err
		}
		result["query"] = string(q)
	case *Update:
		q, err := json.Marshal(r.Update)
		if err != nil {
			return nil, err
		}
		result["update"] = string(q)
	case *Reply:
		result["size"] = fmt.Sprintf("%d", r.MessageLength)
	}
	return result, nil
}
