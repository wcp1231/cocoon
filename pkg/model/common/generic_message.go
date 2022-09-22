package common

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message interface {
	ID() int32
	SetId(int32)
	CaptureNow()
	GetCaptureTime() time.Time
	GetBody() *[]byte
	GetMeta() map[string]string
	GetHeader() map[string]string
	GetPayload() map[string]interface{}
	GetRaw() *[]byte
	SetRaw(*[]byte)
	SetMock()
	String() string
}

type GenericMessage struct {
	Id          int32
	CaptureTime time.Time `json:"captureTime"`
	Meta        map[string]string
	Header      map[string]string
	Payload     map[string]interface{}
	Body        *[]byte
	Raw         *[]byte `json:"raw"` // 原始数据 TODO 为啥要用指针？
}

func (g *GenericMessage) ID() int32 {
	return g.Id
}
func (g *GenericMessage) SetId(id int32) {
	g.Id = id
}

func (g *GenericMessage) CaptureNow() {
	g.CaptureTime = time.Now()
}
func (g *GenericMessage) GetCaptureTime() time.Time {
	return g.CaptureTime
}
func (g *GenericMessage) GetBody() *[]byte {
	return g.Body
}
func (g *GenericMessage) GetMeta() map[string]string {
	return g.Meta
}
func (g *GenericMessage) GetHeader() map[string]string {
	return g.Header
}
func (g *GenericMessage) GetPayload() map[string]interface{} {
	return g.Payload
}
func (g *GenericMessage) GetRaw() *[]byte {
	return g.Raw
}
func (g *GenericMessage) SetRaw(raw *[]byte) {
	g.Raw = raw
}

func (g *GenericMessage) SetMock() {
	g.Meta["MOCK"] = "true"
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

func NewGenericMessage(protocol string) GenericMessage {
	message := GenericMessage{
		Meta:    map[string]string{},
		Header:  map[string]string{},
		Payload: map[string]interface{}{},
	}
	message.Meta["PROTOCOL"] = protocol
	return message
}
