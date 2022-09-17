package proto

import (
	"errors"
	"fmt"
)

type HandshakeResponse interface {
	GetClientFlag() uint32
}

type HandshakeResponse320 struct {
	ClientFlag    uint16
	MaxPacketSize uint32 // uint24
	Username      string
	AuthResponse  string
	Database      string
}

type HandshakeResponse41 struct {
	ClientFlag         uint32
	MaxPacketSize      uint32
	CharacterSet       uint8
	Filter             []byte
	Username           string
	AuthResponseLength uint8
	AuthResponse       string
	Database           string
	ClientPluginName   string
	KeyValLength       uint64
	AttributeKeys      []string
	AttributeValues    []string
}

// UnPackHandshakeResponse used to unpack the Handshake Response packet.
func UnPackHandshakeResponse(data []byte) (HandshakeResponse, error) {
	buf := ReadBuffer(data)

	// client flag
	clientFlag, err := buf.ReadU16()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid handshake response packet version: %v", data))
	}
	buf.Reset(data)
	if uint32(clientFlag)&CLIENT_PROTOCOL_41 != 0 {
		return unPackHandshakeResponse41(buf)
	}
	return unPackHandshakeResponse320(buf)
}

func unPackHandshakeResponse320(buf *Buffer) (*HandshakeResponse320, error) {
	var err error
	handshake := &HandshakeResponse320{}
	if handshake.ClientFlag, err = buf.ReadU16(); err != nil {
		return nil, err
	}
	if handshake.MaxPacketSize, err = buf.ReadU24(); err != nil {
		return nil, err
	}
	if handshake.Username, err = buf.ReadStringNUL(); err != nil {
		return nil, err
	}
	if uint32(handshake.ClientFlag)&CLIENT_CONNECT_WITH_DB != 0 {
		if handshake.AuthResponse, err = buf.ReadStringNUL(); err != nil {
			return nil, err
		}
		if handshake.Database, err = buf.ReadStringNUL(); err != nil {
			return nil, err
		}
	} else {
		if handshake.AuthResponse, err = buf.ReadStringEOF(); err != nil {
			return nil, err
		}
	}
	return handshake, nil
}

func unPackHandshakeResponse41(buf *Buffer) (*HandshakeResponse41, error) {
	var err error
	handshake := &HandshakeResponse41{}
	if handshake.ClientFlag, err = buf.ReadU32(); err != nil {
		return nil, err
	}
	if handshake.MaxPacketSize, err = buf.ReadU32(); err != nil {
		return nil, err
	}
	if handshake.CharacterSet, err = buf.ReadU8(); err != nil {
		return nil, err
	}
	if handshake.Filter, err = buf.ReadBytes(23); err != nil {
		return nil, err
	}
	if handshake.Username, err = buf.ReadStringNUL(); err != nil {
		return nil, err
	}
	if handshake.ClientFlag&CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA != 0 {
		if handshake.AuthResponse, err = buf.ReadLenEncodeString(); err != nil {
			return nil, err
		}
	} else {
		if handshake.AuthResponseLength, err = buf.ReadU8(); err != nil {
			return nil, err
		}
		//if handshake.AuthResponse, err = buf.ReadLenEncodeString(); err != nil {
		//	return nil, err
		//}
		if handshake.AuthResponse, err = buf.ReadString(int(handshake.AuthResponseLength)); err != nil {
			return nil, err
		}
	}
	if handshake.ClientFlag&CLIENT_CONNECT_WITH_DB != 0 {
		if handshake.Database, err = buf.ReadStringNUL(); err != nil {
			return nil, err
		}
	}
	if handshake.ClientFlag&CLIENT_PLUGIN_AUTH != 0 {
		if handshake.ClientPluginName, err = buf.ReadStringNUL(); err != nil {
			return nil, err
		}
	}
	if handshake.ClientFlag&CLIENT_CONNECT_ATTRS != 0 {
		if handshake.KeyValLength, err = buf.ReadLenEncode(); err != nil {
			return nil, err
		}
		var i uint64
		var key, value string
		for i = 0; i < handshake.KeyValLength; i++ {
			if key, err = buf.ReadLenEncodeString(); err != nil {
				return nil, err
			}
			if value, err = buf.ReadLenEncodeString(); err != nil {
				return nil, err
			}
			handshake.AttributeKeys = append(handshake.AttributeKeys, key)
			handshake.AttributeValues = append(handshake.AttributeValues, value)
		}
	}
	return handshake, nil
}

func (h *HandshakeResponse320) GetClientFlag() uint32 {
	return uint32(h.ClientFlag)
}

func (h *HandshakeResponse41) GetClientFlag() uint32 {
	return h.ClientFlag
}

func (h *HandshakeResponse320) String() string {
	return fmt.Sprintf(`HandshakeResponse41(
  ClientFlag=%b
  MaxPacketSize=%d
  Username=%s
  AuthResponse=%s
  Database=%s
`, h.ClientFlag, h.MaxPacketSize, h.Username, h.AuthResponse, h.Database)
}

func (h *HandshakeResponse41) String() string {
	return fmt.Sprintf(`HandshakeResponse41(
  ClientFlag=%b
  MaxPacketSize=%d
  CharacterSet=%d
  Filter=%s
  Username=%s
  AuthResponseLength=%d
  AuthResponse=%s
  Database=%s
  ClientPluginName=%s
  KeyValLength=%d
`, h.ClientFlag, h.MaxPacketSize, h.CharacterSet, h.Filter, h.Username,
		h.AuthResponseLength, h.AuthResponse, h.Database, h.ClientPluginName, h.KeyValLength)
}
