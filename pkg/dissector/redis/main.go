package redis

import (
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"fmt"
)

type Dissector struct {
	resultC chan *api.DissectResult
	reader  *RESPReader
}

func NewDissector(resultC chan *api.DissectResult) *Dissector {
	return &Dissector{
		resultC: resultC,
	}
}

func (d *Dissector) Name() string {
	return "Redis Dissector"
}

func (d *Dissector) Dissect(reader api.TcpReader, isRequest bool) error {
	connectionInfo := reader.Connection()
	d.reader = NewReader(reader)
	return d.handleRedis(connectionInfo, isRequest)
}

func (d *Dissector) handleRedis(connectionInfo *common.ConnectionInfo, isRequest bool) error {
	object, err := d.reader.ReadObject()
	if err != nil {
		return err
	}
	if object == nil {
		return nil
	}

	// TODO
	headers := map[string]string{}
	if isRequest {
		if object.Type == ARRAY {
			headers["0"] = object.Array[0].Data
			if object.Count > 1 {
				headers["1"] = object.Array[1].Data
			}
		}
	} else {
		headers["size"] = fmt.Sprintf("%d", len(object.Data))
	}

	result := &api.DissectResult{}
	result.ConnectionInfo = connectionInfo
	result.IsRequest = isRequest
	result.Protocol = common.PROTOCOL_REDIS
	result.Payload = &common.GenericMessage{
		Header: headers,
		Raw:    &object.Raw,
	}
	d.resultC <- result
	return nil
}
