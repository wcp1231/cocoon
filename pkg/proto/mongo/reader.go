package mongo

import (
	"bytes"
	"cocoon/pkg/model/common"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
)

type Parser struct{}

func ParseQuery(message *MongoMessage, r io.Reader) error {
	msg, err := unpackQueryMessage(r)
	if err != nil {
		return err
	}

	message.SetMongoMessage(msg)
	message.SetOpType("Query")
	return nil
}
func DumpQuery(message interface{}, w io.Writer) error {
	msg := message.(*OpQueryMessage)
	return packQueryMessage(w, msg)
}

func ParseInsert(message *MongoMessage, r io.Reader) error {
	msg, err := unpackInsertMessage(r)
	if err != nil {
		return err
	}

	message.SetOpType("Insert")
	message.SetMongoMessage(msg)
	return nil
}
func DumpInsert(message interface{}, w io.Writer) error {
	msg := message.(*OpInsertMessage)
	return packInsertMessage(w, msg)
}

func ParseUpdate(message *MongoMessage, r io.Reader) error {
	msg, err := unpackUpdateMessage(r)
	if err != nil {
		return err
	}

	message.SetOpType("Update")
	message.SetMongoMessage(msg)
	return nil
}
func DumpUpdate(message interface{}, w io.Writer) error {
	msg := message.(*OpUpdateMessage)
	return packUpdateMessage(w, msg)
}

func ParseGetMore(message *MongoMessage, r io.Reader) error {
	msg, err := unpackGetMoreMessage(r)
	if err != nil {
		return err
	}

	message.SetOpType("GetMore")
	message.SetMongoMessage(msg)
	return nil
}
func DumpGetMore(message interface{}, w io.Writer) error {
	msg := message.(*OpGetMoreMessage)
	return packGetMoreMessage(w, msg)
}

func ParseDelete(message *MongoMessage, r io.Reader) error {
	msg, err := unpackDeleteMessage(r)
	if err != nil {
		return err
	}

	message.SetOpType("Delete")
	message.SetMongoMessage(msg)
	return nil
}
func DumpDelete(message interface{}, w io.Writer) error {
	msg := message.(*OpDeleteMessage)
	return packDeleteMessage(w, msg)
}

func ParseKillCursors(message *MongoMessage, r io.Reader) error {
	msg, err := unpackKillCursorsMessage(r)
	if err != nil {
		return err
	}

	message.SetOpType("KillCursors")
	message.SetMongoMessage(msg)
	return nil
}
func DumpKillCursors(message interface{}, w io.Writer) error {
	msg := message.(*OpKillCursorsMessage)
	return packKillCursorsMessage(w, msg)
}

func ParseReply(message *MongoMessage, r io.Reader) error {
	msg, err := unpackReplyMessage(r)
	if err != nil {
		return err
	}

	message.SetOpType("Reply")
	message.SetMongoMessage(msg)
	return nil
}
func DumpReply(message interface{}, w io.Writer) error {
	msg := message.(*OpReplyMessage)
	return packReplyMessage(w, msg)
}

func ParseReserved(message *MongoMessage, r io.Reader) {}
func DumpReserved(message interface{}, r io.Reader)    {}

func ParseMsg(message *MongoMessage, r io.Reader) error {
	msg, err := unpackMsgMessage(r)
	if err != nil {
		return err
	}

	message.SetOpType("Msg")
	message.SetMongoMessage(msg)
	return nil
}
func DumpMsg(message interface{}, w io.Writer) error {
	msg := message.(*OpMsgMessage)
	return packMsgMessage(w, msg)
}

func readMsgHeader(r io.Reader) (*MessageHeader, error) {
	h := MessageHeader{}
	err := binary.Read(r, binary.LittleEndian, &h)

	if err != nil {
		return nil, err
	}
	return &h, nil
}

func Parse(r io.Reader) (common.Message, error) {
	result := NewMongoGenericMessage()
	header, err := unpackMessageHeader(r)
	if err != nil {
		if err != io.EOF {
			log.Printf("unexpected error:%v\n", err)
		}
		return nil, err
	}
	result.SetMongoHeader(header)
	rd := io.LimitReader(r, int64(header.MessageLength-4*4))
	switch header.OpCode {
	case OP_QUERY:
		err = ParseQuery(result, rd)
	case OP_INSERT:
		err = ParseInsert(result, rd)
	case OP_DELETE:
		err = ParseDelete(result, rd)
	case OP_UPDATE:
		err = ParseUpdate(result, rd)
	case OP_REPLY:
		err = ParseReply(result, rd)
	case OP_GET_MORE:
		err = ParseGetMore(result, rd)
	case OP_KILL_CURSORS:
		err = ParseKillCursors(result, rd)
	case OP_RESERVED:
		ParseReserved(result, rd)
	case OP_MSG:
		err = ParseMsg(result, rd)
	default:
		log.Printf("unknown OpCode: %d", header.OpCode)
		_, err = io.Copy(ioutil.Discard, rd)
		if err != nil {
			log.Printf("read failed: %v", err)
			return nil, err
		}
	}
	return result, err
}

func Dump(message common.Message) ([]byte, error) {
	buf := &bytes.Buffer{}
	mongoMsg := message.(*MongoMessage)
	header := mongoMsg.GetMongoHeader()
	err := packMessageHeader(buf, header)
	if err != nil {
		return nil, err
	}
	msg := mongoMsg.GetMongoMessage()
	switch header.OpCode {
	case OP_QUERY:
		err = DumpQuery(msg, buf)
	case OP_INSERT:
		err = DumpInsert(msg, buf)
	case OP_DELETE:
		err = DumpDelete(msg, buf)
	case OP_UPDATE:
		err = DumpUpdate(msg, buf)
	case OP_REPLY:
		err = DumpReply(msg, buf)
	case OP_GET_MORE:
		err = DumpGetMore(msg, buf)
	case OP_KILL_CURSORS:
		err = DumpKillCursors(msg, buf)
	case OP_RESERVED:
		DumpReserved(msg, buf)
	case OP_MSG:
		err = DumpMsg(msg, buf)
	}

	return buf.Bytes(), nil
}
