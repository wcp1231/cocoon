package proto

import (
	"cocoon/pkg/proto/mysql/packet"
	"errors"
	"fmt"
)

const (
	// EOF_PACKET is the EOF packet.
	EOF_PACKET byte = 0xfe
)

func readEOF(stream *packet.Stream) ([]byte, error) {
	pkt, err := stream.NextPacket()
	if err != nil {
		return nil, err
	}
	data := pkt.Datas
	switch data[0] {
	case EOF_PACKET:
		return data, nil
	case ERR_PACKET:
		return nil, UnPackERR(data).ToError()
	default:
		return nil, errors.New(fmt.Sprintf("unexpected.eof.packet[%+v]", data))
	}
}

func writeEOF(buf *Buffer, origin []byte, sequenceID uint8) {
	buf.WriteBytes(packet.ToPacketBytesWithSequenceID(origin, sequenceID))
}
