package dissector

import (
	"cocoon/pkg/dissector/http"
	"cocoon/pkg/model/api"
)

type DissectProcessor struct {
	ID        string
	isRequest bool
	resultC   chan *api.DissectResult

	dissectorIdx int
	dissectors   []api.Dissector

	nextSeq uint64
}

func (d *DissectProcessor) Process(b api.TcpReader) {
	for d.dissectorIdx < len(d.dissectors) {
		err := d.dissectors[d.dissectorIdx].Dissect(b, d.isRequest)
		if err != nil {
			b.Reset()
			d.dissectorIdx += 1
		}
	}
}

func (d *DissectProcessor) findMatchedDissector(b api.TcpReader) api.Dissector {
	for _, dissector := range d.dissectors {
		if dissector.Match(b) {
			return dissector
		}
	}
	// 返回默认处理器
	return d.dissectors[len(d.dissectors)-1]
}

func NewDissectProcessor(tcpId string, isRequest bool, resultC chan *api.DissectResult) *DissectProcessor {
	var dissectors []api.Dissector
	dissectors = append(dissectors, http.NewDissector(resultC))
	dissectors = append(dissectors, newDefaultDissector(resultC))
	return &DissectProcessor{
		ID:        tcpId,
		isRequest: isRequest,
		resultC:   resultC,

		dissectorIdx: 0,
		dissectors:   dissectors,
	}
}
