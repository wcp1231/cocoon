package common

import (
	"encoding/json"
	"fmt"
	"time"
)

type GenericMessage struct {
	Id          int32
	CaptureTime time.Time `json:"captureTime"`
	Meta        map[string]string
	Header      map[string]string
	Body        *[]byte
	Raw         *[]byte `json:"raw"` // 原始数据
}

func NewHTTPGenericMessage() *GenericMessage {
	return NewGenericMessage(PROTOCOL_HTTP.Name)
}

func NewDubboGenericMessage() *GenericMessage {
	return NewGenericMessage(PROTOCOL_DUBBO.Name)
}

func NewRedisGenericMessage() *GenericMessage {
	return NewGenericMessage(PROTOCOL_REDIS.Name)
}

func NewMongoGenericMessage() *GenericMessage {
	return NewGenericMessage(PROTOCOL_MONGO.Name)
}

func NewMysqlGenericMessage() *GenericMessage {
	return NewGenericMessage(PROTOCOL_MYSQL.Name)
}

func NewGenericMessage(protocol string) *GenericMessage {
	message := &GenericMessage{
		Meta:   map[string]string{},
		Header: map[string]string{},
	}
	message.Meta["PROTOCOL"] = protocol
	return message
}

func (g *GenericMessage) CaptureNow() {
	g.CaptureTime = time.Now()
}

func (g *GenericMessage) String() string {
	return fmt.Sprintf("[%d] Header=%d", g.Id, len(g.Header))
}

func (g *GenericMessage) ToJSON() ([]byte, error) {
	result := map[string]interface{}{
		"id":          g.Id,
		"captureTime": g.CaptureTime,
		"header":      g.Header,
		"body":        g.Body,
	}
	return json.Marshal(result)
}
