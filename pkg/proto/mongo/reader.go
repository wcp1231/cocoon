package mongo

import (
	"cocoon/pkg/model/common"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"log"
)

type Opcode int32

const (
	OP_REPLY                    = 1
	OP_MSG                      = 1000
	OP_UPDATE                   = 2001
	OP_INSERT                   = 2002
	OP_RESERVED                 = 2003
	OP_QUERY                    = 2004
	OP_GET_MORE                 = 2005
	OP_DELETE                   = 2006
	OP_KILL_CURSORS             = 2007
	OP_COMMAND_DEPRECATED       = 2008
	OP_COMMAND_REPLY_DEPRECATED = 2009
	OP_COMMAND                  = 2010
	OP_COMMAND_REPLY            = 2011
	OP_MSG_NEW                  = 2013
)

type MsgHeader struct {
	MessageLength int32
	RequestID     int32
	ResponseTo    int32
	OpCode        int32
}

func (h *MsgHeader) Valid() bool {
	if h.OpCode == 1 || h.OpCode == 1000 {
		return true
	}
	if 2001 <= h.OpCode && h.OpCode <= 2013 {
		return true
	}
	return false
}

type Parser struct{}

func ParseQuery(header MsgHeader, r io.Reader) common.Message {
	flag := MustReadInt32(r)
	fullCollectionName := ReadCString(r)
	numberToSkip := MustReadInt32(r)
	numberToReturn := MustReadInt32(r)
	query := ReadDocument(r)
	queryJson := ToJson(query)
	selector := ToJson(ReadDocument(r))

	result := NewMongoGenericMessage()
	result.SetOpType("Query")
	result.SetRequestId(header.RequestID)
	result.SetFlag(flag)
	result.SetCollection(fullCollectionName)
	result.SetNumberToSkip(numberToSkip)
	result.SetNumberToReturn(numberToReturn)
	result.SetQuery(queryJson)
	result.SetSelector(selector)

	queryMap := parseQuery(query)
	for k, v := range queryMap {
		result.Header[k] = v
	}

	return result
}

func ParseInsert(header MsgHeader, r io.Reader) common.Message {
	flag := MustReadInt32(r)
	fullCollectionName := ReadCString(r)
	docs := ReadDocuments(r)
	var docsStr string
	if len(docs) == 1 {
		docsStr = ToJson(docs[0])
	} else {
		docsStr = ToJson(docs)
	}
	//fmt.Printf("%s INSERT id:%d coll:%s flag:%b docs:%v\n",
	//	currentTime(), header.RequestID, fullCollectionName, flag, docsStr)

	result := NewMongoGenericMessage()
	result.SetOpType("Insert")
	result.SetRequestId(header.RequestID)
	result.SetFlag(flag)
	result.SetCollection(fullCollectionName)
	result.SetMongoBody(docsStr)
	return result
}

func ParseUpdate(header MsgHeader, r io.Reader) common.Message {
	_ = MustReadInt32(r)
	fullCollectionName := ReadCString(r)
	flag := MustReadInt32(r)
	selector := ToJson(ReadDocument(r))
	update := ToJson(ReadDocument(r))
	//fmt.Printf("%s UPDATE id:%d coll:%s flag:%b sel:%v update:%v\n",
	//	currentTime(), header.RequestID, fullCollectionName, flag, selector, update)

	result := NewMongoGenericMessage()
	result.SetOpType("Update")
	result.SetRequestId(header.RequestID)
	result.SetFlag(flag)
	result.SetCollection(fullCollectionName)
	result.SetSelector(selector)
	result.SetMongoBody(update)
	return result
}

func ParseGetMore(header MsgHeader, r io.Reader) common.Message {
	_ = MustReadInt32(r)
	fullCollectionName := ReadCString(r)
	numberToReturn := MustReadInt32(r)
	cursorID := ReadInt64(r)
	//fmt.Printf("%s GETMORE id:%d coll:%s toret:%d curID:%d\n",
	//	currentTime(), header.RequestID, fullCollectionName, numberToReturn, cursorID)

	result := NewMongoGenericMessage()
	result.SetOpType("GetMore")
	result.SetRequestId(header.RequestID)
	result.SetCollection(fullCollectionName)
	result.SetNumberToReturn(numberToReturn)
	result.SetCursorId(*cursorID)
	// TODO raw
	return result
}

func ParseDelete(header MsgHeader, r io.Reader) common.Message {
	_ = MustReadInt32(r)
	fullCollectionName := ReadCString(r)
	//flag := MustReadInt32(r)
	selector := ToJson(ReadDocument(r))
	//fmt.Printf("%s DELETE id:%d coll:%s flag:%b sel:%v \n",
	//	currentTime(), header.RequestID, fullCollectionName, flag, selector)
	result := NewMongoGenericMessage()
	result.SetOpType("Delete")
	result.SetRequestId(header.RequestID)
	result.SetCollection(fullCollectionName)
	result.SetSelector(selector)
	// TODO raw
	return result
}

