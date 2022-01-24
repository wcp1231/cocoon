package common

import (
	"fmt"
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
