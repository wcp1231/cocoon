package proto

type StmtClose struct {
	Status      uint8
	StatementId uint32
}

func UnPackStmtClose(data []byte) (*StmtClose, error) {
	buf := ReadBuffer(data)
	var err error
	stmtClose := &StmtClose{}
	if stmtClose.Status, err = buf.ReadU8(); err != nil {
		return nil, err
	}
	if stmtClose.StatementId, err = buf.ReadU32(); err != nil {
		return nil, err
	}

	return stmtClose, err
}