func ParseKillCursors(header MsgHeader, r io.Reader) common.Message {
	_ = MustReadInt32(r)
	//numberOfCursorIDs := MustReadInt32(r)
	var cursorIDs []int64
	for {
		n := ReadInt64(r)
		if n != nil {
			cursorIDs = append(cursorIDs, *n)
			continue
		}
		break
	}
	//fmt.Printf("%s KILLCURSORS id:%d numCurID:%d curIDs:%s\n",
	//	currentTime(), header.RequestID, numberOfCursorIDs, cursorIDs)

	result := NewMongoGenericMessage()
	result.SetOpType("KillCursors")
	result.SetRequestId(header.RequestID)
	result.SetCursorIds(cursorIDs)
	// TODO raw
	return result
}

func ParseReply(header MsgHeader, r io.Reader) common.Message {
	flag := MustReadInt32(r)
	cursorID := ReadInt64(r)
	startingFrom := MustReadInt32(r)
	numberReturned := MustReadInt32(r)
	docs := ReadDocuments(r)
	var docsStr string
	if len(docs) == 1 {
		docsStr = ToJson(docs[0])
	} else {
		docsStr = ToJson(docs)
	}
	//fmt.Printf("%s REPLY to:%d flag:%b curID:%d from:%d reted:%d docs:%v\n",
	//	currentTime(),
	//	header.ResponseTo,
	//	flag,
	//	cursorID,
	//	startingFrom,
	//	numberReturned,
	//	docsStr,
	//)

	result := NewMongoGenericMessage()
	result.SetOpType("Reply")
	result.SetRequestId(header.RequestID)
	result.SetFlag(flag)
	result.SetCursorId(*cursorID)
	result.SetStartingFrom(startingFrom)
	result.SetNumberToReturn(numberReturned)
	result.SetMongoBody(docsStr)
	return result
}

func ParseMsg(header MsgHeader, r io.Reader) common.Message {
	msg := ReadCString(r)
	//fmt.Printf("%s MSG %d %s\n", currentTime(), header.RequestID, msg)
	result := NewMongoGenericMessage()
	result.SetOpType("Msg")
	result.SetRequestId(header.RequestID)
	result.SetMsg(msg)
	// TODO raw
	return result
}
func ParseReserved(header MsgHeader, r io.Reader) {
	fmt.Printf("%s RESERVED header:%v data:%v\n", currentTime(), header.RequestID, ToJson(header))
}

//func (self *Parser) ParseCommandDeprecated(header MsgHeader, r io.Reader) {
//	fmt.Printf("%s MsgHeader %v\n", currentTime(), ToJson(header))
//	// TODO: no document, current not understand
//	_, err := io.Copy(ioutil.Discard, r)
//	if err != nil {
//		fmt.Printf("read failed: %v", err)
//		return
//	}
//}
//func (self *Parser) ParseCommandReplyDeprecated(header MsgHeader, r io.Reader) {
//	fmt.Printf("%s MsgHeader %v\n", currentTime(), ToJson(header))
//	// TODO: no document, current not understand
//	_, err := io.Copy(ioutil.Discard, r)
//	if err != nil {
//		fmt.Printf("read failed: %v", err)
//		return
//	}
//}

func ParseCommand(header MsgHeader, r io.Reader) common.Message {
	database := ReadCString(r)
	commandName := ReadCString(r)
	metadata := ToJson(ReadDocument(r))
	commandArgs := ToJson(ReadDocument(r))
	inputDocs := ToJson(ReadDocuments(r))
	//fmt.Printf("%s COMMAND id:%v db:%v meta:%v cmd:%v args:%v docs %v\n",
	//	currentTime(),
	//	header.RequestID,
	//	database,
	//	metadata,
	//	commandName,
	//	commandArgs,
	//	inputDocs,
	//)

	result := NewMongoGenericMessage()
	result.SetOpType("Command")
	result.SetRequestId(header.RequestID)
	result.SetDatabase(database)
	result.SetCommand(commandName)
	result.SetMetadata(metadata)
	result.SetCommandArgs(commandArgs)
	result.SetMongoBody(inputDocs)
	return result
}

