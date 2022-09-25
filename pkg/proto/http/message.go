package http

import (
	"cocoon/pkg/model/common"
	"net/http"
	"strconv"
	"strings"
)

const (
	HTTP_STATUS_KEY = "HTTP_STATUS"
	HTTP_HOST_KEY   = "HTTP_HOST"
	HTTP_METHOD_KEY = "HTTP_METHOD"
	HTTP_URL_KEY    = "HTTP_URL"
	HTTP_PROTO_KEY  = "HTTP_PROTO"
	HTTP_HEADER_KEY = "HTTP_HEADERS"
	HTTP_BODY_KEY   = "HTTP_BODY"
)

type HTTPMessage struct {
	common.GenericMessage

	raw []byte
}

func NewHTTPGenericMessage() *HTTPMessage {
	return &HTTPMessage{
		GenericMessage: common.NewGenericMessage(common.PROTOCOL_HTTP.Name),
	}
}

func (h *HTTPMessage) SetStatusCode(code int) {
	h.Payload[HTTP_STATUS_KEY] = code
	h.Meta["STATUS"] = strconv.Itoa(code)
}
func (h *HTTPMessage) GetStatusCode() int {
	return h.Payload[HTTP_STATUS_KEY].(int)
}
func (h *HTTPMessage) SetHost(host string) {
	h.Payload[HTTP_HOST_KEY] = host
	h.Meta["HOST"] = host
}
func (h *HTTPMessage) GetHost() string {
	return h.Payload[HTTP_HOST_KEY].(string)
}
func (h *HTTPMessage) SetMethod(method string) {
	h.Payload[HTTP_METHOD_KEY] = method
	h.Meta["METHOD"] = method
}
func (h *HTTPMessage) GetMethod() string {
	return h.Payload[HTTP_METHOD_KEY].(string)
}
func (h *HTTPMessage) SetUrl(url string) {
	h.Payload[HTTP_URL_KEY] = url
	h.Meta["URL"] = url // 过渡
}
func (h *HTTPMessage) GetUrl() string {
	return h.Payload[HTTP_URL_KEY].(string)
}
func (h *HTTPMessage) SetProto(proto string) {
	h.Payload[HTTP_PROTO_KEY] = proto
	h.Meta["PROTO"] = proto
}
func (h *HTTPMessage) GetProto() string {
	return h.Payload[HTTP_PROTO_KEY].(string)
}
func (h *HTTPMessage) SetHttpHeader(header http.Header) {
	h.Payload[HTTP_HEADER_KEY] = header
	for k, vv := range header {
		h.Header[k] = strings.Join(vv, ";")
	}
}
func (h *HTTPMessage) GetHttpHeader() http.Header {
	return h.Payload[HTTP_HEADER_KEY].(http.Header)
}
func (h *HTTPMessage) SetBody(body []byte) {
	h.Payload[HTTP_BODY_KEY] = string(body)
	h.Body = &body
}
func (h *HTTPMessage) GetBody() *[]byte {
	return h.Body
}
func (h *HTTPMessage) setRaw(raw []byte) {
	h.raw = raw
}
func (h *HTTPMessage) GetRaw() []byte {
	return h.raw
}
