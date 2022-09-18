package proto

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	// NULL represents the NULL value.
	NULL = Value{}
	// DontEscape tells you if a character should not be escaped.
	DontEscape = byte(255)
	nullstr    = []byte("null")
)

// MakeTrusted makes a new Value based on the type.
// If the value is an integral, then val must be in its canonical
// form. This function should only be used if you know the value
// and type conform to the rules.  Every place this function is
// called, a comment is needed that explains why it's justified.
// Functions within this package are exempt.
func MakeTrusted(typ Type, val []byte) Value {
	if typ == Null {
		return NULL
	}
	return Value{Type: typ, Value: val}
}

// BuildValue builds a value from any go type. sqltype.Value is
// also allowed.
func BuildValue(goval interface{}, typ Type) (v Value, err error) {
	// Look for the most common types first.
	switch goval := goval.(type) {
	case nil:
		// no op
	case []byte:
		switch typ {
		case Timestamp, Date, Datetime, Time:
			v = MakeTrusted(typ, goval)
		default:
			v = MakeTrusted(VarBinary, goval)
		}
	case int64:
		v = MakeTrusted(Int64, strconv.AppendInt(nil, int64(goval), 10))
	case uint64:
		v = MakeTrusted(Uint64, strconv.AppendUint(nil, uint64(goval), 10))
	case float64:
		v = MakeTrusted(Float64, strconv.AppendFloat(nil, goval, 'f', -1, 64))
	case int:
		v = MakeTrusted(Int64, strconv.AppendInt(nil, int64(goval), 10))
	case int8:
		v = MakeTrusted(Int8, strconv.AppendInt(nil, int64(goval), 10))
	case int16:
		v = MakeTrusted(Int16, strconv.AppendInt(nil, int64(goval), 10))
	case int32:
		v = MakeTrusted(Int32, strconv.AppendInt(nil, int64(goval), 10))
	case uint:
		v = MakeTrusted(Uint64, strconv.AppendUint(nil, uint64(goval), 10))
	case uint8:
		v = MakeTrusted(Uint8, strconv.AppendUint(nil, uint64(goval), 10))
	case uint16:
		v = MakeTrusted(Uint16, strconv.AppendUint(nil, uint64(goval), 10))
	case uint32:
		v = MakeTrusted(Uint32, strconv.AppendUint(nil, uint64(goval), 10))
	case float32:
		v = MakeTrusted(Float32, strconv.AppendFloat(nil, float64(goval), 'f', -1, 32))
	case string:
		switch typ {
		case Decimal, Text, Blob, VarChar, Char,
			Bit, Enum, Set, Geometry, TypeJSON:
			v = MakeTrusted(typ, []byte(goval))
		default:
			v = MakeTrusted(VarBinary, []byte(goval))
		}
	case time.Time:
		v = MakeTrusted(Datetime, []byte(goval.Format("2006-01-02 15:04:05")))
	case Value:
		v = goval
	case *BindVariable:
		return ValueFromBytes(goval.Type, goval.Value)
	default:
		return v, fmt.Errorf("unexpected type %T: %v", goval, goval)
	}
	return v, nil
}

// ValueFromBytes builds a Value using typ and val. It ensures that val
// matches the requested type. If type is an integral it's converted to
// a canonical form. Otherwise, the original representation is preserved.
func ValueFromBytes(typ Type, val []byte) (v Value, err error) {
	switch {
	case IsSigned(typ):
		signed, err := strconv.ParseInt(string(val), 0, 64)
		if err != nil {
			return NULL, err
		}
		v = MakeTrusted(typ, strconv.AppendInt(nil, signed, 10))
	case IsUnsigned(typ):
		unsigned, err := strconv.ParseUint(string(val), 0, 64)
		if err != nil {
			return NULL, err
		}
		v = MakeTrusted(typ, strconv.AppendUint(nil, unsigned, 10))
	case typ == Tuple:
		return NULL, errors.New("tuple not allowed for ValueFromBytes")
	case IsFloat(typ) || typ == Decimal:
		_, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return NULL, err
		}
		// After verification, we preserve the original representation.
		fallthrough
	default:
		v = MakeTrusted(typ, val)
	}
	return v, nil
}

