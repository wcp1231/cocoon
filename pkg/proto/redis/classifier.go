package redis

import (
	"bufio"
	"cocoon/pkg/model/common"
	"regexp"
)

type Classifier struct{}

var regex *regexp.Regexp

func init() {
	expr := "^([-+]\\w+|[*:]\\d+|\\$-?\\d+)\r\n"
	regex, _ = regexp.Compile(expr)
}

func (c *Classifier) Match(r *bufio.Reader) bool {
	b, err := r.Peek(10)
	if err != nil {
		return false
	}
	return regex.Match(b)
}

func (c *Classifier) Protocol() *common.Protocol {
	return common.PROTOCOL_REDIS
}
