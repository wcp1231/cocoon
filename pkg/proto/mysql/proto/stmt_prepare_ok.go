package proto

import "cocoon/pkg/proto/mysql/packet"

type StmtPrepareOK struct {
	Status       uint8
	StatementId  uint32
	ColumnsCount uint16
	ParamsCount  uint16
	Reserved     uint8
	WarningCount uint16
	Params       []*Field
	ParamsEOF    []byte
	Columns      []*Field
	ColumnsEOF   []byte
}

func UnPackStmtPrepareOk(firstPacket []byte, stream *packet.Stream) (*StmtPrepareOK, error) {
	prepareOk := &StmtPrepareOK{}
	err := unPackStmtPrepareOkFirstPacket(firstPacket, prepareOk)
	if err != nil {
		return nil, err
	}
	err = unPackStmtPrepareOkParams(stream, prepareOk)
	if err != nil {
		return nil, err
	}
	err = unPackStmtPrepareOkColumns(stream, prepareOk)
	return prepareOk, err
}

func unPackStmtPrepareOkFirstPacket(data []byte, prepareOk *StmtPrepareOK) error {
	buf := ReadBuffer(data)
	var err error
	if prepareOk.Status, err = buf.ReadU8(); err != nil {
		return err
	}
	if prepareOk.StatementId, err = buf.ReadU32(); err != nil {
		return err
	}
	if prepareOk.ColumnsCount, err = buf.ReadU16(); err != nil {
		return err
	}
	if prepareOk.ParamsCount, err = buf.ReadU16(); err != nil {
		return err
	}
	if prepareOk.Reserved, err = buf.ReadU8(); err != nil {
		return err
	}
	if prepareOk.WarningCount, err = buf.ReadU16(); err != nil {
		return err
	}
	return nil
}

func unPackStmtPrepareOkParams(stream *packet.Stream, prepareOk *StmtPrepareOK) error {
	var err error
	var pkt *packet.Packet
	var field *Field
	var i uint16
	for i = 0; i < prepareOk.ParamsCount; i++ {
		if pkt, err = stream.NextPacket(); err != nil {
			return err
		}
		if field, err = UnpackColumn(pkt.Datas); err != nil {
			return err
		}
		prepareOk.Params = append(prepareOk.Params, field)
	}
	if prepareOk.ParamsEOF, err = readEOF(stream); err != nil {
		return err
	}
	return nil
}

func unPackStmtPrepareOkColumns(stream *packet.Stream, prepareOk *StmtPrepareOK) error {
	var err error
	var pkt *packet.Packet
	var field *Field
	var i uint16
	for i = 0; i < prepareOk.ColumnsCount; i++ {
		if pkt, err = stream.NextPacket(); err != nil {
			return err
		}
		if field, err = UnpackColumn(pkt.Datas); err != nil {
			return err
		}
		prepareOk.Columns = append(prepareOk.Columns, field)
	}
	if prepareOk.ColumnsEOF, err = readEOF(stream); err != nil {
		return err
	}
	return nil
}

func PackStmtPrepareOk(prepareOk *StmtPrepareOK) []byte {
	buf := NewBuffer(256)
	var sequenceID uint8 = 1
	firstBytes, sequenceID := packStmtPrepareOkFirstPacket(prepareOk, sequenceID)
	buf.WriteBytes(firstBytes)
	paramsBytes, sequenceID := packStmtPrepareOkParams(prepareOk, sequenceID)
	buf.WriteBytes(paramsBytes)
	columnsBytes, sequenceID := packStmtPrepareOkColumns(prepareOk, sequenceID)
	buf.WriteBytes(columnsBytes)
	return buf.Datas()
}

func packStmtPrepareOkFirstPacket(prepareOk *StmtPrepareOK, sequenceId uint8) ([]byte, uint8) {
	buf := NewBuffer(32)
	buf.WriteU8(prepareOk.Status)
	buf.WriteU32(prepareOk.StatementId)
	buf.WriteU16(prepareOk.ColumnsCount)
	buf.WriteU16(prepareOk.ParamsCount)
	buf.WriteU8(prepareOk.Reserved)
	buf.WriteU16(prepareOk.WarningCount)
	return packet.ToPacketBytesWithSequenceID(buf.Datas(), sequenceId), sequenceId + 1
}

func packStmtPrepareOkParams(prepareOk *StmtPrepareOK, sequenceId uint8) ([]byte, uint8) {
	buf := NewBuffer(32)
	for _, field := range prepareOk.Params {
		buf.WriteBytes(packet.ToPacketBytesWithSequenceID(PackColumn(field), sequenceId))
		sequenceId += 1
	}
	if prepareOk.ParamsCount > 0 {
		writeEOF(buf, prepareOk.ParamsEOF, sequenceId)
		sequenceId += 1
	}
	return buf.Datas(), sequenceId
}

func packStmtPrepareOkColumns(prepareOk *StmtPrepareOK, sequenceId uint8) ([]byte, uint8) {
	buf := NewBuffer(32)
	for _, field := range prepareOk.Columns {
		buf.WriteBytes(packet.ToPacketBytesWithSequenceID(PackColumn(field), sequenceId))
		sequenceId += 1
	}
	if prepareOk.ColumnsCount > 0 {
		writeEOF(buf, prepareOk.ColumnsEOF, sequenceId)
		sequenceId += 1
	}
	return buf.Datas(), sequenceId
}