func ParseMySQLValues(buf *Buffer, typ Type) (interface{}, error) {
	switch typ {
	case Null:
		return nil, nil
	case Int8, Uint8:
		return buf.ReadU8()
	case Uint16:
		return buf.ReadU16()
	case Int16, Year:
		val, err := buf.ReadU16()
		if err != nil {
			return nil, err
		}
		return int16(val), nil
	case Uint24, Uint32:
		return buf.ReadU32()
	case Int24, Int32:
		val, err := buf.ReadU32()
		if err != nil {
			return nil, err
		}
		return int32(val), nil
	case Float32:
		val, err := buf.ReadU32()
		if err != nil {
			return nil, err
		}
		return math.Float32frombits(val), nil
	case Uint64:
		return buf.ReadU64()
	case Int64:
		val, err := buf.ReadU64()
		if err != nil {
			return nil, err
		}
		return int64(val), nil
	case Float64:
		val, err := buf.ReadU64()
		if err != nil {
			return nil, err
		}
		return math.Float64frombits(val), nil
	case Timestamp, Date, Datetime:
		var out []byte

		size, err := buf.ReadU8()
		if err != nil {
			return nil, err
		}
		switch size {
		case 0x00:
			out = append(out, ' ')
		case 0x0b:
			year, err := buf.ReadU16()
			if err != nil {
				return nil, err
			}

			month, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			day, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			hour, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			minute, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			second, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			microSecond, err := buf.ReadU32()
			if err != nil {
				return nil, err
			}

			val := fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d.%06d",
				year, month, day, hour, minute, second, microSecond)
			out = []byte(val)
			return out, nil
		case 0x07:
			year, err := buf.ReadU16()
			if err != nil {
				return nil, err
			}

			month, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			day, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			hour, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			minute, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			second, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			val := fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d",
				year, month, day, hour, minute, second)
			out = []byte(val)
			return out, nil
		case 0x04:
			year, err := buf.ReadU16()
			if err != nil {
				return nil, err
			}

			month, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			day, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}
			val := fmt.Sprintf("%4d-%02d-%02d", year, month, day)
			out = []byte(val)
			return out, nil
		default:
			return nil, fmt.Errorf("datetime.error")
		}
	case Time:
		var out []byte

		size, err := buf.ReadU8()
		if err != nil {
			return nil, err
		}
		switch size {
		case 0x00:
			copy(out, "00:00:00")
		case 0x0c:
			isNegative, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			days, err := buf.ReadU32()
			if err != nil {
				return nil, err
			}

			hour, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			hours := uint32(hour) + days*uint32(24)

			minute, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			second, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			microSecond, err := buf.ReadU32()
			if err != nil {
				return nil, err
			}

			var val string
			if isNegative == 0x01 {
				val = fmt.Sprintf("-%02d:%02d:%02d.%06d", hours, minute, second, microSecond)
			} else {
				val = fmt.Sprintf("%02d:%02d:%02d.%06d", hours, minute, second, microSecond)
			}
			out = []byte(val)
			return out, nil
		case 0x08:
			isNegative, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			days, err := buf.ReadU32()
			if err != nil {
				return nil, err
			}

			hour, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			hours := uint32(hour) + days*uint32(24)

			minute, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			second, err := buf.ReadU8()
			if err != nil {
				return nil, err
			}

			var val string
			if isNegative == 0x01 {
				val = fmt.Sprintf("-%02d:%02d:%02d", hours, minute, second)
			} else {
				val = fmt.Sprintf("%02d:%02d:%02d", hours, minute, second)
			}
			out = []byte(val)
			return out, nil
		default:
			return nil, fmt.Errorf("time.error")
		}
	case Decimal, Text, Blob, VarChar, Char,
		Bit, Enum, Set, Geometry, TypeJSON:
		return buf.ReadLenEncodeString()
	case VarBinary, Binary:
		return buf.ReadLenEncodeBytes()
	default:
		return nil, fmt.Errorf("type.unhandle.error")
	}
	return nil, fmt.Errorf("type.unhandle.error")
}

