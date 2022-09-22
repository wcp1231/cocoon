package proto

type StmtExecute struct {
	Status         uint8
	StatementId    uint32
	Flag           uint8
	IterationCount uint32
	Params         []Value
}

func UnPackStmtExecute(data []byte, stmtMap map[uint32]uint16) (*StmtExecute, error) {
	buf := ReadBuffer(data)
	var err error
	stmtExecute := &StmtExecute{}
	if stmtExecute.Status, err = buf.ReadU8(); err != nil {
		return nil, err
	}
	if stmtExecute.StatementId, err = buf.ReadU32(); err != nil {
		return nil, err
	}
	if stmtExecute.Flag, err = buf.ReadU8(); err != nil {
		return nil, err
	}
	if stmtExecute.IterationCount, err = buf.ReadU32(); err != nil {
		return nil, err
	}

	paramsCount := int(stmtMap[stmtExecute.StatementId])
	if paramsCount <= 0 {
		return stmtExecute, nil
	}

	nullMask, err := buf.ReadBytes((paramsCount + 7) / 8)
	if err != nil {
		return nil, err
	}
	boundFlag, err := buf.ReadU8()
	if err != nil {
		return nil, err
	}
	if boundFlag == 0x01 {
		values, err := unPackStmtExecute(buf, nullMask, paramsCount)
		if err != nil {
			return nil, err
		}
		stmtExecute.Params = values
	}
	return stmtExecute, err
}

func unPackStmtExecute(buf *Buffer, nullMask []byte, paramsCount int) ([]Value, error) {
	result := make([]Value, paramsCount)
	typeBytes, err := buf.ReadBytes(int(paramsCount) * 2)
	if err != nil {
		return nil, err
	}
	typeBuf := ReadBuffer(typeBytes)
	for i := 0; i < paramsCount; i++ {
		typeByte, err := typeBuf.ReadU8()
		if err != nil {
			return nil, err
		}
		flag, err := typeBuf.ReadU8()
		if err != nil {
			return nil, err
		}
		if ((nullMask[i>>3] >> uint(i&7)) & 1) == 1 {
			result[i] = Value{}
			continue
		}
		typ, err := MySQLToType(int64(typeByte), int64(flag))
		if err != nil {
			return nil, err
		}
		v, err := ParseMySQLValues(buf, typ)
		if err != nil {
			return nil, err
		}
		val, err := BuildValue(v, typ)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}
