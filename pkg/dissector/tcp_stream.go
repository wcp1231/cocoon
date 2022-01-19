package dissector

import (
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
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

func newTcpStream(logger *zap.Logger, connectionInfo *common.ConnectionInfo, resultC chan *traffic.StreamItem) *tcpStream {
	dissectC := make(chan *api.DissectResult)

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

func (s *tcpStream) Accept(packet *common.TcpPacket) {
	if packet.Seq < s.nextSeq || s.nextSeq == 0 {
		// TODO 新链接, 初始化
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

func (s *tcpStream) Close() {
	s.requestReader.Close()
	s.responseReader.Close()
	s.wg.Wait()
	close(s.dissectC)
}

func (s *tcpStream) handleDissectResult() {
	for {
		diss, more := <-s.dissectC
		if !more {
			break
		}
		s.logger.Info("Dissect result", zap.String("req", diss.String()))
	}
}
