package dubbo

import (
	"cocoon/pkg/model/common"
)

const (
	DUBBO_HEARTBEAT_KEY = "DUBBO_HEARTBEAT"
	DUBBO_REQUEST_KEY   = "DUBBO_REQUEST"
	DUBBO_RESPONSE_KEY  = "DUBBO_RESPONSE"
)

type DubboMessage struct {
	common.GenericMessage

	raw []byte
}

func NewDubboGenericMessage() *DubboMessage {
	return &DubboMessage{
		GenericMessage: common.NewGenericMessage(common.PROTOCOL_DUBBO.Name),
	}
}

func (d *DubboMessage) SetHeartbeat() {
	d.Meta["HEARTBEAT"] = "true"
	d.Payload[DUBBO_HEARTBEAT_KEY] = true
}

func (d *DubboMessage) IsHeartbeat() bool {
	hb, exist := d.Payload[DUBBO_HEARTBEAT_KEY]
	return exist && hb.(bool)
}

func (d *DubboMessage) SetRequest(request *DubboRequest) {
	d.Payload[DUBBO_REQUEST_KEY] = request
}
func (d *DubboMessage) GetRequest() *DubboRequest {
	return d.Payload[DUBBO_REQUEST_KEY].(*DubboRequest)
}
func (d *DubboMessage) SetResponse(response *DubboResponse) {
	if response.Exception != "" {
		d.MarkException()
	}
	d.Payload[DUBBO_RESPONSE_KEY] = response
}
func (d *DubboMessage) GetResponse() *DubboResponse {
	return d.Payload[DUBBO_RESPONSE_KEY].(*DubboResponse)
}

func (d *DubboMessage) setRaw(raw []byte) {
	d.raw = raw
}
func (d *DubboMessage) GetRaw() []byte {
	return d.raw
}
