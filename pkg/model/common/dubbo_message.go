package common

import (
	"encoding/json"
	"fmt"
)

const (
	DUBBO_HEARTBEAT_KEY       = "DUBBO_HEARTBEAT"
	DUBBO_VERSION_KEY         = "DUBBO_VERSION"
	DUBBO_SERVICE_VERSION_KEY = "DUBBO_SERVICE_VERSION"
	DUBBO_METHOD_KEY          = "DUBBO_METHOD"
	DUBBO_TARGET_KEY          = "DUBBO_TARGET"
	DUBBO_ARGS_KEY            = "DUBBO_ARGS"
	DUBBO_ATTACHMENTS_KEY     = "DUBBO_ATTACHMENTS"
	DUBBO_EXCEPTION_KEY       = "DUBBO_EXCEPTION"
	DUBBO_RESPONSE_KEY        = "DUBBO_RESPONSE"
)

type DubboMessage struct {
	GenericMessage
}

func NewDubboGenericMessage() *DubboMessage {
	return &DubboMessage{
		NewGenericMessage(PROTOCOL_DUBBO.Name),
	}
}

func (d *DubboMessage) SetHeartbeat() {
	d.Meta["HEARTBEAT"] = "true"
	d.Payload[DUBBO_HEARTBEAT_KEY] = true
}

func (d *DubboMessage) IsHeartbeat() bool {
	hb, exist := d.Payload[DUBBO_HEARTBEAT_KEY]
	return !exist || hb.(bool)
}

func (d *DubboMessage) SetDubboVersion(dubboVersion string) {
	d.Header["dubboVersion"] = dubboVersion
	d.Payload[DUBBO_VERSION_KEY] = dubboVersion
}

func (d *DubboMessage) GetDubboVersion() string {
	return d.Payload[DUBBO_VERSION_KEY].(string)
}

func (d *DubboMessage) SetServiceVersion(serviceVersion string) {
	d.Header["serviceVersion"] = serviceVersion
	d.Payload[DUBBO_SERVICE_VERSION_KEY] = serviceVersion
}

func (d *DubboMessage) GetServiceVersion() string {
	return d.Payload[DUBBO_SERVICE_VERSION_KEY].(string)
}

func (d *DubboMessage) SetMethod(method string) {
	d.Header["method"] = method
	d.Payload[DUBBO_METHOD_KEY] = method
}

func (d *DubboMessage) GetMethod() string {
	return d.Payload[DUBBO_METHOD_KEY].(string)
}

func (d *DubboMessage) SetTarget(target string) {
	d.Header["target"] = target
	d.Payload[DUBBO_TARGET_KEY] = target
}

func (d *DubboMessage) GetTarget() string {
	return d.Payload[DUBBO_TARGET_KEY].(string)
}

func (d *DubboMessage) SetArgs(args map[string]interface{}) {
	d.Header["args"] = formatDubboAttachments(args)
	d.Payload[DUBBO_ARGS_KEY] = args
}

func (d *DubboMessage) GetArgs() map[string]interface{} {
	return d.Payload[DUBBO_ARGS_KEY].(map[string]interface{})
}

func (d *DubboMessage) SetAttachments(attachments map[string]interface{}) {
	d.Header["attachments"] = formatDubboAttachments(attachments)
	d.Payload[DUBBO_ATTACHMENTS_KEY] = attachments
}

func (d *DubboMessage) GetAttachments() map[string]interface{} {
	return d.Payload[DUBBO_ATTACHMENTS_KEY].(map[string]interface{})
}

func (d *DubboMessage) SetException(exception string) {
	d.Header["exception"] = exception
	d.Payload[DUBBO_EXCEPTION_KEY] = exception
}

func (d *DubboMessage) GetException() string {
	return d.Payload[DUBBO_EXCEPTION_KEY].(string)
}

func (d *DubboMessage) SetResponse(response interface{}) {
	d.Header["respObj"] = fmt.Sprintf("%v", response)
	d.Payload[DUBBO_RESPONSE_KEY] = response
}

func (d *DubboMessage) GetResponse() interface{} {
	return d.Payload[DUBBO_RESPONSE_KEY]
}

func formatDubboAttachments(attachments map[string]interface{}) string {
	result := map[string]string{}
	for k, v := range attachments {
		result[k] = fmt.Sprintf("%v", v)
	}

	str, _ := json.Marshal(result)
	return string(str)
}
