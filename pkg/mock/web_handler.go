package mock

import (
	mockModel "cocoon/pkg/model/mock"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (m *MockService) ListMocks(w http.ResponseWriter, r *http.Request) {
	configs := m.GetConfig()
	result := map[string]interface{}{
		"ok":   true,
		"data": configs,
	}
	data, err := json.Marshal(result)
	if err != nil {
		m.logger.Warn("List mock failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"list mocks error\"}")
		return
	}
	w.Write(data)
}

func (m *MockService) AddMocks(w http.ResponseWriter, r *http.Request) {
	var body mockModel.MockConfig
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		m.logger.Warn("Add mock parse body failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"parse request body error\"}")
		return
	}
	m.CreateImposter(body)
	fmt.Fprintf(w, "{\"ok\":true}")
}

func (m *MockService) EditMocks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"error\":\"Not Supported\"}")
}

func (m *MockService) DeleteMocks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"missing id\"}")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		m.logger.Warn("Delete mock Atoi id failed", zap.String("id", idStr), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"invaild id\"}")
		return
	}
	m.DeleteImposter(int32(id))
	fmt.Fprintf(w, "{\"ok\":true}")
}
