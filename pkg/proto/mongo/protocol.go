package mongo

import (
	"bufio"
	"encoding/binary"
	"go.mongodb.org/mongo-driver/bson"
	"io"
)

type Opcode int32

const (
	OP_REPLY        = 1
	OP_UPDATE       = 2001
	OP_INSERT       = 2002
	OP_RESERVED     = 2003
	OP_QUERY        = 2004
	OP_GET_MORE     = 2005
	OP_DELETE       = 2006
	OP_KILL_CURSORS = 2007
	OP_COMPRESSED   = 2012
	OP_MSG          = 2013
)

type MessageHeader struct {
	MessageLength int32
	RequestID     int32
	ResponseTo    int32
	OpCode        int32
}

func (h *MessageHeader) Valid() bool {
	if h.OpCode == 1 || h.OpCode == 1000 {
		return true
	}
	if 2001 <= h.OpCode && h.OpCode <= 2013 {
		return true
	}
	return false
}

func unpackMessageHeader(r io.Reader) (header MessageHeader, err error) {
	err = binary.Read(r, binary.LittleEndian, &header)
	return
}
func packMessageHeader(w io.Writer, header MessageHeader) error {
	return binary.Write(w, binary.LittleEndian, header)
}

type OpUpdateMessage struct {
	Zero         int32
	FullCollName string
	Flags        int32
	Selector     bson.D
	Update       bson.D
}

func unpackUpdateMessage(r io.Reader) (msg *OpUpdateMessage, err error) {
	msg = &OpUpdateMessage{}
	_, _ = MustReadInt32(r)
	msg.FullCollName, err = ReadCString(r)
	if err != nil {
		return nil, err
	}
	msg.Flags, err = MustReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.Selector, err = ReadDocument(r)
	if err != nil {
		return nil, err
	}
	msg.Update, err = ReadDocument(r)
	return
}
func packUpdateMessage(w io.Writer, msg *OpUpdateMessage) error {
	err := WriteInt32(w, 0)
	if err != nil {
		return err
	}
	err = WriteCString(w, msg.FullCollName)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.Flags)
	if err != nil {
		return err
	}
	err = WriteDocument(w, msg.Selector)
	if err != nil {
		return err
	}
	return WriteDocument(w, msg.Selector)
}

type OpInsertMessage struct {
	Flags        int32
	FullCollName string
	Documents    []bson.D
}

func unpackInsertMessage(r io.Reader) (msg *OpInsertMessage, err error) {
	msg = &OpInsertMessage{}
	msg.Flags, err = MustReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.FullCollName, err = ReadCString(r)
	if err != nil {
		return nil, err
	}
	msg.Documents, err = ReadDocuments(r)
	return
}
func packInsertMessage(w io.Writer, msg *OpInsertMessage) error {
	err := WriteInt32(w, msg.Flags)
	if err != nil {
		return err
	}
	err = WriteCString(w, msg.FullCollName)
	if err != nil {
		return err
	}
	return WriteDocuments(w, msg.Documents)
}

type OpQueryMessage struct {
	Flags          int32
	FullCollName   string
	NumberToSkip   int32
	NumberToReturn int32
	Query          bson.D
	ReturnFields   bson.D
}

func unpackQueryMessage(r io.Reader) (msg *OpQueryMessage, err error) {
	msg = &OpQueryMessage{}
	msg.Flags, err = MustReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.FullCollName, err = ReadCString(r)
	if err != nil {
		return nil, err
	}
	msg.NumberToSkip, err = MustReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.NumberToReturn, err = MustReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.Query, err = ReadDocument(r)
	if err != nil {
		return nil, err
	}
	msg.ReturnFields, err = ReadDocument(r)
	return
}
func packQueryMessage(w io.Writer, msg *OpQueryMessage) error {
	err := WriteInt32(w, msg.Flags)
	if err != nil {
		return err
	}
	err = WriteCString(w, msg.FullCollName)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.NumberToSkip)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.NumberToReturn)
	if err != nil {
		return err
	}
	err = WriteDocument(w, msg.Query)
	if err != nil {
		return err
	}
	return WriteDocument(w, msg.ReturnFields)
}

type OpGetMoreMessage struct {
	FullCollName   string
	NumberToReturn int32
	CursorID       int64
}

func unpackGetMoreMessage(r io.Reader) (msg *OpGetMoreMessage, err error) {
	msg = &OpGetMoreMessage{}
	_, _ = MustReadInt32(r)
	msg.FullCollName, err = ReadCString(r)
	if err != nil {
		return nil, err
	}
	msg.NumberToReturn, err = MustReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.CursorID, err = ReadInt64(r)
	return
}
func packGetMoreMessage(w io.Writer, msg *OpGetMoreMessage) error {
	err := WriteInt32(w, 0)
	if err != nil {
		return err
	}
	err = WriteCString(w, msg.FullCollName)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.NumberToReturn)
	if err != nil {
		return err
	}
	return WriteInt64(w, msg.CursorID)
}

type OpDeleteMessage struct {
	FullCollName string
	Flags        int32
	Selector     bson.D
}

func unpackDeleteMessage(r io.Reader) (msg *OpDeleteMessage, err error) {
	msg = &OpDeleteMessage{}
	_, _ = MustReadInt32(r)
	msg.FullCollName, err = ReadCString(r)
	if err != nil {
		return nil, err
	}
	msg.Flags, err = MustReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.Selector, err = ReadDocument(r)
	return
}
func packDeleteMessage(w io.Writer, msg *OpDeleteMessage) error {
	err := WriteInt32(w, 0)
	if err != nil {
		return err
	}
	err = WriteCString(w, msg.FullCollName)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.Flags)
	if err != nil {
		return err
	}
	return WriteDocument(w, msg.Selector)
}

