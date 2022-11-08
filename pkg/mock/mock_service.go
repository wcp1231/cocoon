package mock

import (
	"cocoon/pkg/model/common"
	mockModel "cocoon/pkg/model/mock"
	"cocoon/pkg/proto/http"
	"cocoon/pkg/proto/redis"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync/atomic"
)

type MockService struct {
	logger    *zap.Logger
	importers map[string]*ProtocolImposter

	id int32
}

func NewMockService(logger *zap.Logger) *MockService {
	return &MockService{
		logger:    logger,
		importers: map[string]*ProtocolImposter{},
	}
}

func (m *MockService) Mock(proto string, msg common.Message) *common.MockResult {
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

func (m *MockService) DeleteAllImposters() {
	for _, protoImposter := range m.importers {
		protoImposter.DeleteAllRequestImposters()
	}
	m.logger.Info("Delete all imposters")
}

func (m *MockService) DeleteImposter(id int32) {
	for _, protoImposter := range m.importers {
		protoImposter.DeleteRequestImposter(id)
	}
	m.logger.Info("Delete imposter", zap.Int32("id", id))
}

func (m *MockService) CreateImposter(config mockModel.MockConfig) {
	m.logger.Info("Create mock by config", zap.String("config", fmt.Sprintf("%v", config)))
	for _, httpConfig := range config.Http {
		id := atomic.AddInt32(&m.id, 1)
		imposter := http.NewHttpRequestMatcherFromConfig(httpConfig, id)
		m.AddImposter(common.PROTOCOL_HTTP.Name, imposter)
	}
	for _, redisConfig := range config.Redis {
		id := atomic.AddInt32(&m.id, 1)
		imposter := redis.NewRedisRequestMatcherFromConfig(redisConfig, id)
		m.AddImposter(common.PROTOCOL_REDIS.Name, imposter)
	}
}

func (m *MockService) GetConfig() map[string][]interface{} {
	result := make(map[string][]interface{})
	ims, ok := m.importers[common.PROTOCOL_HTTP.Name]
	if ok && ims != nil {
		result[common.PROTOCOL_HTTP.Name] = ims.GetConfig()
	}
	ims, ok = m.importers[common.PROTOCOL_REDIS.Name]
	if ok && ims != nil {
		result[common.PROTOCOL_REDIS.Name] = ims.GetConfig()
	}
	return result
}

// TODO 临时测试
func (m *MockService) InitFromFile() error {
	fileData, err := os.ReadFile("./mock.json")
	if err != nil {
		m.logger.Error("Read mock file failed", zap.Error(err))
		return err
	}
	var config mockModel.MockConfig
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		m.logger.Error("Parse mock json failed", zap.Error(err))
		return err
	}
	m.CreateImposter(config)
	return nil
}
