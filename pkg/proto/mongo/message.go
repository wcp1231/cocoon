package mongo

import (
	"cocoon/pkg/model/common"
	"strconv"
	"strings"
)

const (
	MONGO_OP_KEY               = "MONGO_OP_TYPE"
	MONGO_REQ_ID_KEY           = "MONGO_REQ_ID"
	MONGO_RESPONSE_TO_KEY      = "MONGO_RESPONSE_TO"
	MONGO_FLAG_KEY             = "MONGO_FLAG"
	MONGO_FLAGS_KEY            = "MONGO_FLAGS"
	MONGO_DATABASE_KEY         = "MONGO_DATABASE"
	MONGO_COLLECTION_KEY       = "MONGO_COLLECTION"
	MONGO_BODY_KEY             = "MONGO_BODY"
	MONGO_SELECTOR_KEY         = "MONGO_SELECTOR"
	MONGO_MSG_KEY              = "MONGO_MSG"
	MONGO_NUMBER_TO_RETURN_KEY = "MONGO_NUMBER_TO_RETURN"
	MONGO_CURSOR_ID_KEY        = "MONGO_CURSOR_ID"
	MONGO_CURSOR_IDS_KEY       = "MONGO_CURSOR_IDS"
	MONGO_QUERY_KEY            = "MONGO_QUERY"
	MONGO_COMMAND_KEY          = "MONGO_COMMAND"
	MONGO_COMMAND_ARGS_KEY     = "MONGO_COMMAND_ARGS"
	MONGO_COMMAND_REPLY_KEY    = "MONGO_COMMAND_REPLY"
	MONGO_METADATA_KEY         = "MONGO_METADATA"
	MONGO_NUMBER_TO_SKIP_KEY   = "MONGO_NUMBER_TO_SKIP"
	MONGO_STARTING_FROM_KEY    = "MONGO_STARTING_FROM"
)

type MongoMessage struct {
	common.GenericMessage
}

func NewMongoGenericMessage() *MongoMessage {
	return &MongoMessage{
		common.NewGenericMessage(common.PROTOCOL_MONGO.Name),
	}
}

