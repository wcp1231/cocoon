package mock

import (
	"cocoon/pkg/model/common"
	"sync"
)

type RequestImposter interface {
	ID() int32
	Match(req common.Message) bool
	Data() common.Message
	GetConfig() interface{}
}

type ProtocolImposter struct {
	imposters []RequestImposter

	mutex sync.RWMutex
}

func (i *ProtocolImposter) Mock(msg common.Message) *common.MockResult {
	result := &common.MockResult{
		Pass: true,
	}
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	for _, imposter := range i.imposters {
		if imposter.Match(msg) {
			result.Pass = false
			result.Data = imposter.Data()
			result.Data.CaptureNow()
			break
		}
	}
	return result
}

func (i *ProtocolImposter) AddRequestImposter(importer RequestImposter) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.imposters = append(i.imposters, importer)
}

func (i *ProtocolImposter) DeleteAllRequestImposters() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.imposters = make([]RequestImposter, 0)
}

func (i *ProtocolImposter) DeleteRequestImposter(id int32) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.imposters = i.deleteById(id)
}

func (i *ProtocolImposter) deleteById(id int32) []RequestImposter {
	ims := i.imposters
	for i, imposter := range ims {
		if imposter.ID() == id {
			return append(ims[:i], ims[i+1:]...)
		}
	}
	return ims
}

func (i *ProtocolImposter) GetConfig() []interface{} {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	var configs []interface{}
	for _, im := range i.imposters {
		configs = append(configs, im.GetConfig())
	}
	return configs
}
