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
	InternalError error

	Header       byte // always 0xff
	ErrorCode    uint16
	SQLState     string
	ErrorMessage string
}

func (e *ERR) ToError() error {
	return errors.New(fmt.Sprintf("%d %s %s", e.ErrorCode, e.SQLState, e.ErrorMessage))
}

// UnPackERR parses the error packet and returns a sqldb.SQLError.
// https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
// TODO 不要转化成 golang 的 error
func UnPackERR(data []byte) *ERR {
	var err error
	e := &ERR{}
	buf := ReadBuffer(data)
	if e.Header, err = buf.ReadU8(); err != nil {
		e.InternalError = errors.New(fmt.Sprintf("invalid error packet header: %v", data))
		return e
	}
	if e.Header != ERR_PACKET {
		e.InternalError = errors.New(fmt.Sprintf("invalid error packet header: %v", e.Header))
		return e
	}
	if e.ErrorCode, err = buf.ReadU16(); err != nil {
		e.InternalError = errors.New(fmt.Sprintf("invalid error packet code: %v", data))
		return e
	}

	// Skip SQLStateMarker
	if _, err = buf.ReadString(1); err != nil {
		e.InternalError = errors.New(fmt.Sprintf("invalid error packet marker: %v", data))
		return e
	}
	if e.SQLState, err = buf.ReadString(5); err != nil {
		e.InternalError = errors.New(fmt.Sprintf("invalid error packet sqlstate: %v", data))
		return e
	}
	msgLen := len(data) - buf.Seek()
	if e.ErrorMessage, err = buf.ReadString(msgLen); err != nil {
		e.InternalError = errors.New(fmt.Sprintf("invalid error packet message: %v", data))
		return e
	}
	return e
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
