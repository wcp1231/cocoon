package dubbo

import (
	"bufio"
	"cocoon/pkg/model/common"
)

const (
	Zero = byte(0x00)

	HEADER_LENGTH = 16

	MAGIC      = uint16(0xdabb)
	MAGIC_HIGH = byte(0xda)
	MAGIC_LOW  = byte(0xbb)

	FLAG_REQUEST = byte(0x80)
	FLAG_TWOWAY  = byte(0x40)
	FLAG_EVENT   = byte(0x20)
	SERIAL_MASK  = 0x1f

	Response_OK                byte = 20
	Response_CLIENT_TIMEOUT    byte = 30
	Response_SERVER_TIMEOUT    byte = 31
	Response_BAD_REQUEST       byte = 40
	Response_BAD_RESPONSE      byte = 50
	Response_SERVICE_NOT_FOUND byte = 60
	Response_SERVICE_ERROR     byte = 70
	Response_SERVER_ERROR      byte = 80
	Response_CLIENT_ERROR      byte = 90
)

type Classifier struct{}

func (c *Classifier) Match(r *bufio.Reader) bool {
	buf, err := r.Peek(16)
	if err != nil {
		return false
	}

	if buf[0] != MAGIC_HIGH && buf[1] != MAGIC_LOW {
		return false
	}

	serialId := buf[2] & SERIAL_MASK
	if serialId == Zero {
		return false
	}

	status := buf[3]
	if status != Response_OK &&
		status != Response_CLIENT_TIMEOUT &&
		status != Response_SERVER_TIMEOUT &&
		status != Response_BAD_REQUEST &&
		status != Response_BAD_RESPONSE &&
		status != Response_SERVICE_NOT_FOUND &&
		status != Response_SERVICE_ERROR &&
		status != Response_SERVER_ERROR &&
		status != Response_CLIENT_ERROR {
		return false
	}

	return true
}

func (c *Classifier) Protocol() *common.Protocol {
	return common.PROTOCOL_DUBBO
}
