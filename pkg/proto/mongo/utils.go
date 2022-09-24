package mongo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net"
)

func MustReadInt32(r io.Reader) (n int32, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}
func ReadInt32(r io.Reader) (n int32, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}
func WriteInt32(w io.Writer, n int32) error {
	return binary.Write(w, binary.LittleEndian, n)
}
func ReadUint32(r io.Reader) (n uint32, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}
func WriteUint32(w io.Writer, n uint32) error {
	return binary.Write(w, binary.LittleEndian, n)
}
func ReadUint8(r io.Reader) (n uint8, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}
func WriteUint8(w io.Writer, n uint8) error {
	return binary.Write(w, binary.LittleEndian, n)
}

func ReadInt64(r io.Reader) (int64, error) {
	var n int64
	err := binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return 0, err
	}
	return n, nil
}
func WriteInt64(w io.Writer, n int64) error {
	return binary.Write(w, binary.LittleEndian, n)
}

func ReadBytes(r io.Reader, n int) []byte {
	b := make([]byte, n)
	_, err := r.Read(b)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}
	return b
}

func ReadCString(r io.Reader) (string, error) {
	var b []byte
	var one = make([]byte, 1)
	for {
		_, err := r.Read(one)
		if err != nil {
			return "", err
		}
		if one[0] == '\x00' {
			break
		}
		b = append(b, one[0])
	}
	return string(b), nil
}
func WriteCString(w io.Writer, str string) error {
	_, err := w.Write(append([]byte(str), '\x00'))
	return err
}

func ReadOne(r io.Reader) ([]byte, error) {
	docLen, err := ReadInt32(r)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	buf := make([]byte, int(docLen))
	binary.LittleEndian.PutUint32(buf, uint32(docLen))
	if _, err := io.ReadFull(r, buf[4:]); err != nil {
		panic(err)
	}
	return buf, nil
}

func ReadDocument(r io.Reader) (m bson.D, err error) {
	one, err := ReadOne(r)
	if err != nil {
		return nil, err
	}
	if one != nil {
		err := bson.Unmarshal(one, &m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
func WriteDocument(w io.Writer, m bson.D) error {
	if m == nil {
		return nil
	}
	bs, err := bson.Marshal(m)
	if err != nil {
		return err
	}
	// bson marshal 时已经是最终的结果
	_, err = w.Write(bs)
	return err
}

func ReadDocuments(r io.Reader) (ms []bson.D, err error) {
	for {
		m, err := ReadDocument(r)
		if err != nil {
			return nil, err
		}
		if m == nil {
			break
		}
		ms = append(ms, m)
	}
	return
}
func WriteDocuments(w io.Writer, ms []bson.D) error {
	for _, m := range ms {
		err := WriteDocument(w, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func ToJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("{\"error\":%s}", err.Error())
	}
	return string(b)
}

func ToJsonB(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte(fmt.Sprintf("{\"error\":%s}", err.Error()))
	}
	return b
}

func isClosedErr(err error) bool {
	if e, ok := err.(*net.OpError); ok {
		if e.Err.Error() == "use of closed network connection" {
			return true
		}
	}
	return false
}
