package dubbo

import (
	"bufio"
	"cocoon/pkg/model/common"
	"encoding/binary"
	"fmt"
	hessian "github.com/apache/dubbo-go-hessian2"
	"io"
)

func ReadPacket(reader *bufio.Reader) (*common.GenericMessage, error) {
	header, err := ReadHeader(reader)
	if err != nil {
		return nil, err
	}

	body := make([]byte, header.BodyLen)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return nil, err
	}

	if header.isRequest() {
		request, err := readRequestBody(header, body)
		fmt.Printf("Dubbo body decode. err=%v, request=%+v\n", err, request)
	} else {
		response, err := readResponseBody(header, body)
		fmt.Printf("Dubbo body decode. err=%v, response=%+v\n", err, response)
	}

	headerBytes := EncodeHeader(header)
	raw := make([]byte, len(headerBytes)+len(body))
	copy(raw, headerBytes)
	copy(raw[len(headerBytes):], body)

	message := common.NewDubboGenericMessage()
	if header.isHeartbeat() {
		message.Meta["HEARTBEAT"] = "true"
	}
	message.Body = &body // FIXME
	message.Raw = &raw
	return message, nil
}

func ReadHeader(reader *bufio.Reader) (*dubboHeader, error) {
	header := &dubboHeader{}
	var err error
	buf, err := reader.Peek(HEADER_LENGTH)
	if err != nil { // this is impossible
		return nil, err
	}
	_, err = reader.Discard(HEADER_LENGTH)
	if err != nil { // this is impossible
		return nil, err
	}

	//// read header
	if buf[0] != MAGIC_HIGH && buf[1] != MAGIC_LOW {
		return nil, ErrIllegalPackage
	}

	// Header{serialization id(5 bit), event, two way, req/response}
	if header.SerialID = buf[2] & SERIAL_MASK; header.SerialID == Zero {
		return nil, fmt.Errorf("serialization ID:%v", header.SerialID)
	}

	flag := buf[2] & FLAG_EVENT
	if flag != Zero {
		header.Type |= PackageHeartbeat
	}
	flag = buf[2] & FLAG_REQUEST
	if flag != Zero {
		header.Type |= PackageRequest
		flag = buf[2] & FLAG_TWOWAY
		if flag != Zero {
			header.Type |= PackageRequest_TwoWay
		}
	} else {
		header.Type |= PackageResponse
		header.ResponseStatus = buf[3]
		if header.ResponseStatus != Response_OK {
			header.Type |= PackageResponse_Exception
		}
	}

	// Header{req id}
	header.ID = int64(binary.BigEndian.Uint64(buf[4:]))

	// Header{body len}
	header.BodyLen = int(binary.BigEndian.Uint32(buf[12:]))
	if header.BodyLen < 0 {
		return nil, ErrIllegalPackage
	}

	return header, err
}

func EncodeHeader(header *dubboHeader) []byte {
	bs := make([]byte, 0)
	switch {
	case header.isHeartbeat():
		if header.ResponseStatus == Zero {
			bs = append(bs, DubboRequestHeartbeatHeader[:]...)
		} else {
			bs = append(bs, DubboResponseHeartbeatHeader[:]...)
		}
	case header.isResponse():
		bs = append(bs, DubboResponseHeaderBytes[:]...)
		if header.ResponseStatus != 0 {
			bs[3] = header.ResponseStatus
		}
	case header.Type&PackageRequest_TwoWay != 0x00:
		bs = append(bs, DubboRequestHeaderBytesTwoWay[:]...)
	}
	bs[2] |= header.SerialID & SERIAL_MASK
	binary.BigEndian.PutUint64(bs[4:], uint64(header.ID))
	binary.BigEndian.PutUint32(bs[12:], uint32(header.BodyLen))
	return bs
}

func readRequestBody(header *dubboHeader, bb []byte) (*dubboRequest, error) {
	request := &dubboRequest{
		header: header,
	}
	decoder := hessian.NewDecoder(bb)
	dubboVersion, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	if dubboVersion != nil {
		request.dubboVersion = dubboVersion.(string)
	}
	target, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	request.target = target.(string)
	serviceVersion, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	request.serviceVersion = serviceVersion.(string)
	method, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	request.method = method.(string)
	args, err := readRequestArgs(decoder)
	if err != nil {
		return nil, err
	}
	request.args = args
	attachements, err := readAttachments(decoder)
	if err != nil {
		return nil, err
	}
	request.attachments = attachements

	return request, nil
}

func readResponseBody(header *dubboHeader, bb []byte) (*dubboResponse, error) {
	response := &dubboResponse{
		header: header,
	}

	decoder := hessian.NewDecoder(bb)
	if header.hasException() {
		exception, err := decoder.Decode()
		if err != nil {
			return nil, err
		}
		response.exception = exception.(string)
		return response, nil
	}

	respType, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	if respType == RESPONSE_WITH_EXCEPTION || respType == RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS {
		exception, err := decoder.Decode()
		if err != nil {
			return nil, err
		}
		response.exception = exception.(string)
		if respType == RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS {
			attachements, err := readAttachments(decoder)
			if err != nil {
				return nil, err
			}
			response.attachments = attachements
		}
		return response, nil
	}

	if respType == RESPONSE_VALUE || respType == RESPONSE_VALUE_WITH_ATTACHMENTS {
		resp, err := decoder.Decode()
		if err != nil {
			return nil, err
		}
		response.respObj = resp
		if respType == RESPONSE_VALUE_WITH_ATTACHMENTS {
			attachements, err := readAttachments(decoder)
			if err != nil {
				return nil, err
			}
			response.attachments = attachements
		}
		return response, nil
	}

	response.respObj = nil
	if respType == RESPONSE_NULL_VALUE_WITH_ATTACHMENTS {
		attachements, err := readAttachments(decoder)
		if err != nil {
			return nil, err
		}
		response.attachments = attachements
	}

	return response, nil
}

func readRequestArgs(decoder *hessian.Decoder) (map[string]interface{}, error) {
	argsTypes, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	args := map[string]interface{}{}
	ats := DescRegex.FindAllString(argsTypes.(string), -1)
	for _, t := range ats {
		arg, err := decoder.Decode()
		if err != nil {
			return nil, err
		}
		args[t] = arg
	}
	return args, nil
}

func readAttachments(decoder *hessian.Decoder) (map[string]interface{}, error) {
	attachments, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	attachmentsMap, ok := attachments.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("get wrong attachements: %+v", attachments)
	}
	result := map[string]interface{}{}
	for k, v := range attachmentsMap {
		if kv, ok := k.(string); ok {
			if v == nil {
				result[kv] = ""
				continue
			}
			result[kv] = v
		}
	}
	return result, nil
}
