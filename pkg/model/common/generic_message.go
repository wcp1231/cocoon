package common

import (
	"fmt"
	"time"
)

type GenericMessage struct {
	CaptureTime time.Time `json:"captureTime"`
	Header      map[string]string
	Body        *[]byte
	Raw         *[]byte `json:"raw"` // 原始数据
}

func (g *GenericMessage) String() string {
	return fmt.Sprintf("Header=%d", len(g.Header))
}