func ParseMsgNew(header MsgHeader, r io.Reader) common.Message {
	flags := ToJson(MustReadInt32(r))
	//fmt.Printf("%s MSG start id:%v flags: %v\n", currentTime(), header.RequestID, flags)
	var msgs []map[string]interface{}
	for {
		t := ReadBytes(r, 1)
		if t == nil {
			//fmt.Printf("%s MSG end id:%v \n",
			//	currentTime(),
			//	header.RequestID,
			//)
			break
		}
		switch t[0] {
		case 0: // body
			body := ReadDocument(r)
			//bodyJson := ToJson(body)
			checksum, _ := ReadUint32(r)
			//fmt.Printf("%s MSG id:%v type:0 body: %v checksum:%v\n",
			//	currentTime(),
			//	header.RequestID,
			//	bodyJson,
			//	checksum,
			//)
			item := map[string]interface{}{
				"body":     body,
				"checksum": checksum,
			}
			msgs = append(msgs, item)
		case 1:
			sectionSize := MustReadInt32(r)
			r1 := io.LimitReader(r, int64(sectionSize))
			documentSequenceIdentifier := ReadCString(r1)
			objects := ReadDocuments(r1)
			//objectsJson := ToJson(objects)
			//fmt.Printf("%s MSG id:%v type:1 documentSequenceIdentifier: %v objects:%v\n",
			//	currentTime(),
			//	header.RequestID,
			//	documentSequenceIdentifier,
			//	objectsJson,
			//)
			item := map[string]interface{}{
				"objects":                      objects,
				"document_sequence_identifier": documentSequenceIdentifier,
			}
			msgs = append(msgs, item)
		default:
			log.Panic(fmt.Sprint("unknown body kind:", t[0]))
		}
	}

	result := NewMongoGenericMessage()
	result.SetOpType("Msg")
	result.SetRequestId(header.RequestID)
	result.SetFlags(flags)
	result.SetMongoBody(ToJson(msgs))
	return result
}

func ParseCommandReply(header MsgHeader, r io.Reader) common.Message {
	metadata := ToJson(ReadDocument(r))
	commandReply := ToJson(ReadDocument(r))
	outputDocs := ToJson(ReadDocument(r))
	//fmt.Printf("%s COMMANDREPLY to:%d id:%v meta:%v cmdReply:%v outputDocs:%v\n",
	//	currentTime(), header.ResponseTo, header.RequestID, metadata, commandReply, outputDocs)

	result := NewMongoGenericMessage()
	result.SetOpType("CommandReply")
	result.SetRequestId(header.RequestID)
	result.SetResponseTo(header.ResponseTo)
	result.SetMetadata(metadata)
	result.SetCommandReply(commandReply)
	result.SetMongoBody(outputDocs)
	// TODO raw
	return result
}

func readMsgHeader(r io.Reader) (*MsgHeader, error) {
	h := MsgHeader{}
	err := binary.Read(r, binary.LittleEndian, &h)

	if err != nil {
		return nil, err
	}
	return &h, nil
}

func Parse(r io.Reader) (common.Message, error) {

	header := MsgHeader{}
	err := binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		if err != io.EOF {
			log.Printf("unexpected error:%v\n", err)
		}
		return nil, err
	}
	rd := io.LimitReader(r, int64(header.MessageLength-4*4))
	var ret common.Message
	switch header.OpCode {
	case OP_QUERY:
		ret = ParseQuery(header, rd)
	case OP_INSERT:
		ret = ParseInsert(header, rd)
	case OP_DELETE:
		ret = ParseDelete(header, rd)
	case OP_UPDATE:
		ret = ParseUpdate(header, rd)
	case OP_MSG:
		ret = ParseMsg(header, rd)
	case OP_REPLY:
		ret = ParseReply(header, rd)
	case OP_GET_MORE:
		ret = ParseGetMore(header, rd)
	case OP_KILL_CURSORS:
		ret = ParseKillCursors(header, rd)
	case OP_RESERVED:
		ParseReserved(header, rd)
	//case OP_COMMAND_DEPRECATED:
	//	self.ParseCommandDeprecated(header, rd)
	//case OP_COMMAND_REPLY_DEPRECATED:
	//	self.ParseCommandReplyDeprecated(header, rd)
	case OP_COMMAND:
		ret = ParseCommand(header, rd)
	case OP_COMMAND_REPLY:
		ret = ParseCommandReply(header, rd)
	case OP_MSG_NEW:
		ret = ParseMsgNew(header, rd)
	default:
		log.Printf("unknown OpCode: %d", header.OpCode)
		_, err = io.Copy(ioutil.Discard, rd)
		if err != nil {
			log.Printf("read failed: %v", err)
			return nil, err
		}
	}
	return ret, nil
}

func parseQuery(query primitive.M) map[string]string {
	result := map[string]string{}
	for k, v := range query {
		bs, err := json.Marshal(v)
		if err != nil {
			result[k] = fmt.Sprintf("%v", v)
		} else {
			result[k] = string(bs)
		}

	}
	return result
}
