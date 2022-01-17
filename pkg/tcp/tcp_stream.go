package tcp

import (
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"cocoon/pkg/model/rpc"
	"cocoon/pkg/model/traffic"
	"go.uber.org/zap"
	"sync"
)

type tcpStream struct {
	logger         *zap.Logger
	connectionInfo *common.ConnectionInfo
	resultC        chan *traffic.StreamItem
	dissectC       chan *api.DissectResult

	wg             sync.WaitGroup
	nextSeq        uint64
	requestReader  *tcpReader
	responseReader *tcpReader
}

func NewTcpStream(logger *zap.Logger, connectionInfo *common.ConnectionInfo, resultC chan *traffic.StreamItem) *tcpStream {
	dissectC := make(chan *api.DissectResult, 1024)

	stream := &tcpStream{
		logger:         logger,
		connectionInfo: connectionInfo,
		resultC:        resultC,
		dissectC:       dissectC,
		nextSeq:        0,
		wg:             sync.WaitGroup{},
	}

	requestReader := newTcpReader(connectionInfo, stream, true, dissectC)
	responseReader := newTcpReader(connectionInfo, stream, false, dissectC)

	stream.requestReader = requestReader
	stream.responseReader = responseReader

	stream.wg.Add(2)
	go requestReader.run(&stream.wg)
	go responseReader.run(&stream.wg)
	go stream.handleDissectResult()

	return stream
}

func (s *tcpStream) Accept(packet *rpc.TcpPacket) {
	if packet.Seq < s.nextSeq || s.nextSeq == 0 {
		// 新链接
		// TODO 初始化
		s.nextSeq = packet.Seq
	}

	s.nextSeq = packet.Seq + 1
	dataMsg := tcpReaderDataMsg{
		bytes:     packet.Payload,
		timestamp: packet.Timestamp,
	}
	if packet.IsRequest() {
		s.requestReader.msgQueue <- dataMsg
	} else {
		s.responseReader.msgQueue <- dataMsg
	}
}

func (s *tcpStream) handleDissectResult() {
	for {
		result := <-s.dissectC
		s.logger.Info("Dissect result", zap.String("result", result.String()))
	}
}
