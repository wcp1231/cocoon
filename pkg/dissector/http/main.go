package http

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"net/http"
	"strings"
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
	return "HTTP Dissector"
}

func (d *Dissector) Dissect(reader api.TcpReader, isRequest bool) error {
	connectionInfo := reader.Connection()
	br := reader.BufferReader()
	if isRequest {
		return d.handleRequest(connectionInfo, br)
	} else {
		return d.handleResponse(connectionInfo, br)
	}
}

func (d *Dissector) handleRequest(connectionInfo *common.ConnectionInfo, br *bufio.Reader) error {
	request, err := http.ReadRequest(br)
	if err != nil {
		return err
	}

	headers := d.parseHeaders(request.Header)
	raw := d.parseRawRequest(request)

	result := &api.DissectResult{}
	result.ConnectionInfo = connectionInfo
	result.IsRequest = true
	result.Protocol = common.PROTOCOL_HTTP
	result.Payload = &common.GenericMessage{
		Header: headers,
		Raw:    &raw,
	}
	d.resultC <- result
	return nil
}

func (d *Dissector) handleResponse(connectionInfo *common.ConnectionInfo, br *bufio.Reader) error {
	response, err := http.ReadResponse(br, nil)
	if err != nil {
		return err
	}

	headers := d.parseHeaders(response.Header)
	raw := d.parseRawResponse(response)

	result := &api.DissectResult{}
	result.ConnectionInfo = connectionInfo
	result.IsRequest = false
	result.Protocol = common.PROTOCOL_HTTP
	result.Payload = &common.GenericMessage{
		Header: headers,
		Raw:    &raw,
	}
	d.resultC <- result
	return nil
}

func (d *Dissector) parseHeaders(headers http.Header) map[string]string {
	ret := make(map[string]string)
	for k, v := range headers {
		ret[k] = strings.Join(v, ";")
	}
	return ret
}

func (d *Dissector) parseRawRequest(r *http.Request) []byte {
	buf := new(bytes.Buffer)
	r.Write(buf)
	return buf.Bytes()
}

func (d *Dissector) parseRawResponse(r *http.Response) []byte {
	buf := new(bytes.Buffer)
	r.Write(buf)
	return buf.Bytes()
}
