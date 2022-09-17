package proto

import (
	"errors"
	"fmt"
)

const (
	// OK_PACKET is the OK byte.
	OK_PACKET byte = 0x00
)

// OK used for OK packet.
type OK struct {
	Header       byte // 0x00
	AffectedRows uint64
	LastInsertID uint64
	StatusFlags  uint16
	Warnings     uint16
}

// UnPackOK used to unpack the OK packet.
// https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
func UnPackOK(data []byte) (*OK, error) {
	var err error
	o := &OK{}
	buf := ReadBuffer(data)

	// header
	if o.Header, err = buf.ReadU8(); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid ok packet header: %v", data))
	}
	if o.Header != OK_PACKET {
		return nil, errors.New(fmt.Sprintf("invalid ok packet header: %v", o.Header))
	}

	// AffectedRows
	if o.AffectedRows, err = buf.ReadLenEncode(); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid ok packet affectedrows: %v", data))
	}

	// LastInsertID
	if o.LastInsertID, err = buf.ReadLenEncode(); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid ok packet lastinsertid: %v", data))
	}

	// Status
	if o.StatusFlags, err = buf.ReadU16(); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid ok packet statusflags: %v", data))
	}

	// Warnings
	if o.Warnings, err = buf.ReadU16(); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid ok packet warnings: %v", data))
	}
	return o, nil
}

// PackOK used to pack the OK packet.
func PackOK(o *OK) []byte {
	buf := NewBuffer(64)

	// OK
	buf.WriteU8(OK_PACKET)

	// affected rows
	buf.WriteLenEncode(o.AffectedRows)

	// last insert id
	buf.WriteLenEncode(o.LastInsertID)

	// status
	buf.WriteU16(o.StatusFlags)

	// warnings
	buf.WriteU16(o.Warnings)
	return buf.Datas()
}
