package mock

import (
	"cocoon/pkg/model/common"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
)

type MockService struct {
	logger    *zap.Logger
	importers map[string]*ProtocolImposter
}

func NewMockService(logger *zap.Logger) *MockService {
	return &MockService{
		logger:    logger,
		importers: map[string]*ProtocolImposter{},
	}
}

func (m *MockService) Mock(proto string, msg *common.GenericMessage) *common.MockResult {
	protoImposter, exist := m.importers[proto]
	if !exist {
		result := &common.MockResult{
			Pass: true,
		}
		return result
	}
	return protoImposter.Mock(msg)
}

func (m *MockService) AddImposter(proto string, importer RequestImposter) {
	protoImposter, exist := m.importers[proto]
	if !exist {
		protoImposter = &ProtocolImposter{}
		m.importers[proto] = protoImposter
	}
	protoImposter.AddRequestImposter(importer)

	m.logger.Info("Add imposter", zap.String("proto", proto), zap.String("map", fmt.Sprintf("%+v", m.importers)))
}

// TODO 临时测试
func (m *MockService) InitFromFile() error {
	fileData, err := os.ReadFile("./mock.json")
	if err != nil {
		m.logger.Error("Read mock file failed", zap.Error(err))
		return err
	}
	var config mockConfig
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		m.logger.Error("Parse mock json failed", zap.Error(err))
		return err
	}
	m.logger.Info("Init mock by config", zap.String("config", fmt.Sprintf("%v", config)))
	for _, httpConfig := range config.Http {
		imposter := newHttpRequestMatcherFromConfig(httpConfig)
		m.AddImposter(common.PROTOCOL_HTTP.Name, imposter)
	}
	for _, redisConfig := range config.Redis {
		imposter := newRedisRequestMatcherFromConfig(redisConfig)
		m.AddImposter(common.PROTOCOL_REDIS.Name, imposter)
	}
	return nil
}
