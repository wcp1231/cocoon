package proto

import (
	"bufio"
	"cocoon/pkg/model/common"
)

type DefaultClassifier struct{}

func (c *DefaultClassifier) Match(r *bufio.Reader) bool {
	return true
}

func (c *DefaultClassifier) Protocol() *common.Protocol {
	return common.PROTOCOL_NOT_SUPPORTED
}
