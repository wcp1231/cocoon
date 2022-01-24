package rpc

import (
	"cocoon/pkg/model/common"
)

const COCOON_SERVER_NAME = "Cocoon"
const (
	OP_MOCK = 1
	OP_PASS = 2
)

type UploadReq struct {
	Session string
	Packet  *common.TcpPacket
}

type UploadResp struct {
	Error error
}

type ConnCloseReq struct {
	Source      string
	Destination string
	Direction   *common.Direction
}

type ConnCloseResp struct {
	Error error
}

type PostStartReq struct {
	Appname string
	Session string
}

type PostStartResp struct {
	Error error
}

type AnalysisReq struct {
	Session string
}

type AnalysisResp struct {
	Error error
}

type OutboundReq struct {
	Session string
	Proto   *common.Protocol
	Request *common.GenericMessage
}

type OutboundResp struct {
	OpType int
	Data   *[]byte
}

type RecordReq struct {
	Session    string
	IsOutgoing bool
	Proto      *common.Protocol
	ReqHeader  map[string]string
	RespHeader map[string]string
	ReqBody    *[]byte
	RespBody   *[]byte
}

type RecordResp struct {
	Error error
}
