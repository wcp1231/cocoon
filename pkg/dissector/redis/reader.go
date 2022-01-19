package redis

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/api"
	"errors"
	"io"
	"strconv"
)

const (
	SIMPLE_STRING = '+'
	BULK_STRING   = '$'
	INTEGER       = ':'
	ARRAY         = '*'
	ERROR         = '-'
	EMPTY_LINE    = '\r'
)

var (
	ErrInvalidSyntax = errors.New("resp: invalid syntax")
)

type RedisObject struct {
	Type  int32
	Count int
	Data  string
	Array []*RedisObject
	Raw   []byte
}

func NewBulkString() *RedisObject {
	object := &RedisObject{}
	object.Type = BULK_STRING
	object.Count = 1
	return object
}

func NewArray(count int) *RedisObject {
	object := &RedisObject{}
	object.Type = ARRAY
	object.Count = count
	object.Array = make([]*RedisObject, count)
	return object
}

func NewStringOrIntegerOrError(t int32, line []byte) *RedisObject {
	object := &RedisObject{}
	object.Type = t
	object.Data = string(line[1:])
	object.Raw = line
	return object
}

type RESPReader struct {
	*bufio.Reader
}

func NewReader(reader api.TcpReader) *RESPReader {
	return &RESPReader{
		Reader: reader.BufferReader(),
	}
}

func (r *RESPReader) ReadObject() (*RedisObject, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case SIMPLE_STRING:
		return NewStringOrIntegerOrError(SIMPLE_STRING, line), nil
	case INTEGER:
		return NewStringOrIntegerOrError(INTEGER, line), nil
	case ERROR:
		return NewStringOrIntegerOrError(ERROR, line), nil
	case BULK_STRING:
		return r.readBulkString(line)
	case ARRAY:
		return r.readArray(line)
	case EMPTY_LINE:
		return nil, nil
	default:
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) readLine() (line []byte, err error) {
	line, err = r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(line) > 1 && line[len(line)-2] == '\r' {
		return line, nil
	} else {
		// Line was too short or \n wasn't preceded by \r.
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) readBulkString(line []byte) (*RedisObject, error) {
	object := NewBulkString()
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	if count == -1 {
		object.Raw = make([]byte, len(line))
		copy(object.Raw, line)
		return object, nil
	}
	raw := make([]byte, len(line)+count+2)
	copy(raw, line)
	data := make([]byte, count+2)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	copy(raw[len(line):], data)
	object.Data = string(data)
	object.Raw = raw
	return object, nil
}

func (r *RESPReader) getCount(line []byte) (int, error) {
	end := bytes.IndexByte(line, '\r')
	return strconv.Atoi(string(line[1:end]))
}

func (r *RESPReader) readArray(line []byte) (*RedisObject, error) {
	// Get number of array elements.
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	object := NewArray(count)
	// Read `count` number of RESP objects in the array.
	for i := 0; i < count; i++ {
		item, err := r.ReadObject()
		if err != nil {
			return nil, err
		}
		if item == nil {
			continue
		}
		line = append(line, item.Raw...)
		object.Array[i] = item
	}
	object.Raw = make([]byte, len(line))
	copy(object.Raw, line)
	return object, nil
}
