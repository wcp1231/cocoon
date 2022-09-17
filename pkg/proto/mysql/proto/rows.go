package proto

import (
	"cocoon/pkg/proto/mysql/packet"
	"errors"
	"fmt"
)

var _ Rows = &TextRows{}

type RowMode int

const (
	TextRowMode RowMode = iota
	BinaryRowMode
)

// NewSimpleRows creates BinaryRows.
func NewSimpleRows(stream *packet.Stream) *SimpleRows {
	binaryRows := &SimpleRows{}
	binaryRows.c = stream
	binaryRows.buffer = NewBuffer(8)
	return binaryRows
}

type SimpleRows struct {
	c            *packet.Stream
	end          bool
	err          *ERR
	data         []byte
	bytes        int
	RowsAffected uint64
	InsertID     uint64
	buffer       *Buffer
	Fields       []*Field
	raw          []byte
}

// Next implements the Rows interface.
// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::ResultsetRow
func (r *SimpleRows) Next() bool {

	if r.end {
		return false
	}

	// if Fields count is 0
	// the packet is OK-Packet without Resultset.
	if len(r.Fields) == 0 {
		r.end = true
		return false
	}

	var err error
	var pkt *packet.Packet
	if pkt, err = r.c.NextPacket(); err != nil {
		r.err = &ERR{InternalError: err}
		r.end = true
		return false
	}
	r.data = pkt.Datas
	r.raw = append(r.raw, pkt.Raw()...)

	switch r.data[0] {
	case EOF_PACKET:
		// This packet may be one of two kinds:
		// - an EOF packet,
		// - an OK packet with an EOF header if
		// sqldb.CLIENT_DEPRECATE_EOF is set.
		r.end = true
		return false

	case ERR_PACKET:
		r.err = UnPackERR(r.data)
		r.end = true
		return false
	}
	r.buffer.Reset(r.data)
	return true
}

// RowValues implements the Rows interface.
// https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::ResultsetRow
func (r *SimpleRows) RowValues() ([]byte, error) {
	if r.Fields == nil {
		return nil, errors.New("rows.Fields is NIL")
	}
	data := r.buffer.Datas()
	return data, nil
}

func (r *SimpleRows) Raw() []byte {
	return r.raw
}

// Rows presents row cursor interface.
type Rows interface {
	Next() bool
	Close() *ERR
	Datas() []byte
	Bytes() int
	//RowsAffected() uint64
	//LastInsertID() uint64
	LastError() *ERR
	//Fields() []*Field
	RowValues() ([]Value, error)
}

// BaseRows --
type BaseRows struct {
	c            packet.Stream
	end          bool
	err          *ERR
	data         []byte
	bytes        int
	RowsAffected uint64
	InsertID     uint64
	buffer       *Buffer
	Fields       []*Field
	raw          []byte
}

// TextRows presents row tuple.
type TextRows struct {
	BaseRows
}

// BinaryRows presents binary row tuple.
type BinaryRows struct {
	BaseRows
}

// Next implements the Rows interface.
// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::ResultsetRow
func (r *BaseRows) Next() bool {

	if r.end {
		return false
	}

	// if Fields count is 0
	// the packet is OK-Packet without Resultset.
	if len(r.Fields) == 0 {
		r.end = true
		return false
	}

	var err error
	var pkt *packet.Packet
	if pkt, err = r.c.NextPacket(); err != nil {
		r.err = &ERR{InternalError: err}
		r.end = true
		return false
	}
	r.data = pkt.Datas
	r.raw = append(r.raw, pkt.Raw()...)

	switch r.data[0] {
	case EOF_PACKET:
		// This packet may be one of two kinds:
		// - an EOF packet,
		// - an OK packet with an EOF header if
		// sqldb.CLIENT_DEPRECATE_EOF is set.
		r.end = true
		return false

	case ERR_PACKET:
		r.err = UnPackERR(r.data)
		r.end = true
		return false
	}
	r.buffer.Reset(r.data)
	return true
}

func (r *BaseRows) Raw() []byte {
	return r.raw
}

// Close drain the rest packets and check the error.
func (r *BaseRows) Close() *ERR {
	for r.Next() {
	}
	return r.LastError()
}

// RowValues implements the Rows interface.
// https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::ResultsetRow
func (r *BaseRows) RowValues() ([]Value, error) {
	if r.Fields == nil {
		return nil, errors.New("rows.Fields is NIL")
	}

	colNumber := len(r.Fields)
	result := make([]Value, colNumber)
	for i := 0; i < colNumber; i++ {
		v, err := r.buffer.ReadLenEncodeBytes()
		if err != nil {
			//r.c.Cleanup()
			return nil, err
		}

		if v != nil {
			r.bytes += len(v)
			result[i] = MakeTrusted(r.Fields[i].Type, v)
		}
	}
	return result, nil
}

// Datas implements the Rows interface.
func (r *BaseRows) Datas() []byte {
	return r.buffer.Datas()
}

// Fields implements the Rows interface.
//func (r *BaseRows) Fields() []*Field {
//	return r.Fields
//}

// Bytes returns all the memory usage which read by this row cursor.
func (r *BaseRows) Bytes() int {
	return r.bytes
}

// RowsAffected implements the Rows interface.
//func (r *BaseRows) RowsAffected() uint64 {
//	return r.RowsAffected
//}

// LastInsertID implements the Rows interface.
//func (r *BaseRows) LastInsertID() uint64 {
//	return r.InsertID
//}

// LastError implements the Rows interface.
func (r *BaseRows) LastError() *ERR {
	return r.err
}

// NewTextRows creates TextRows.
func NewTextRows(stream *packet.Stream) *TextRows {
	textRows := &TextRows{}
	textRows.c = *stream
	textRows.buffer = NewBuffer(8)
	return textRows
}

// NewBinaryRows creates BinaryRows.
func NewBinaryRows(stream *packet.Stream) *BinaryRows {
	binaryRows := &BinaryRows{}
	binaryRows.c = *stream
	binaryRows.buffer = NewBuffer(8)
	return binaryRows
}

// RowValues implements the Rows interface.
// https://dev.mysql.com/doc/internals/en/binary-protocol-resultset-row.html
func (r *BinaryRows) RowValues() ([]Value, error) {
	if r.Fields == nil {
		return nil, errors.New("rows.Fields is NIL")
	}

	header, err := r.buffer.ReadU8()
	if err != nil {
		return nil, err
	}
	if header != OK_PACKET {
		return nil, fmt.Errorf("binary.rows.header.is.not.ok[%v]", header)
	}

	colCount := len(r.Fields)
	// NULL-bitmap,  [(column-count + 7 + 2) / 8 bytes]
	nullMask, err := r.buffer.ReadBytes(int((colCount + 7 + 2) / 8))
	if err != nil {
		return nil, err
	}

	result := make([]Value, colCount)
	for i := 0; i < colCount; i++ {
		// Field is NULL
		// (byte >> bit-pos) % 2 == 1
		if ((nullMask[(i+2)>>3] >> uint((i+2)&7)) & 1) == 1 {
			result[i] = Value{}
			continue
		}

		v, err := ParseMySQLValues(r.buffer, r.Fields[i].Type)
		if err != nil {
			//r.c.Cleanup()
			return nil, err
		}

		if v != nil {
			val, err := BuildValue(v)
			if err != nil {
				//r.c.Cleanup()
				return nil, err
			}
			//r.bytes += val.Len()
			result[i] = val
		} else {
			result[i] = Value{}
		}
	}
	return result, nil
}
