package tcp

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/model/rpc"
	"cocoon/pkg/model/traffic"
	"go.uber.org/zap"
)

type TcpManager struct {
	logger  *zap.Logger
	streams map[string]*tcpStream // TODO concurrent map
	resultC chan *traffic.StreamItem
}

func NewTcpManager(logger *zap.Logger, resultC chan *traffic.StreamItem) *TcpManager {
	return &TcpManager{
		logger:  logger,
		streams: map[string]*tcpStream{},
		resultC: resultC,
	}
}

func (t *TcpManager) Accept(packet *rpc.TcpPacket) {
	connectionInfo := &common.ConnectionInfo{
		Source:      packet.Source,
		Destination: packet.Destination,
		IsOutgoing:  packet.IsOutgoing,
	}
	stream := t.getOrCreateStream(connectionInfo)
	stream.Accept(packet)
}

func (t *TcpManager) getOrCreateStream(connectionInfo *common.ConnectionInfo) *tcpStream {
	tcpId := connectionInfo.ID()
	stream, ok := t.streams[tcpId]
	if !ok {
		stream = NewTcpStream(t.logger, connectionInfo, t.resultC)
		t.streams[tcpId] = stream
	}
	return stream
}
