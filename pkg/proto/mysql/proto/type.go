package proto

import "fmt"

// This file provides wrappers and support
// functions for querypb.Type.

// These bit flags can be used to query on the
// common properties of types.
const (
	flagIsIntegral = int(Flag_ISINTEGRAL)
	flagIsUnsigned = int(Flag_ISUNSIGNED)
	flagIsFloat    = int(Flag_ISFLOAT)
	flagIsQuoted   = int(Flag_ISQUOTED)
	flagIsText     = int(Flag_ISTEXT)
	flagIsBinary   = int(Flag_ISBINARY)
)

// IsIntegral returns true if Type is an integral
// (signed/unsigned) that can be represented using
// up to 64 binary bits.
func IsIntegral(t Type) bool {
	return int(t)&flagIsIntegral == flagIsIntegral
}

// IsSigned returns true if Type is a signed integral.
func IsSigned(t Type) bool {
	return int(t)&(flagIsIntegral|flagIsUnsigned) == flagIsIntegral
}

// IsUnsigned returns true if Type is an unsigned integral.
// Caution: this is not the same as !IsSigned.
func IsUnsigned(t Type) bool {
	return int(t)&(flagIsIntegral|flagIsUnsigned) == flagIsIntegral|flagIsUnsigned
}

// IsFloat returns true is Type is a floating point.
func IsFloat(t Type) bool {
	return int(t)&flagIsFloat == flagIsFloat
}

// IsQuoted returns true if Type is a quoted text or binary.
func IsQuoted(t Type) bool {
	return int(t)&flagIsQuoted == flagIsQuoted
}

// IsText returns true if Type is a text.
func IsText(t Type) bool {
	return int(t)&flagIsText == flagIsText
}

// IsBinary returns true if Type is a binary.
func IsBinary(t Type) bool {
	return int(t)&flagIsBinary == flagIsBinary
}

// isNumber returns true if the type is any type of number.
func isNumber(t Type) bool {
	return IsIntegral(t) || IsFloat(t) || t == Decimal
}

// IsTemporal returns true if Value is time type.
func IsTemporal(t Type) bool {
	switch t {
	case Timestamp, Date, Time, Datetime:
		return true
	}
	return false
}

// Vitess data types. These are idiomatically
// named synonyms for the Type values.
const (
	Null       = Type_NULL_TYPE
	Int8       = Type_INT8
	Uint8      = Type_UINT8
	Int16      = Type_INT16
	Uint16     = Type_UINT16
	Int24      = Type_INT24
	Uint24     = Type_UINT24
	Int32      = Type_INT32
	Uint32     = Type_UINT32
	Int64      = Type_INT64
	Uint64     = Type_UINT64
	Float32    = Type_FLOAT32
	Float64    = Type_FLOAT64
	Timestamp  = Type_TIMESTAMP
	Date       = Type_DATE
	Time       = Type_TIME
	Datetime   = Type_DATETIME
	Year       = Type_YEAR
	Decimal    = Type_DECIMAL
	Text       = Type_TEXT
	Blob       = Type_BLOB
	VarChar    = Type_VARCHAR
	VarBinary  = Type_VARBINARY
	Char       = Type_CHAR
	Binary     = Type_BINARY
	Bit        = Type_BIT
	Enum       = Type_ENUM
	Set        = Type_SET
	Tuple      = Type_TUPLE
	Geometry   = Type_GEOMETRY
	TypeJSON   = Type_JSON
	Expression = Type_EXPRESSION
)

// bit-shift the mysql flags by two byte so we
// can merge them with the mysql or vitess types.
const (
	mysqlUnsigned = 32
	mysqlBinary   = 128
	mysqlEnum     = 256
	mysqlSet      = 2048
)

