package mock

import (
	"bytes"
	"cocoon/pkg/model/common"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type HttpRequestMatcher struct {
	id     int32
	config httpMockConfig

	method    FieldMatcher
	host      FieldMatcher
	url       FieldMatcher
	reqHeader map[string]FieldMatcher

	status     string
	respHeader map[string]string
	respBody   string
}

func newHttpRequestMatcherFromConfig(config httpMockConfig, id int32) *HttpRequestMatcher {
	matcher := &HttpRequestMatcher{
		id:     id,
		config: config,
	}
	if config.Request.Method != "" {
		matcher.method = &StringMatcher{
			expect: config.Request.Method,
		}
	}
	if config.Request.Host != nil {
		hostConfig := config.Request.Host
		matcher.host = newFieldMatcher(hostConfig)
	}
	if config.Request.Url != nil {
		urlConfig := config.Request.Url
		matcher.url = newFieldMatcher(urlConfig)
	}
	headerConfig := config.Request.Header
	if len(headerConfig) > 0 {
		matcher.reqHeader = map[string]FieldMatcher{}
	}
	for k, v := range config.Request.Header {
		matcher.reqHeader[k] = newFieldMatcher(v)
	}

	matcher.status = config.Response.Status
	matcher.respHeader = config.Response.Header
	matcher.respBody = config.Response.Body
	return matcher
}

func newFieldMatcher(field *fieldMockConfig) FieldMatcher {
	if field.Equals != "" {
		return &StringMatcher{
			expect: field.Equals,
		}
	}
	regex, err := regexp.Compile(field.Regex)
	if err != nil {
		// TODO
		os.Exit(1)
	}
	return &RegexMatcher{
		regex: regex,
	}
}

func (h *HttpRequestMatcher) Match(r common.Message) bool {
	req := r.(*common.HTTPMessage)
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
	response := http.Response{}
	response.StatusCode, _ = strconv.Atoi(h.status)
	response.Header = http.Header{}
	response.ProtoMajor = 1
	response.ProtoMinor = 1

	message := common.NewHTTPGenericMessage()
	message.SetMock()
	message.Meta["STATUS"] = h.status
	for k, v := range h.respHeader {
		message.Header[k] = v
		response.Header[k] = strings.Split(v, ";;")
	}
	body := []byte(h.respBody)
	message.Body = &body

	bodyBuf := bytes.NewBuffer(body)
	response.ContentLength = int64(bodyBuf.Len())
	response.Body = io.NopCloser(bodyBuf)
	buf := new(bytes.Buffer)
	_ = response.Write(buf)
	bs := buf.Bytes()
	message.Raw = &bs
	return message
}

func (h *HttpRequestMatcher) ID() int32 {
	return h.id
}

func (h *HttpRequestMatcher) GetConfig() interface{} {
	h.config.Id = h.id
	return h.config
}