type OpKillCursorsMessage struct {
	NumberOfCursorIDs int32
	CursorIDs         []int64
}

func unpackKillCursorsMessage(r io.Reader) (msg *OpKillCursorsMessage, err error) {
	msg = &OpKillCursorsMessage{}
	_, _ = MustReadInt32(r)
	msg.NumberOfCursorIDs, err = ReadInt32(r)
	if err != nil {
		return nil, err
	}
	var cursorIDs []int64
	for {
		n, err := ReadInt64(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		cursorIDs = append(cursorIDs, n)
	}
	msg.CursorIDs = cursorIDs
	return
}
func packKillCursorsMessage(w io.Writer, msg *OpKillCursorsMessage) error {
	err := WriteInt32(w, 0)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.NumberOfCursorIDs)
	if err != nil {
		return err
	}
	for _, id := range msg.CursorIDs {
		err = WriteInt64(w, id)
		if err != nil {
			return err
		}
	}
	return nil
}

type OPCompressedMessage struct {
	OriginalOpcode    int32
	UncompressedSize  int32
	CompressorId      uint8
	CompressedMessage []byte
}

func unpackCompressedMessage(r io.Reader) (msg *OPCompressedMessage, err error) {
	msg = &OPCompressedMessage{}
	_, _ = MustReadInt32(r)
	msg.OriginalOpcode, err = ReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.UncompressedSize, err = ReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.CompressorId, err = ReadUint8(r)
	if err != nil {
		return nil, err
	}

	switch msg.CompressorId {
	case 0: // noop
	case 1: // snappy
	case 2: // zlib
	case 3: // zstd
	}
	return
}
func packCompressedMessage(w io.Writer, msg *OPCompressedMessage) error {
	err := WriteInt32(w, msg.OriginalOpcode)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.UncompressedSize)
	if err != nil {
		return err
	}
	err = WriteUint8(w, msg.CompressorId)
	return nil
}

type OpMsgMessage struct {
	Flags    uint32
	Sections []Section
	Checksum uint32
}
type Section struct {
	Kind uint8
	Body bson.D

	Identifier string
	Objects    []bson.D
}

func unpackMsgMessage(r io.Reader) (msg *OpMsgMessage, err error) {
	msg = &OpMsgMessage{}
	msg.Flags, err = ReadUint32(r)
	if err != nil {
		return nil, err
	}
	msg.Sections = []Section{}
	for {
		section := Section{}
		section.Kind, err = ReadUint8(r)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return nil, err
		}
		if section.Kind == 0 {
			section.Body, err = ReadDocument(r)
			if err != nil {
				return nil, err
			}
		} else if section.Kind == 1 {
			sectionSize, err := MustReadInt32(r)
			if err != nil {
				return nil, err
			}
			limitReader := io.LimitReader(r, int64(sectionSize))
			section.Identifier, err = ReadCString(limitReader)
			if err != nil {
				return nil, err
			}
			section.Objects, err = ReadDocuments(limitReader)
			if err != nil {
				return nil, err
			}
		}
		msg.Sections = append(msg.Sections, section)
	}

	if msg.Flags&1 != 0 {
		msg.Checksum, _ = ReadUint32(r)
	}

	return
}
func packMsgMessage(w io.Writer, msg *OpMsgMessage) error {
	err := WriteUint32(w, msg.Flags)
	if err != nil {
		return err
	}
	for _, section := range msg.Sections {
		err = WriteUint8(w, section.Kind)
		if err != nil {
			return err
		}
		if section.Kind == 0 {
			err = WriteDocument(w, section.Body)
			if err != nil {
				return err
			}
		} else if section.Kind == 1 {
			buf := bufio.NewWriter(w)
			err = WriteCString(buf, section.Identifier)
			if err != nil {
				return err
			}
			err = WriteDocuments(buf, section.Objects)
			if err != nil {
				return err
			}
			err = WriteInt32(w, int32(buf.Buffered()+4))
			if err != nil {
				return err
			}
			err = buf.Flush()
			if err != nil {
				return err
			}
		}
	}

	if msg.Flags&1 != 0 {
		return WriteUint32(w, msg.Checksum)
	}
	return nil
}

type OpReplyMessage struct {
	Flags          int32
	CursorID       int64
	StartingFrom   int32
	NumberReturned int32
	Documents      []bson.D
}

func unpackReplyMessage(r io.Reader) (msg *OpReplyMessage, err error) {
	msg = &OpReplyMessage{}
	msg.Flags, err = ReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.CursorID, err = ReadInt64(r)
	if err != nil {
		return nil, err
	}
	msg.StartingFrom, err = ReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.NumberReturned, err = ReadInt32(r)
	if err != nil {
		return nil, err
	}
	msg.Documents, err = ReadDocuments(r)
	return
}
func packReplyMessage(w io.Writer, msg *OpReplyMessage) error {
	err := WriteInt32(w, msg.Flags)
	if err != nil {
		return err
	}
	err = WriteInt64(w, msg.CursorID)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.StartingFrom)
	if err != nil {
		return err
	}
	err = WriteInt32(w, msg.NumberReturned)
	if err != nil {
		return err
	}
	return WriteDocuments(w, msg.Documents)
}
