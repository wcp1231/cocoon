package redis

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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

type RedisObject interface {
	// 返回原始 bytes 数据
	Raw() []byte
	// hunmen readable
	Pretty() string
	// raw text
	Text() string
}

type RedisSimpleString struct {
	String string
}

type RedisInteger struct {
	Integer int64
}

type RedisError struct {
	Error string
}

type RedisBulkString struct {
	Len    int64
	String string
}

type RedisArray struct {
	Len   int
	Items []RedisObject
}

func (r *RedisSimpleString) Raw() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(SIMPLE_STRING)
	buf.WriteString(r.String)
	buf.WriteString("\r\n")
	return buf.Bytes()
}
func (r *RedisSimpleString) Pretty() string {
	return fmt.Sprintf(`"%s"`, r.String)
}
func (r *RedisSimpleString) Text() string {
	return r.String
}

func (r *RedisInteger) Raw() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(INTEGER)
	buf.WriteString(strconv.FormatInt(r.Integer, 10))
	buf.WriteString("\r\n")
	return buf.Bytes()
}
func (r *RedisInteger) Pretty() string {
	return r.Text()
}
func (r *RedisInteger) Text() string {
	return strconv.FormatInt(r.Integer, 10)
}

func (r *RedisError) Raw() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(ERROR)
	buf.WriteString(r.Error)
	buf.WriteString("\r\n")
	return buf.Bytes()
}
func (r *RedisError) Pretty() string {
	return fmt.Sprintf("-%s", r.Error)
}
func (r *RedisError) Text() string {
	return r.Error
}

func (r *RedisBulkString) Raw() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(BULK_STRING)
	if r.Len == -1 {
		buf.WriteString("-1\r\n")
		return buf.Bytes()
	}
	buf.WriteString(fmt.Sprintf("%d\r\n%s\r\n", r.Len, r.String))
	return buf.Bytes()
}
func (r *RedisBulkString) Pretty() string {
	if r.Len < 0 {
		return "nil"
	}
	return fmt.Sprintf(`"%s"`, r.String)
}
func (r *RedisBulkString) Text() string {
	return r.String
}

func (r *RedisArray) Raw() []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(ARRAY)
	if r.Len == -1 {
		buf.WriteString("-1\r\n")
		return buf.Bytes()
	}
	buf.WriteString(fmt.Sprintf("%d\r\n", r.Len))
	for _, item := range r.Items {
		buf.Write(item.Raw())
	}
	return buf.Bytes()
}
func (r *RedisArray) Pretty() string {
	if r.Len < 0 {
		return "nil"
	}
	buf := bytes.Buffer{}
	buf.WriteByte('[')
	for i, item := range r.Items {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(item.Pretty())
	}
	buf.WriteByte(']')
	return buf.String()
}
func (r *RedisArray) Text() string {
	return ""
}

type RESPReader struct {
	*bufio.Reader
}

func NewReader(reader *bufio.Reader) *RESPReader {
	return &RESPReader{
		Reader: reader,
	}
}

func (r *RESPReader) ReadObject() (RedisObject, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case SIMPLE_STRING:
		return r.readSimpleString(line)
	case INTEGER:
		return r.readInteger(line)
	case ERROR:
		return r.readError(line)
	case BULK_STRING:
		return r.readBulkString(line)
	case ARRAY:
		return r.readArray(line)
	case EMPTY_LINE:
		return nil, nil
	default:
		nextline, _ := r.readLine()
		fmt.Printf("RESP type not match. \ncurrent %v. str %s\n nextline %v. str %s\n", line, string(line), nextline, string(nextline))
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
		fmt.Printf("RESP read %v. str %s\n", line, string(line))
		// Line was too short or \n wasn't preceded by \r.
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) getCount(line []byte) (int64, error) {
	end := bytes.IndexByte(line, '\r')
	return strconv.ParseInt(string(line[1:end]), 10, 64)
}

func (r *RESPReader) readSimpleString(line []byte) (RedisObject, error) {
	object := &RedisSimpleString{}
	object.String = string(line[1 : len(line)-2])
	return object, nil
}

func (r *RESPReader) readInteger(line []byte) (RedisObject, error) {
	object := &RedisInteger{}
	var err error
	object.Integer, err = strconv.ParseInt(string(line[1:len(line)-2]), 10, 64)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (r *RESPReader) readError(line []byte) (RedisObject, error) {
	object := &RedisError{}
	object.Error = string(line[1 : len(line)-2])
	return object, nil
}

func (r *RESPReader) readBulkString(line []byte) (RedisObject, error) {
	object := &RedisBulkString{}
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	object.Len = count
	if count == -1 {
		return object, nil
	}
	stringLine, err := r.readLine()
	if err != nil {
		return nil, err
	}
	object.String = string(stringLine[:len(stringLine)-2])
	return object, nil
}

func (r *RESPReader) readArray(line []byte) (RedisObject, error) {
	// Get number of array elements.
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	object := &RedisArray{}
	object.Len = int(count)
	// Read `count` number of RESP objects in the array.
	for i := 0; i < object.Len; i++ {
		item, err := r.ReadObject()
		if err != nil {
			return nil, err
		}
		if item == nil {
			continue
		}
		object.Items = append(object.Items, item)
	}
	return object, nil
}
