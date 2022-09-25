package http

import (
	"testing"
)

func TestRequestImporter_Match(t *testing.T) {
	//req := &common.GenericMessage{
	//	Header: map[string]string{
	//		"METHOD":          "GET",
	//		"HOST":            "box.http.svc.dev.keep",
	//		"URL":             "/internal/user/61f170000000000000000000",
	//		"Accept":          "*/*",
	//		"Accept-Encoding": "gzip",
	//		"Connection":      "Keep-Alive",
	//		"User-Agent":      "okhttp/3.14.2",
	//	},
	//}
	//config := mock.HttpMockConfig{
	//	Request: mock.HttpRequestMockConfig{
	//		Url: &mock.FieldMockConfig{
	//			Equals: "/internal/user/61f170000000000000000000",
	//		},
	//	},
	//	Response: mock.HttpResponseMockConfig{
	//		Status: "200",
	//		Header: map[string]string{
	//			"Content-Type": "application/json;charset=UTF-8",
	//		},
	//		Body: "{}",
	//	},
	//}
	//matcher := NewHttpRequestMatcherFromConfig(config, 0)
	//
	//if !matcher.Match(req) {
	//	t.Fatal("HttpRequestMatcher failed")
	//}
}
