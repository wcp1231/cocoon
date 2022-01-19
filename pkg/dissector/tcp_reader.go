package dissector

import (
	"bufio"
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"io"
	"sync"
	"time"
)

type tcpReaderDataMsg struct {
	bytes     []byte
	timestamp time.Time
}

/* tcpReader gets reads from a channel of bytes of tcp payload, and parses it into requests and responses.
 * The payload is written to the channel by a tcpStream object that is dedicated to one tcp connection.
 * An tcpReader object is unidirectional: it parses either a client stream or a server stream.
 * Implements io.Reader interface (Read)
 */
type tcpReader struct {
	connectionInfo *common.ConnectionInfo
	isClosed       bool
	isClient       bool
	msgQueue       chan tcpReaderDataMsg // Channel of captured reassembled tcp payload
	data           []byte
	preData        []byte
	preLen         int
	//superTimer  *api.SuperTimer
	parent      *tcpStream
	packetsSeen uint
	//outboundLinkWriter *OutboundLinkWriter
	//extension          *api.Extension
	//emitter            api.Emitter
	//counterPair        *api.CounterPair

	dissector *DissectProcessor
	dissectC  chan *api.DissectResult
	br        *bufio.Reader

	sync.Mutex
}

func newTcpReader(connectionInfo *common.ConnectionInfo, parent *tcpStream, isRequest bool, dissectC chan *api.DissectResult) *tcpReader {
	tcpReader := &tcpReader{
		connectionInfo: connectionInfo,
		isClient:       true,
		isClosed:       false,
		parent:         parent,
		msgQueue:       make(chan tcpReaderDataMsg),
		data:           []byte{},
		packetsSeen:    0,
		dissectC:       dissectC,
		dissector:      NewDissectProcessor(connectionInfo.ID(), isRequest, dissectC),
	}
	tcpReader.br = bufio.NewReader(tcpReader)
	return tcpReader
}

func (h *tcpReader) Read(p []byte) (int, error) {
	var msg tcpReaderDataMsg

	ok := true
	for ok && len(h.data) == 0 {
		msg, ok = <-h.msgQueue
		h.data = msg.bytes

		//h.superTimer.CaptureTime = msg.timestamp
		if len(h.data) > 0 {
			h.packetsSeen += 1
		}
		//if h.packetsSeen < checkTLSPacketAmount && len(msg.bytes) > 5 { // packets with less than 5 bytes cause tlsx to panic
		//	clientHello := tlsx.ClientHello{}
		//	err := clientHello.Unmarshall(msg.bytes)
		//	if err == nil {
		//		// logger.Log.Debugf("Detected TLS client hello with SNI %s", clientHello.SNI)
		//		// TODO: Throws `panic: runtime error: invalid memory address or nil pointer dereference` error.
		//		// numericPort, _ := strconv.Atoi(h.tcpID.DstPort)
		//		// h.outboundLinkWriter.WriteOutboundLink(h.tcpID.SrcIP, h.tcpID.DstIP, numericPort, clientHello.SNI, TLSProtocol)
		//	}
		//}
	}
	if !ok || len(h.data) == 0 {
		return 0, io.EOF
	}

	l := copy(p, h.data)
	h.preData = make([]byte, l)
	h.preLen = copy(h.preData, p)
	h.data = h.data[l:]
	return l, nil
}

func (h *tcpReader) Connection() *common.ConnectionInfo {
	return h.connectionInfo
}

func (h *tcpReader) BufferReader() *bufio.Reader {
	return h.br
}

// 将上次读取的内容回填再读
func (h *tcpReader) Reset() {
	h.data = append(h.preData[:h.preLen], h.data...)
	h.br.Reset(h)
}

func (h *tcpReader) Close() {
	h.Lock()
	if !h.isClosed {
		h.isClosed = true
		close(h.msgQueue)
	}
	h.Unlock()
}

func (h *tcpReader) run(wg *sync.WaitGroup) {
	defer wg.Done()
	h.dissector.Process(h)

	/*
		err := h.extension.Dissector.Dissect(b, h.isClient, h.tcpID, h.counterPair, h.superTimer, h.parent.superIdentifier, h.emitter, filteringOptions)
		if err != nil {
			io.Copy(ioutil.Discard, b)
		}
	*/
}
