package traffic

import (
	"cocoon/pkg/model/common"
	"fmt"
)

type RequestResponsePair struct {
	Request  *common.GenericMessage `json:"request"`
	Response *common.GenericMessage `json:"response"`
}

func (r *RequestResponsePair) String() string {
	req := ""
	resp := ""
	if r.Request != nil {
		req = r.Request.String()
	}
	if r.Response != nil {
		resp = r.Response.String()
	}
	return fmt.Sprintf("Req=[%s] Resp=[%s]", req, resp)
}

type StreamItem struct {
	ID             string           `json:"ID"` // src_ip:src_port|dst_ip:dst_port
	Timestamp      int64            `json:"timestamp"`
	Protocol       *common.Protocol `json:"protocol"`
	ConnectionInfo *common.ConnectionInfo
	Pair           *RequestResponsePair
}

func (s *StreamItem) String() string {
	return fmt.Sprintf("%s %s %s", s.ID, s.Protocol.String(), s.Pair.String())
}
