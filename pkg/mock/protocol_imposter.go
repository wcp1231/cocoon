package mock

import (
	"cocoon/pkg/model/common"
	"fmt"
)

type RequestImposter interface {
	Match(req *common.GenericMessage) bool
	Data() *[]byte
}

type ProtocolImposter struct {
	imposters []RequestImposter
}

func (i *ProtocolImposter) Mock(msg *common.GenericMessage) *common.MockResult {
	result := &common.MockResult{
		Pass: true,
	}
	for _, imposter := range i.imposters {
		if imposter.Match(msg) {
			result.Pass = false
			result.Data = imposter.Data()
			break
		}
	}
	return result
}

func (i *ProtocolImposter) AddRequestImposter(importer RequestImposter) {
	i.imposters = append(i.imposters, importer)
	fmt.Printf("Add request imposters, len=%d\n", len(i.imposters))
}