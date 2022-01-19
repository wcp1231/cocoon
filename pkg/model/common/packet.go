package common

import (
	"time"
)

type TcpPacket struct {
	Source      string
	Destination string
	IsOutgoing  bool
	Direction   *Direction
	Seq         uint64
	Timestamp   time.Time
	Payload     []byte
}

func (t *TcpPacket) IsRequest() bool {
	isRequest := *t.Direction == ClientToRemote
	if !t.IsOutgoing {
		isRequest = *t.Direction == RemoteToClient
	}
	return isRequest
}
