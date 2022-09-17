package proto

import (
	"cocoon/pkg/proto/mysql/packet"
	"errors"
	"fmt"
)

type ResultSet struct {
	// 区分是 TextResultSet 还是 BinaryResultSet
	reqCmd string

	Columns    uint64
	ColumnsEOR []byte
	Fields     []*Field
	Rows       [][]Value
}

type resultSetReader struct {
	stream *packet.Stream
	data   []byte

	// Client capabilities
	capabilities uint32
}

// UnPackResultSet used to unpack the ResultSet packet.
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html
func UnPackResultSet(reqCmd string, capabilities uint32, firstPacket []byte, stream *packet.Stream) (*ResultSet, error) {
	reader := resultSetReader{
		stream:       stream,
		data:         firstPacket,
		capabilities: capabilities,
	}
	return reader.read(reqCmd)
}

func (r *resultSetReader) read(reqCmd string) (*ResultSet, error) {
	columnCount, err := ColumnCount(r.data)
	if err != nil {
		return nil, err
	}

	if columnCount <= 0 {
		return nil, nil
	}

	resultSet := &ResultSet{
		reqCmd:  reqCmd,
		Columns: columnCount,
	}
	err = r.readColumns(resultSet)
	if err != nil {
		return nil, err
	}
	if r.capabilities&CLIENT_DEPRECATE_EOF == 0 {
		if resultSet.ColumnsEOR, err = r.readEOF(); err != nil {
			return nil, err
		}
	}
	if reqCmd == "COM_STMT_EXECUTE" {
		err := r.readBinaryResultSet(resultSet)
		if err != nil {
			return nil, err
		}
	} else if reqCmd == "COM_QUERY" {
		err := r.readTextResultSet(resultSet)
		if err != nil {
			return nil, err
		}
	}
	return resultSet, nil
}

func (r *resultSetReader) readColumns(resultSet *ResultSet) error {
	var err error
	var pkt *packet.Packet
	var field *Field
	resultSet.Fields = make([]*Field, 0, resultSet.Columns)
	var i uint64
	for i = 0; i < resultSet.Columns; i++ {
		if pkt, err = r.stream.NextPacket(); err != nil {
			return err
		}
		if field, err = UnpackColumn(pkt.Datas); err != nil {
			return err
		}
		resultSet.Fields = append(resultSet.Fields, field)
	}
	return nil
}

func (r *resultSetReader) readBinaryResultSet(resultSet *ResultSet) error {
	rowReader := NewBinaryRows(r.stream)
	rowReader.Fields = resultSet.Fields
	var rows [][]Value
	for rowReader.Next() {
		values, err := rowReader.RowValues()
		if err != nil {
			return err
		}
		rows = append(rows, values)
	}
	resultSet.Rows = rows
	return nil
}

func (r *resultSetReader) readTextResultSet(resultSet *ResultSet) error {
	rowReader := NewTextRows(r.stream)
	rowReader.Fields = resultSet.Fields
	var rows [][]Value
	for rowReader.Next() {
		values, err := rowReader.RowValues()
		if err != nil {
			return err
		}
		rows = append(rows, values)
	}
	resultSet.Rows = rows
	return nil
}

func (r *resultSetReader) readEOF() ([]byte, error) {
	pkt, err := r.stream.NextPacket()
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

// resultSetWriter use to pack the ResultSet packets.
type resultSetWriter struct {
	buf        *Buffer
	sequenceID uint8

	resultSet *ResultSet

	// Client capabilities
	capabilities uint32
}

func PackResultSet(resultSet *ResultSet) []byte {
	writer := resultSetWriter{
		buf:        NewBuffer(256),
		sequenceID: 1,
		resultSet:  resultSet,
	}
	if resultSet.reqCmd == "COM_STMT_EXECUTE" {
		return writer.packBinaryResultSet()
	}
	return writer.packTextResultSet()
}

func (w *resultSetWriter) packBinaryResultSet() []byte {
	buf := NewBuffer(64)

	return buf.Datas()
}

func (w *resultSetWriter) packTextResultSet() []byte {
	w.packResultSetColumnCount()
	w.packColumns()
	w.packTextRows()
	return w.buf.Datas()
}

func (w *resultSetWriter) packResultSetColumnCount() {
	data := packet.ToPacketBytesWithSequenceID(PackColumnCount(w.resultSet.Columns), w.sequenceID)
	w.sequenceID += 1
	w.buf.WriteBytes(data)
}

func (w *resultSetWriter) packColumns() {
	for _, field := range w.resultSet.Fields {
		w.buf.WriteBytes(packet.ToPacketBytesWithSequenceID(PackColumn(field), w.sequenceID))
		w.sequenceID += 1
	}
	if w.capabilities&CLIENT_DEPRECATE_EOF == 0 {
		w.writeEOF()
	}
}

func (w *resultSetWriter) packTextRows() {
	for _, row := range w.resultSet.Rows {
		w.buf.WriteBytes(packet.ToPacketBytesWithSequenceID(w.packTextResultSetRow(row), w.sequenceID))
		w.sequenceID += 1
	}
	w.writeEOF()
}

func (w *resultSetWriter) packTextResultSetRow(row []Value) []byte {
	buf := NewBuffer(64)
	for _, value := range row {
		buf.WriteLenEncodeBytes(value.Value)
	}
	return buf.Datas()
}

func (w *resultSetWriter) writeEOF() {
	w.buf.WriteBytes(packet.ToPacketBytesWithSequenceID(w.resultSet.ColumnsEOR, w.sequenceID))
	w.sequenceID += 1
}
