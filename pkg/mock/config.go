package mock

type mockConfig struct {
	Http []httpMockConfig `json:"http"`
	Redis []redisMockConfig `json:"redis"`
}

type fieldMockConfig struct {
	Match string `json:"match,omitempty"`
	Regex string `json:"regex,omitempty"`
}

type httpRequestMockConfig struct {
	Method string `json:"method,omitempty"`
	Host *fieldMockConfig `json:"host,omitempty"`
	Url *fieldMockConfig `json:"url,omitempty"`
	Header map[string]*fieldMockConfig `json:"header,omitempty"`
}

type httpResponseMockConfig struct {
	Status string `json:"status"`
	Header map[string]string `json:"header,omitempty"`
	Body string `json:"body"`
}

type httpMockConfig struct {
	Id int32 `json:"id,omitempty"`
	Request httpRequestMockConfig `json:"request"`
	Response httpResponseMockConfig `json:"response"`
}

type redisRequestMockConfig struct {
	Cmd *fieldMockConfig `json:"cmd,omitempty"`
	Key *fieldMockConfig `json:"key,omitempty"`
}

type redisResponseMockConfig struct {
	Type string `json:"type,omitempty"`
	String string `json:"string,omitempty"`
	Array []string `json:"array,omitempty"`
	Hash map[string]string `json:"hash,omitempty"`
	Zset interface{} `json:"zset,omitempty"`
	Err string `json:"err,omitempty"`
}

type redisMockConfig struct {
	Id int32 `json:"id,omitempty"`
	Request redisRequestMockConfig `json:"request"`
	Response redisResponseMockConfig `json:"response"`
}