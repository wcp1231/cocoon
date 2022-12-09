package http

import (
	"bytes"
	"cocoon/pkg/model/common"
	"io"
	"net/http"
	"net/url"
)

const (
	HTTP_REQUEST_KEY  = "HTTP_REQUEST"
	HTTP_RESPONSE_KEY = "HTTP_RESPONSE"

	HTTP_STATUS_KEY = "HTTP_STATUS"
	HTTP_HEADER_KEY = "HTTP_HEADERS"
	HTTP_BODY_KEY   = "HTTP_BODY"
)

type HttpReuqest struct {
	Header http.Header
	Host   string
	Method string
	URL    string
	Body   []byte
}

func (r *HttpReuqest) raw() []byte {
	uri, _ := url.Parse(r.URL)
	request := &http.Request{
		Method: r.Method,
		Host:   r.Host,
		Header: r.Header,
		URL:    uri,
	}
	bodyBuf := bytes.NewBuffer(r.Body)
	request.Body = io.NopCloser(bodyBuf)
	buf := new(bytes.Buffer)
	_ = request.Write(buf)
	return buf.Bytes()
}

type HttpResponse struct {
	StatusCode int
	Proto      string
	ProtoMajor int
	ProtoMinor int
	Header     http.Header
	Body       []byte
}

func (r *HttpResponse) raw() []byte {
	response := &http.Response{
		StatusCode: r.StatusCode,
		Proto:      r.Proto,
		ProtoMajor: r.ProtoMajor,
		ProtoMinor: r.ProtoMinor,
		Header:     r.Header,
	}
	bodyBuf := bytes.NewBuffer(r.Body)
	response.ContentLength = int64(bodyBuf.Len())
	response.Body = io.NopCloser(bodyBuf)
	buf := new(bytes.Buffer)
	_ = response.Write(buf)
	return buf.Bytes()
}

type HTTPMessage struct {
	common.GenericMessage
}

func NewHTTPGenericMessage() *HTTPMessage {
	return &HTTPMessage{
		GenericMessage: common.NewGenericMessage(common.PROTOCOL_HTTP.Name),
	}
}

func (h *HTTPMessage) SetRequest(request *HttpReuqest) {
	h.Payload[HTTP_REQUEST_KEY] = request
}
func (h *HTTPMessage) SetResponse(response *HttpResponse) {
	if response.StatusCode >= 400 {
		h.MarkException()
	}
	h.Payload[HTTP_RESPONSE_KEY] = response
}

func (h *HTTPMessage) SetStatusCode(code int) {
	h.Payload[HTTP_STATUS_KEY] = code
}
func (h *HTTPMessage) GetStatusCode() int {
	return h.Payload[HTTP_STATUS_KEY].(int)
}
func (h *HTTPMessage) SetHttpHeader(header http.Header) {
	h.Payload[HTTP_HEADER_KEY] = header
}
func (h *HTTPMessage) GetHttpHeader() http.Header {
	return h.Payload[HTTP_HEADER_KEY].(http.Header)
}
func (h *HTTPMessage) SetBody(body []byte) {
	h.Payload[HTTP_BODY_KEY] = string(body)
	h.Body = body
}
func (h *HTTPMessage) GetBody() []byte {
	return h.Body
}
func (h *HTTPMessage) GetRaw() []byte {
	if h.Payload[HTTP_REQUEST_KEY] != nil {
		return h.Payload[HTTP_REQUEST_KEY].(*HttpReuqest).raw()
	}
	return h.Payload[HTTP_RESPONSE_KEY].(*HttpResponse).raw()
}