func WriteMySQLValues(buf *Buffer, value Value) error {
	switch value.Type {
	case Null:
		return nil
	case Int8, Uint8:
		buf.WriteU8(value.Value[0])
		return nil
	case Uint16:
		n, err := strconv.ParseUint(string(value.Value), 10, 16)
		if err != nil {
			return err
		}
		buf.WriteU16(uint16(n))
		return nil
	case Int16, Year:
		n, err := strconv.ParseInt(string(value.Value), 10, 16)
		if err != nil {
			return err
		}
		buf.WriteU16(uint16(n))
		return nil
	case Uint24, Uint32:
		n, err := strconv.ParseUint(string(value.Value), 10, 32)
		if err != nil {
			return err
		}
		buf.WriteU32(uint32(n))
		return nil
	case Int24, Int32:
		n, err := strconv.ParseInt(string(value.Value), 10, 32)
		if err != nil {
			return err
		}
		buf.WriteU32(uint32(n))
		return nil
	case Float32:
		n, err := strconv.ParseFloat(string(value.Value), 32)
		if err != nil {
			return err
		}
		buf.WriteU32(math.Float32bits(float32(n)))
		return nil
	case Uint64:
		n, err := strconv.ParseUint(string(value.Value), 10, 64)
		if err != nil {
			return err
		}
		buf.WriteU64(n)
		return nil
	case Int64:
		n, err := strconv.ParseInt(string(value.Value), 10, 64)
		if err != nil {
			return err
		}
		buf.WriteU64(uint64(n))
		return nil
	case Float64:
		n, err := strconv.ParseFloat(string(value.Value), 64)
		if err != nil {
			return err
		}
		buf.WriteU64(math.Float64bits(n))
		return nil
	case Timestamp, Date, Datetime:
		var timestamp time.Time
		var err error
		var byteLen uint8 = 0
		if bytes.ContainsRune(value.Value, '.') {
			byteLen = 11
			timestamp, err = time.Parse("2006-01-02 15:04:05.000000", string(value.Value))
			if err != nil {
				return err
			}
		} else if bytes.ContainsRune(value.Value, ' ') {
			byteLen = 7
			timestamp, err = time.Parse("2006-01-02 15:04:05", string(value.Value))
			if err != nil {
				return err
			}
		} else if bytes.ContainsRune(value.Value, '-') {
			byteLen = 4
			timestamp, err = time.Parse("2006-01-02", string(value.Value))
			if err != nil {
				return err
			}
		}

		buf.WriteU8(byteLen)
		if byteLen >= 4 {
			buf.WriteU16(uint16(timestamp.Year()))
			buf.WriteU8(uint8(timestamp.Month()))
			buf.WriteU8(uint8(timestamp.Day()))
		}
		if byteLen >= 7 {
			buf.WriteU8(uint8(timestamp.Hour()))
			buf.WriteU8(uint8(timestamp.Minute()))
			buf.WriteU8(uint8(timestamp.Second()))
		}
		if byteLen >= 11 {
			buf.WriteU32(uint32(timestamp.Nanosecond() / int(time.Microsecond)))
		}
		return nil
	case Time:
		negative := value.Value[0] == '-'
		var byteLen uint8 = 0
		timeStr := string(value.Value)
		if timeStr == "_00:00:00_" {
			byteLen = 0
		} else if bytes.ContainsRune(value.Value, '.') {
			byteLen = 12
		} else if bytes.ContainsRune(value.Value, ':') {
			byteLen = 8
		}

		if byteLen == 0 {
			return nil
		}

		terms := strings.FieldsFunc(timeStr, func(r rune) bool {
			return r == ':' || r == '.'
		})

		hours, err := strconv.ParseInt(terms[0], 10, 32)
		if err != nil {
			return err
		}
		if negative {
			hours = -hours
		}
		days := hours / 24
		hours = hours % 24

		minutes, err := strconv.ParseInt(terms[1], 10, 8)
		if err != nil {
			return err
		}
		seconds, err := strconv.ParseInt(terms[2], 10, 8)
		if err != nil {
			return err
		}

		if byteLen == 8 {
			buf.WriteU8(byteLen)
			if negative {
				buf.WriteU8(1)
			} else {
				buf.WriteU8(0)
			}
			buf.WriteU32(uint32(days))
			buf.WriteU8(uint8(hours))
			buf.WriteU8(uint8(minutes))
			buf.WriteU8(uint8(seconds))
			return nil
		}

		microSeconds, err := strconv.ParseInt(terms[3], 10, 32)
		if err != nil {
			return err
		}

		if byteLen == 12 {
			buf.WriteU8(byteLen)
			if negative {
				buf.WriteU8(1)
			} else {
				buf.WriteU8(0)
			}
			buf.WriteU32(uint32(days))
			buf.WriteU8(uint8(hours))
			buf.WriteU8(uint8(minutes))
			buf.WriteU8(uint8(seconds))
			buf.WriteU32(uint32(microSeconds))
			return nil
		}

		return fmt.Errorf("time.error")

	case Decimal, Text, Blob, VarChar, Char,
		Bit, Enum, Set, Geometry, TypeJSON:
		buf.WriteLenEncodeBytes(value.Value)
		return nil
	case VarBinary, Binary:
		buf.WriteLenEncodeBytes(value.Value)
		return nil
	default:
		return fmt.Errorf("type.unhandle.error")
	}
	return fmt.Errorf("type.unhandle.error")
}
