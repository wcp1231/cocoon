package mongo

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/common"
	"fmt"
)

type Classifier struct{}

func (c *Classifier) Match(r *bufio.Reader) bool {
	head, err := r.Peek(16)
	if err != nil {
		return false
	}

	br := bytes.NewBuffer(head)
	header, err := readMsgHeader(br)
	if err != nil {
		fmt.Printf("Mongo classifier read header failed. %v\n", err)
		return false
	}
	fmt.Printf("Mongo classifier read header %v\n", header)
	if header.MessageLength <= 0 {
		return false
	}
	return header.Valid()
}

func (c *Classifier) Protocol() *common.Protocol {
	return common.PROTOCOL_MONGO
}