func (h *MongoMessage) SetOpType(opType string) {
	h.Payload[MONGO_OP_KEY] = opType
	h.Meta["OP_TYPE"] = opType
}
func (h *MongoMessage) GetOpType() string {
	return h.Payload[MONGO_OP_KEY].(string)
}
func (h *MongoMessage) SetRequestId(reqId int32) {
	h.Payload[MONGO_REQ_ID_KEY] = reqId
	h.Meta["REQUEST_ID"] = strconv.Itoa(int(reqId))
}
func (h *MongoMessage) GetRequestId() int32 {
	return h.Payload[MONGO_REQ_ID_KEY].(int32)
}
func (h *MongoMessage) SetResponseTo(responseTo int32) {
	h.Payload[MONGO_RESPONSE_TO_KEY] = responseTo
	h.Meta["RESPONSE_TO"] = strconv.Itoa(int(responseTo))
}
func (h *MongoMessage) GetResponseTo() int32 {
	return h.Payload[MONGO_RESPONSE_TO_KEY].(int32)
}
func (h *MongoMessage) SetFlag(flag int32) {
	h.Payload[MONGO_FLAG_KEY] = flag
	h.Meta["FLAG"] = strconv.Itoa(int(flag))
}
func (h *MongoMessage) GetFlag() int32 {
	return h.Payload[MONGO_FLAG_KEY].(int32)
}
func (h *MongoMessage) SetFlags(flags string) {
	h.Payload[MONGO_FLAGS_KEY] = flags
	h.Meta["FLAGS"] = flags
}
func (h *MongoMessage) GetFlags() string {
	return h.Payload[MONGO_FLAGS_KEY].(string)
}
func (h *MongoMessage) SetDatabase(database string) {
	h.Payload[MONGO_DATABASE_KEY] = database
	h.Meta["DATABASE"] = database
}
func (h *MongoMessage) GetDatabase() string {
	return h.Payload[MONGO_DATABASE_KEY].(string)
}
func (h *MongoMessage) SetCollection(collection string) {
	h.Payload[MONGO_COLLECTION_KEY] = collection
	h.Meta["COLLECTION"] = collection
}
func (h *MongoMessage) GetCollection() string {
	return h.Payload[MONGO_COLLECTION_KEY].(string)
}
func (h *MongoMessage) SetSelector(selector string) {
	h.Payload[MONGO_SELECTOR_KEY] = selector
	h.Meta["SELECTOR"] = selector
}
func (h *MongoMessage) GetSelector() string {
	return h.Payload[MONGO_SELECTOR_KEY].(string)
}
func (h *MongoMessage) SetMongoBody(body string) {
	h.Payload[MONGO_BODY_KEY] = body
	bodyBytes := []byte(body)
	h.Body = &bodyBytes
}
func (h *MongoMessage) GetMongoBody() string {
	return h.Payload[MONGO_BODY_KEY].(string)
}
func (h *MongoMessage) SetMsg(msg string) {
	h.Payload[MONGO_MSG_KEY] = msg
	h.Meta["MSG"] = msg
}
func (h *MongoMessage) GetMsg() string {
	return h.Payload[MONGO_MSG_KEY].(string)
}
func (h *MongoMessage) SetNumberToReturn(numberToReturn int32) {
	h.Payload[MONGO_NUMBER_TO_RETURN_KEY] = numberToReturn
	h.Meta["RETURN"] = strconv.Itoa(int(numberToReturn))
}
func (h *MongoMessage) GetNumberToReturn() int32 {
	return h.Payload[MONGO_NUMBER_TO_RETURN_KEY].(int32)
}
func (h *MongoMessage) SetCursorId(cursorId int64) {
	h.Payload[MONGO_CURSOR_ID_KEY] = cursorId
	h.Meta["CURSOR_ID"] = strconv.FormatInt(cursorId, 10)
}
func (h *MongoMessage) GetCursorId() int64 {
	return h.Payload[MONGO_CURSOR_ID_KEY].(int64)
}
func (h *MongoMessage) SetCursorIds(cursorIds []int64) {
	h.Payload[MONGO_CURSOR_IDS_KEY] = cursorIds
	var ids []string
	for _, id := range cursorIds {
		ids = append(ids, strconv.FormatInt(id, 10))
	}
	h.Meta["CURSOR_IDS"] = strings.Join(ids, ",")
}
func (h *MongoMessage) GetCursorIds() []int64 {
	return h.Payload[MONGO_CURSOR_IDS_KEY].([]int64)
}
func (h *MongoMessage) SetQuery(query string) {
	h.Payload[MONGO_QUERY_KEY] = query
	h.Meta["QUERY"] = query
}
func (h *MongoMessage) GetQuery() string {
	return h.Payload[MONGO_QUERY_KEY].(string)
}
func (h *MongoMessage) SetCommand(command string) {
	h.Payload[MONGO_COMMAND_KEY] = command
	h.Meta["COMMAND"] = command
}
func (h *MongoMessage) GetCommand() string {
	return h.Payload[MONGO_COMMAND_KEY].(string)
}
func (h *MongoMessage) SetCommandArgs(args string) {
	h.Payload[MONGO_COMMAND_ARGS_KEY] = args
	h.Meta["ARGS"] = args
}
func (h *MongoMessage) GetCommandArgs() string {
	return h.Payload[MONGO_COMMAND_ARGS_KEY].(string)
}
func (h *MongoMessage) SetCommandReply(args string) {
	h.Payload[MONGO_COMMAND_REPLY_KEY] = args
	h.Meta["REPLY"] = args
}
func (h *MongoMessage) GetCommandReply() string {
	return h.Payload[MONGO_COMMAND_REPLY_KEY].(string)
}
func (h *MongoMessage) SetMetadata(metadata string) {
	h.Payload[MONGO_METADATA_KEY] = metadata
	h.Meta["METADATA"] = metadata
}
func (h *MongoMessage) GetMetadata() string {
	return h.Payload[MONGO_METADATA_KEY].(string)
}
func (h *MongoMessage) SetNumberToSkip(numberToSkip int32) {
	h.Payload[MONGO_NUMBER_TO_SKIP_KEY] = numberToSkip
	h.Meta["SKIP"] = strconv.Itoa(int(numberToSkip))
}
func (h *MongoMessage) GetNumberToSkip() int32 {
	return h.Payload[MONGO_NUMBER_TO_SKIP_KEY].(int32)
}
func (h *MongoMessage) SetStartingFrom(startingFrom int32) {
	h.Payload[MONGO_STARTING_FROM_KEY] = startingFrom
	h.Meta["STARTING"] = strconv.Itoa(int(startingFrom))
}
func (h *MongoMessage) GetStartingFrom() int32 {
	return h.Payload[MONGO_STARTING_FROM_KEY].(int32)
}
