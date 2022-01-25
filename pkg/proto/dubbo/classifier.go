package dubbo

import (
	"bufio"
	"cocoon/pkg/model/common"
)

type Classifier struct{}

func (c *Classifier) Match(r *bufio.Reader) bool {
	buf, err := r.Peek(16)
	if err != nil {
		return false
	}

	if buf[0] != MAGIC_HIGH && buf[1] != MAGIC_LOW {
		return false
	}

	serialId := buf[2] & SERIAL_MASK
	if serialId == Zero {
		return false
	}

	return true
}

func (c *Classifier) Protocol() *common.Protocol {
	return common.PROTOCOL_DUBBO
}
