package api

import (
	"bufio"
	"cocoon/pkg/model/common"
)

type ProtoClassifier interface {
	Match(r *bufio.Reader) bool
	Protocol() *common.Protocol
}

type ClientDissector interface {
	StartRequestDissect(reader *bufio.Reader)
	StartResponseDissect(reader *bufio.Reader)
}
