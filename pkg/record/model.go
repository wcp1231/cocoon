package record

import (
	"cocoon/pkg/model/common"
	"encoding/json"
	"fmt"
	"time"
)

type recordTime time.Time

type record struct {
	Id          int32             `json:"id"`
	CaptureTime recordTime        `json:"captureTime"`
	IsRequest   bool              `json:"isRequest"`
	Meta        map[string]string `json:"meta"`
	Header      map[string]string `json:"header"`
	Body        string            `json:"body"`
}

func fromGenericMessage(message *common.GenericMessage, isRequest bool) *record {
	body := ""
	if message.Body != nil {
		body = string(*message.Body) // FIXME
	}
	return &record{
		Id:          message.Id,
		CaptureTime: recordTime(message.CaptureTime),
		IsRequest:   isRequest,
		Meta:        message.Meta,
		Header:      message.Header,
		Body:        body,
	}
}

func (r recordTime) MarshalJSON() ([]byte, error) {
	timestamp := fmt.Sprintf("%d", time.Time(r).UnixNano())
	return []byte(timestamp), nil
}

func (r *recordTime) UnmarshalJSON(bytes []byte) error {
	var raw int64
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	*(*time.Time)(r) = time.Unix(0, raw)
	return nil
}
