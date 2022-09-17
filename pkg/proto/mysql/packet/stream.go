package packet

import (
	"bufio"
	"io"
)

const (
	// PACKET_BUFFER_SIZE is how much we buffer for reading.
	PACKET_BUFFER_SIZE = 32 * 1024
)

type Stream struct {
	pktMaxSize int
	reader     *bufio.Reader
}

func NewStream(reader *bufio.Reader, pktMaxSize int) *Stream {
	return &Stream{
		pktMaxSize: pktMaxSize,
		reader:     reader,
	}
}

func (s *Stream) Read() (*Packet, error) {
	pkt := NewPacket()

	// TODO 只读 buffered 的部分
	// Header.
	if _, err := io.ReadFull(s.reader, pkt.Header); err != nil {
		return nil, err
	}

	pkt.SequenceID = pkt.Header[3]
	length := int(uint32(pkt.Header[0]) | uint32(pkt.Header[1])<<8 | uint32(pkt.Header[2])<<16)
	if length == 0 {
		return pkt, nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(s.reader, data); err != nil {
		return nil, err
	}
	pkt.Datas = data

	// single packet.
	if length < s.pktMaxSize {
		return pkt, nil
	}

	// There is more than one packet, read them all.
	next, err := s.Read()
	if err != nil {
		return nil, err
	}
	pkt.SequenceID = next.SequenceID
	pkt.Datas = append(pkt.Datas, next.Datas...)
	return pkt, nil
}

func (s *Stream) NextPacket() (*Packet, error) {
	pkt, err := s.Read()
	if err != nil {
		return nil, err
	}
	return pkt, nil
}