// If you add to this map, make sure you add a test case
// in tabletserver/endtoend.
var mysqlToType = map[int64]Type{
	1:   Int8,
	2:   Int16,
	3:   Int32,
	4:   Float32,
	5:   Float64,
	6:   Null,
	7:   Timestamp,
	8:   Int64,
	9:   Int24,
	10:  Date,
	11:  Time,
	12:  Datetime,
	13:  Year,
	16:  Bit,
	245: TypeJSON,
	246: Decimal,
	249: Text,
	250: Text,
	251: Text,
	252: Text,
	253: VarChar,
	254: Char,
	255: Geometry,
}

// modifyType modifies the vitess type based on the
// mysql flag. The function checks specific flags based
// on the type. This allows us to ignore stray flags
// that MySQL occasionally sets.
func modifyType(typ Type, flags int64) Type {
	switch typ {
	case Int8:
		if flags&mysqlUnsigned != 0 {
			return Uint8
		}
		return Int8
	case Int16:
		if flags&mysqlUnsigned != 0 {
			return Uint16
		}
		return Int16
	case Int32:
		if flags&mysqlUnsigned != 0 {
			return Uint32
		}
		return Int32
	case Int64:
		if flags&mysqlUnsigned != 0 {
			return Uint64
		}
		return Int64
	case Int24:
		if flags&mysqlUnsigned != 0 {
			return Uint24
		}
		return Int24
	case Text:
		if flags&mysqlBinary != 0 {
			return Blob
		}
		return Text
	case VarChar:
		if flags&mysqlBinary != 0 {
			return VarBinary
		}
		return VarChar
	case Char:
		if flags&mysqlBinary != 0 {
			return Binary
		}
		if flags&mysqlEnum != 0 {
			return Enum
		}
		if flags&mysqlSet != 0 {
			return Set
		}
		return Char
	}
	return typ
}

// MySQLToType computes the vitess type from mysql type and flags.
func MySQLToType(mysqlType, flags int64) (typ Type, err error) {
	result, ok := mysqlToType[mysqlType]
	if !ok {
		return 0, fmt.Errorf("unsupported type: %d", mysqlType)
	}
	return modifyType(result, flags), nil
}

// typeToMySQL is the reverse of mysqlToType.
var typeToMySQL = map[Type]struct {
	typ   int64
	flags int64
}{
	Int8:      {typ: 1},
	Uint8:     {typ: 1, flags: mysqlUnsigned},
	Int16:     {typ: 2},
	Uint16:    {typ: 2, flags: mysqlUnsigned},
	Int32:     {typ: 3},
	Uint32:    {typ: 3, flags: mysqlUnsigned},
	Float32:   {typ: 4},
	Float64:   {typ: 5},
	Null:      {typ: 6, flags: mysqlBinary},
	Timestamp: {typ: 7},
	Int64:     {typ: 8},
	Uint64:    {typ: 8, flags: mysqlUnsigned},
	Int24:     {typ: 9},
	Uint24:    {typ: 9, flags: mysqlUnsigned},
	Date:      {typ: 10, flags: mysqlBinary},
	Time:      {typ: 11, flags: mysqlBinary},
	Datetime:  {typ: 12, flags: mysqlBinary},
	Year:      {typ: 13, flags: mysqlUnsigned},
	Bit:       {typ: 16, flags: mysqlUnsigned},
	TypeJSON:  {typ: 245},
	Decimal:   {typ: 246},
	Text:      {typ: 252},
	Blob:      {typ: 252, flags: mysqlBinary},
	VarChar:   {typ: 253},
	VarBinary: {typ: 253, flags: mysqlBinary},
	Char:      {typ: 254},
	Binary:    {typ: 254, flags: mysqlBinary},
	Enum:      {typ: 254, flags: mysqlEnum},
	Set:       {typ: 254, flags: mysqlSet},
	Geometry:  {typ: 255},
}

// TypeToMySQL returns the equivalent mysql type and flag for a vitess type.
func TypeToMySQL(typ Type) (mysqlType, flags int64) {
	val := typeToMySQL[typ]
	return val.typ, val.flags
}
