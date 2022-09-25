package record

import (
	"cocoon/pkg/model/common"
	"encoding/json"
	"fmt"
	"time"
)

type recordTime time.Time

type record struct {
	Id          int32                  `json:"id"`
	CaptureTime recordTime             `json:"captureTime"`
	IsRequest   bool                   `json:"isRequest"`
	Meta        map[string]string      `json:"meta"`
	Header      map[string]string      `json:"header"`
	Payload     map[string]interface{} `json:"payload"`
	Body        string                 `json:"body"`
}

func fromGenericMessage(message common.Message, isRequest bool) *record {
	body := ""
	mbytes := message.GetBody()
	if mbytes != nil {
		body = string(mbytes) // FIXME
	}
	return &record{
		Id:          message.ID(),
		CaptureTime: recordTime(message.GetCaptureTime()),
		IsRequest:   isRequest,
		Meta:        message.GetMeta(),
		Header:      message.GetHeader(),
		Payload:     message.GetPayload(),
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
