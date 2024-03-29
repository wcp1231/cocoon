package http

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/model/mock"
	"net/http"
	"strconv"
	"strings"
)

type HttpRequestMatcher struct {
	id     int32
	config mock.HttpMockConfig

	method    mock.FieldMatcher
	host      mock.FieldMatcher
	url       mock.FieldMatcher
	reqHeader map[string]mock.FieldMatcher

	status     int
	respHeader map[string]string
	respBody   string
}

func NewHttpRequestMatcherFromConfig(config mock.HttpMockConfig, id int32) *HttpRequestMatcher {
	matcher := &HttpRequestMatcher{
		id:     id,
		config: config,
	}
	if config.Request.Method != "" {
		matcher.method = &mock.StringMatcher{
			Expect: config.Request.Method,
		}
	}
	if config.Request.Host != nil {
		hostConfig := config.Request.Host
		matcher.host = mock.NewFieldMatcher(hostConfig)
	}
	if config.Request.Url != nil {
		urlConfig := config.Request.Url
		matcher.url = mock.NewFieldMatcher(urlConfig)
	}
	headerConfig := config.Request.Header
	if len(headerConfig) > 0 {
		matcher.reqHeader = map[string]mock.FieldMatcher{}
	}
	for k, v := range config.Request.Header {
		matcher.reqHeader[k] = mock.NewFieldMatcher(v)
	}

	matcher.status, _ = strconv.Atoi(config.Response.Status)
	matcher.respHeader = config.Response.Header
	matcher.respBody = config.Response.Body
	return matcher
}

func (h *HttpRequestMatcher) Match(r common.Message) bool {
	req := r.(*HTTPMessage)
	if h.method != nil {
		method := req.Meta["METHOD"]
		if !h.method.Match(method) {
			return false
		}
	}

	if h.host != nil {
		host := req.Meta["HOST"]
		if !h.host.Match(host) {
			return false
		}
	}

	if h.url != nil {
		url := req.Meta["URL"]
		if !h.url.Match(url) {
			return false
		}
	}

	for key, matcher := range h.reqHeader {
		val, exist := req.Header[key]
		if !exist || !matcher.Match(val) {
			return false
		}
	}

	return true
}

func (h *HttpRequestMatcher) Data() common.Message {
	resp := &HttpResponse{}
	resp.StatusCode = h.status
	resp.ProtoMajor = 1
	resp.ProtoMinor = 1
	resp.Header = http.Header{}
	//headers := make(map[string][]string)
	for k, v := range h.respHeader {
		//headers[k] = strings.Split(v, ";;")
		resp.Header[k] = strings.Split(v, ";;")
	}
	//message.SetHttpHeader(headers)

	resp.Body = []byte(h.respBody)
	message := NewHTTPGenericMessage()
	message.SetMock()
	message.SetResponse(resp)
	return message
}

func (h *HttpRequestMatcher) ID() int32 {
	return h.id
}

func (h *HttpRequestMatcher) GetConfig() interface{} {
	h.config.Id = h.id
	return h.config
}
