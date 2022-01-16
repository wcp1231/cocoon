package dissector

import (
	"cocoon/pkg/dissector/http"
	"cocoon/pkg/model/api"
	"sync"
)

type DissectProcessor struct {
	sync.Mutex
	ID        string
	isRequest bool
	resultC   chan *api.DissectResult

	maybeDissector api.Dissector
	dissectors     []api.Dissector

	nextSeq uint64
}

func (d *DissectProcessor) Process(b api.TcpReader) {
	defer d.Unlock()
	d.Lock()
	if d.maybeDissector == nil {
		// 链接重建，重置 dissector
		d.maybeDissector = d.findMatchedDissector(b)
	}

	d.maybeDissector.Dissect(b, d.isRequest)
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

		dissectors: dissectors,
	}
}
