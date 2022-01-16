package common

import (
	"fmt"
	"time"
)

type ConnectionInfo struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	IsOutgoing  bool   `json:"isOutgoing"` // TCP 是向外发送的，还是接收外部的
}

func (c *ConnectionInfo) ID() string {
	return c.Source + "|" + c.Destination
}

func (c *ConnectionInfo) String() string {
	t := "IN"
	if c.IsOutgoing {
		t = "OUT"
	}
	return fmt.Sprintf("[%s][%s|%s]", t, c.Source, c.Destination)
}

// TODO 单独文件
type GenericMessage struct {
	CaptureTime time.Time         `json:"captureTime"`
	Raw         *[]byte           `json:"raw"` // 原始数据
	Header      map[string]string // TODO value 的类型
	Body        *[]byte
}

func (g *GenericMessage) String() string {
	return fmt.Sprintf("Header=%d, Size=%d", len(g.Header), len(*g.Raw))
}
