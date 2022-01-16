package http

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"fmt"
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

func (d *Dissector) Match(reader api.TcpReader) bool {
	br := reader.BufferReader()
	bytes, err := br.Peek(10)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	text := string(bytes)
	isHttp := strings.HasPrefix(text, "GET /")
	isHttp = isHttp || strings.HasPrefix(text, "HEAD /")
	isHttp = isHttp || strings.HasPrefix(text, "POST /")
	isHttp = isHttp || strings.HasPrefix(text, "PUT /")
	isHttp = isHttp || strings.HasPrefix(text, "PATCH /")
	isHttp = isHttp || strings.HasPrefix(text, "DELETE /")
	isHttp = isHttp || strings.HasPrefix(text, "OPTIONS /")
	isHttp = isHttp || strings.HasPrefix(text, "TRACE /")
	//isHttp = isHttp || strings.HasPrefix(text, "CONNECT /")

	// response
	isHttp = isHttp || strings.HasPrefix(text, "HTTP/1.1")
	return isHttp
}

func (d *Dissector) Dissect(reader api.TcpReader, isRequest bool) {
	connectionInfo := reader.Connection()
	br := reader.BufferReader()
	if isRequest {
		d.handleRequest(connectionInfo, br)
	} else {
		d.handleResponse(connectionInfo, br)
	}
}

func (d *Dissector) handleRequest(connectionInfo *common.ConnectionInfo, br *bufio.Reader) {
	request, err := http.ReadRequest(br)
	if err != nil {
		fmt.Println(err.Error())
		return
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
}

func (d *Dissector) handleResponse(connectionInfo *common.ConnectionInfo, br *bufio.Reader) {
	response, err := http.ReadResponse(br, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
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
}

func (d *Dissector) parseHeaders(headers http.Header) map[string]string {
	ret := make(map[string]string)
	for k, v := range headers {
		ret[k] = strings.Join(v, ",")
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
