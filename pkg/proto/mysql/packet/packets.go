package packet

import (
	"bufio"
	"errors"
)

const (
	// PACKET_MAX_SIZE used for the max packet size.
	PACKET_MAX_SIZE = (1<<24 - 1) // 16MB - 1
)

type Packet struct {
	Header     []byte
	SequenceID byte
	Datas      []byte
}

func NewPacket() *Packet {
	return &Packet{
		Header: []byte{0, 0, 0, 0},
	}
}

func (p *Packet) Raw() []byte {
	l := len(p.Datas)
	raw := make([]byte, 4+l)
	raw[0] = byte(l)
	raw[1] = byte(l >> 8)
	raw[2] = byte(l >> 16)
	raw[3] = p.SequenceID
	copy(raw[4:], p.Datas)
	return raw
}

type Packets struct {
	seq    uint8
	stream *Stream
}

func NewPackets(reader *bufio.Reader) *Packets {
	return &Packets{
		stream: NewStream(reader, PACKET_BUFFER_SIZE),
	}
}

func (p *Packets) Next() ([]byte, error) {
	pkt, err := p.stream.Read()
	if err != nil {
		return nil, err
	}

	if pkt.SequenceID != p.seq {
		return nil, errors.New("pkt.read.seq != pkt.actual.seq")
	}
	p.seq++
	return pkt.Datas, nil
}
