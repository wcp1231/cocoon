package http

import (
	"bufio"
	"cocoon/pkg/model/common"
	"regexp"
	"strings"
)

type Classifier struct{}

var httpVerbs = []string{
	"CONNECT",
	"DELETE",
	"GET",
	"HEAD",
	"OPTIONS",
	"PATCH",
	"POST",
	"PUT",
	"TRACE",
}

var regex *regexp.Regexp

func init() {
	expr := "^(" + strings.Join(httpVerbs, "|") + ")"
	regex, _ = regexp.Compile(expr)
}

func (c *Classifier) Match(r *bufio.Reader) bool {
	b, err := r.Peek(16)
	if err != nil {
		return false
	}
	return regex.Match(b)
}

func (c *Classifier) Protocol() *common.Protocol {
	return common.PROTOCOL_HTTP
}
