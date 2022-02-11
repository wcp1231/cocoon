package proto

import (
	"errors"
	"fmt"
)

const (
	// ERR_PACKET is the error packet byte.
	ERR_PACKET byte = 0xff
)

// ERR is the error packet.
type ERR struct {
	Header       byte // always 0xff
	ErrorCode    uint16
	SQLState     string
	ErrorMessage string
}

// UnPackERR parses the error packet and returns a sqldb.SQLError.
// https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
func UnPackERR(data []byte) error {
	var err error
	e := &ERR{}
	buf := ReadBuffer(data)
	if e.Header, err = buf.ReadU8(); err != nil {
		return errors.New(fmt.Sprintf("invalid error packet header: %v", data))
	}
	if e.Header != ERR_PACKET {
		return errors.New(fmt.Sprintf("invalid error packet header: %v", e.Header))
	}
	if e.ErrorCode, err = buf.ReadU16(); err != nil {
		return errors.New(fmt.Sprintf("invalid error packet code: %v", data))
	}

	// Skip SQLStateMarker
	if _, err = buf.ReadString(1); err != nil {
		return errors.New(fmt.Sprintf("invalid error packet marker: %v", data))
	}
	if e.SQLState, err = buf.ReadString(5); err != nil {
		return errors.New(fmt.Sprintf("invalid error packet sqlstate: %v", data))
	}
	msgLen := len(data) - buf.Seek()
	if e.ErrorMessage, err = buf.ReadString(msgLen); err != nil {
		return errors.New(fmt.Sprintf("invalid error packet message: %v", data))
	}
	return errors.New(fmt.Sprintf("%d %s %s", e.ErrorCode, e.SQLState, e.ErrorMessage))
}

// PackERR used to pack the error packet.
func PackERR(e *ERR) []byte {
	buf := NewBuffer(64)

	buf.WriteU8(ERR_PACKET)

	// error code
	buf.WriteU16(e.ErrorCode)

	// sql-state marker #
	buf.WriteU8('#')

	// sql-state (?) 5 ascii bytes
	if e.SQLState == "" {
		e.SQLState = "HY000"
	}
	if len(e.SQLState) != 5 {
		panic("sqlState has to be 5 characters long")
	}
	buf.WriteString(e.SQLState)

	// error msg
	buf.WriteString(e.ErrorMessage)
	return buf.Datas()
}