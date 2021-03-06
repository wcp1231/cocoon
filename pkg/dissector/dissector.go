package dissector

import (
	"bufio"
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"io"
)

type DefaultDissector struct {
	resultC chan *api.DissectResult
}

func newDefaultDissector(resultC chan *api.DissectResult) api.Dissector {
	return &DefaultDissector{
		resultC: resultC,
	}
}

func (d *DefaultDissector) Name() string {
	return "Default Dissector"
}

func (d *DefaultDissector) Dissect(reader api.TcpReader, isRequest bool) error {
	br := reader.BufferReader()
	buffered := d.tryRead(br)
	if buffered <= 0 {
		return io.EOF
	}
	raw := make([]byte, buffered)
	br.Read(raw)
	result := &api.DissectResult{}
	result.ConnectionInfo = reader.Connection()
	result.IsRequest = isRequest
	result.Protocol = common.PROTOCOL_UNKNOWN
	result.Payload = &common.GenericMessage{
		Raw: &raw,
	}

	d.resultC <- result
	return nil
}

func (d *DefaultDissector) tryRead(br *bufio.Reader) int {
	buffered := br.Buffered()
	if buffered > 0 {
		return buffered
	}
	br.Peek(1)
	return br.Buffered()
}
