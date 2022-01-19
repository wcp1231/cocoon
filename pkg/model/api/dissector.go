package api

import (
	"cocoon/pkg/model/common"
	"errors"
	"fmt"
)

var (
	DissectorNotMatch = errors.New("dissector not match")
)

type Dissector interface {
	Name() string
	Dissect(reader TcpReader, isRequest bool) error
}

type DissectResult struct {
	Error          error
	ConnectionInfo *common.ConnectionInfo
	IsRequest      bool
	Protocol       *common.Protocol
	Payload        *common.GenericMessage
}

func (r *DissectResult) String() string {
	t := "->"
	if !r.IsRequest {
		t = "<-"
	}
	return fmt.Sprintf("%s %s %s %s", r.ConnectionInfo.String(), t, r.Protocol.String(), r.Payload.String())
}
