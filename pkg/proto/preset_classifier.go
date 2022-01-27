package proto

import (
	"cocoon/pkg/model/common"
	"strings"
	"sync"
)

var protocolsMap sync.Map

func InitPresetClassifier(portToProtocol string) {
	protocolsMap.Store("80", common.PROTOCOL_HTTP)
	protocolsMap.Store("6379", common.PROTOCOL_REDIS)
	protocolsMap.Store("3306", common.PROTOCOL_MYSQL)
	protocolsMap.Store("27017", common.PROTOCOL_MONGO)
	processConfigString(portToProtocol)
}

func classifyByDst(dst string) *common.Protocol {
	pair := strings.Split(dst, ":")
	if len(pair) < 2 {
		return nil
	}
	proto, exist := protocolsMap.Load(pair[1])
	if !exist {
		return nil
	}
	return proto.(*common.Protocol)
}

func processConfigString(portToProtocol string) {
	protocols := strings.Split(portToProtocol, ",")
	for _, proto := range protocols {
		if len(proto) < 2 {
			continue
		}
		pairs := strings.Split(proto, ":")
		port := pairs[0]
		protocol := nameToProtocol(pairs[1])
		if protocol != nil {
			protocolsMap.Store(port, protocol)
		}
	}
}

func nameToProtocol(name string) *common.Protocol {
	switch strings.ToLower(name) {
	case "dubbo":
		return common.PROTOCOL_DUBBO
	case "http":
		return common.PROTOCOL_HTTP
	case "mongo":
		return common.PROTOCOL_MONGO
	case "mysql":
		return common.PROTOCOL_MYSQL
	case "redis":
		return common.PROTOCOL_REDIS
	}
	return nil
}
