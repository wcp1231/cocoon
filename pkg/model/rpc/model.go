package rpc

import (
	"cocoon/pkg/model"
	"time"
)

const COCOON_SERVER_NAME = "Cocoon"

type TcpPacket struct {
	Source      string
	Destination string
	IsOutgoing  bool
	Direction   *model.Direction
	Seq         uint64
	Timestamp   time.Time
	Payload     []byte
}

func (t *TcpPacket) IsRequest() bool {
	isRequest := *t.Direction == model.ClientToRemote
	if !t.IsOutgoing {
		isRequest = *t.Direction == model.RemoteToClient
	}
	return isRequest
}

type UploadReq struct {
	Session string
	Packet  *TcpPacket
}

type UploadResp struct {
	Count int
}

type PostStartReq struct {
	Appname string
	Session string
}

type PostStartResp struct {
	Error error
}
