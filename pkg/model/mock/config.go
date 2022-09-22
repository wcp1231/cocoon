package mock

type MockConfig struct {
	Http  []HttpMockConfig  `json:"http"`
	Redis []RedisMockConfig `json:"redis"`
}

type FieldMockConfig struct {
	Equals string `json:"equals,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

type HttpRequestMockConfig struct {
	Method string                      `json:"method,omitempty"`
	Host   *FieldMockConfig            `json:"host,omitempty"`
	Url    *FieldMockConfig            `json:"url,omitempty"`
	Header map[string]*FieldMockConfig `json:"header,omitempty"`
}

type HttpResponseMockConfig struct {
	Status string            `json:"status"`
	Header map[string]string `json:"header,omitempty"`
	Body   string            `json:"body"`
}

type HttpMockConfig struct {
	Id       int32                  `json:"id,omitempty"`
	Request  HttpRequestMockConfig  `json:"request"`
	Response HttpResponseMockConfig `json:"response"`
}

type RedisRequestMockConfig struct {
	Cmd  *FieldMockConfig  `json:"cmd,omitempty"`
	Key  *FieldMockConfig  `json:"key,omitempty"`
	Args []FieldMockConfig `json:"args,omitempty"`
}

type RedisResponseObject struct {
	Type  string                `json:"type,omitempty"`
	Value string                `json:"value,omitempty"`
	Array []RedisResponseObject `json:"array,omitempty"`
}

type RedisMockConfig struct {
	Id       int32                  `json:"id,omitempty"`
	Request  RedisRequestMockConfig `json:"request"`
	Response RedisResponseObject    `json:"response"`
}
