package record

import (
	"cocoon/pkg/model/common"
	"encoding/json"
	"go.uber.org/zap"
)

type RecordService struct {
	logger *zap.Logger

	// Registered clients.
	clients map[*WsConn]bool
	// Register requests from the clients.
	register chan *WsConn
	// Unregister requests from clients.
	unregister chan *WsConn
}

func NewRecordService(logger *zap.Logger) *RecordService {
	return &RecordService{
		logger: logger,

		clients:    map[*WsConn]bool{},
		register:   make(chan *WsConn),
		unregister: make(chan *WsConn),
	}
}

func (r *RecordService) RecordRequest(request *common.GenericMessage) error {
	record := fromGenericMessage(request, true)
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	r.broadcast(data)
	return nil
}

func (r *RecordService) RecordResponse(request, response *common.GenericMessage) error {
	response.Id = request.Id
	record := fromGenericMessage(response, false)
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	r.broadcast(data)
	return nil
}

func (r *RecordService) broadcast(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client)
		}
	}
}

func (r *RecordService) Start() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		}
	}
}
