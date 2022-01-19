package rpc

import (
	"cocoon/pkg/model/common"
)

const COCOON_SERVER_NAME = "Cocoon"

type UploadReq struct {
	Session string
	Packet  *common.TcpPacket
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

type AnalysisReq struct {
	Session string
}

type AnalysisResp struct {
	Error error
}
