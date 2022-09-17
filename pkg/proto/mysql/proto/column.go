package proto

import (
	"errors"
	"fmt"
)

// ColumnCount returns the column count.
func ColumnCount(payload []byte) (count uint64, err error) {
	buff := ReadBuffer(payload)
	if count, err = buff.ReadLenEncode(); err != nil {
		return 0, errors.New("extracting column count failed")
	}
	return
}

func PackColumnCount(count uint64) []byte {
	buf := NewBuffer(8)
	buf.WriteLenEncode(count)
	return buf.Datas()
}

// UnpackColumn used to unpack the column packet.
// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-Protocol::ColumnDefinition41
func UnpackColumn(payload []byte) (*Field, error) {
	var err error
	field := &Field{}
	buff := ReadBuffer(payload)
	// Catalog is ignored, always set to "def"
	if _, err = buff.ReadLenEncodeString(); err != nil {
		return nil, errors.New("skipping col catalog failed")
	}

	// lenenc_str Schema
	if field.Database, err = buff.ReadLenEncodeString(); err != nil {
		return nil, errors.New("extracting col schema failed")
	}

	// lenenc_str Table
	if field.Table, err = buff.ReadLenEncodeString(); err != nil {
		return nil, errors.New("extracting col table failed")
	}

	// lenenc_str Org_Table
	if field.OrgTable, err = buff.ReadLenEncodeString(); err != nil {
		return nil, errors.New("extracting col org_table failed")
	}

	// lenenc_str Name
	if field.Name, err = buff.ReadLenEncodeString(); err != nil {
		return nil, errors.New("extracting col name failed")
	}

	// lenenc_str Org_Name
	if field.OrgName, err = buff.ReadLenEncodeString(); err != nil {
		return nil, errors.New("extracting col org_name failed")
	}

	// lenenc_int length of fixed-length Fields [0c], skip
	if _, err = buff.ReadLenEncode(); err != nil {
		return nil, errors.New("extracting col 0c failed")
	}

	// 2 character set
	charset, err := buff.ReadU16()
	if err != nil {
		return nil, errors.New("extracting col charset failed")
	}
	field.Charset = uint32(charset)

	// 4 column length
	if field.ColumnLength, err = buff.ReadU32(); err != nil {
		return nil, errors.New("extracting col columnlength failed")
	}

	// 1 type
	t, err := buff.ReadU8()
	if err != nil {
		return nil, errors.New("extracting col type failed")
	}

	// 2 flags
	flags, err := buff.ReadU16()
	if err != nil {
		return nil, errors.New("extracting col flags failed")
	}
	field.Flags = uint32(flags)

	// Convert MySQL type
	if field.Type, err = MySQLToType(int64(t), int64(field.Flags)); err != nil {
		return nil, errors.New(fmt.Sprintf("MySQLToType(%v,%v) failed: %v", t, field.Flags, err))
	}

	// 1 Decimals
	decimals, err := buff.ReadU8()
	if err != nil {
		return nil, errors.New("extracting col type failed")
	}
	field.Decimals = uint32(decimals)

	// 2 Filler and Default Values is ignored
	//
	return field, nil
}

// PackColumn used to pack the column packet.
func PackColumn(field *Field) []byte {
	typ, flags := TypeToMySQL(field.Type)
	if field.Flags != 0 {
		flags = int64(field.Flags)
	}

	buf := NewBuffer(256)

	// lenenc_str Catalog, always 'def'
	buf.WriteLenEncodeString("def")

	// lenenc_str Schema
	buf.WriteLenEncodeString(field.Database)

	// lenenc_str Table
	buf.WriteLenEncodeString(field.Table)

	// lenenc_str Org_Table
	buf.WriteLenEncodeString(field.OrgTable)

	// lenenc_str Name
	buf.WriteLenEncodeString(field.Name)

	// lenenc_str Org_Name
	buf.WriteLenEncodeString(field.OrgName)

	// lenenc_int length of fixed-length Fields [0c]
	buf.WriteLenEncode(uint64(0x0c))

	// 2 character set
	buf.WriteU16(uint16(field.Charset))

	// 4 column length
	buf.WriteU32(field.ColumnLength)

	// 1 type
	buf.WriteU8(byte(typ))

	// 2 flags
	buf.WriteU16(uint16(flags))

	//1 Decimals
	buf.WriteU8(uint8(field.Decimals))

	// 2 filler [00] [00]
	buf.WriteU16(uint16(0))
	return buf.Datas()
}
