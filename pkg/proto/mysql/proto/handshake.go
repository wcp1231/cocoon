package proto

import (
	"errors"
	"fmt"
)

type HandshakeV9 struct {
	ProtocolVersion byte
	ServerVersion   string
	ConnectionId    uint32
	Authentication  string
}

type HandshakeV10 struct {
	ProtocolVersion   byte
	ServerVersion     string
	ConnectionId      uint32
	AuthPluginPart1   []byte
	Filter            byte
	CapabilityFlags1  uint16
	CharacterSet      uint8
	StatusFlags       uint16
	CapabilityFlags2  uint16
	AuthPluginDataLen uint8
	Reserved          []byte
	AuthPluginPart2   string
	AuthPluginName    []byte
}

// UnPackHandshake used to unpack the Handshake packet.
func UnPackHandshake(data []byte) (interface{}, error) {
	buf := ReadBuffer(data)

	// protocol version
	protocolVersion, err := buf.ReadU8()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid handshake packet version: %v", data))
	}
	fmt.Printf("HandShake protocol version: %v\n", protocolVersion)
	if protocolVersion == 10 {
		return unPackHandshakeV10(buf)
	}
	return unPackHandshakeV9(buf)
}

func unPackHandshakeV9(buf *Buffer) (*HandshakeV9, error) {
	var err error
	handshake := &HandshakeV9{ProtocolVersion: 9}
	if handshake.ServerVersion, err = buf.ReadStringNUL(); err != nil {
		return nil, err
	}
	if handshake.ConnectionId, err = buf.ReadU32(); err != nil {
		return nil, err
	}
	if handshake.Authentication, err = buf.ReadStringNUL(); err != nil {
		return nil, err
	}
	return handshake, nil
}

func unPackHandshakeV10(buf *Buffer) (*HandshakeV10, error) {
	var err error
	handshake := &HandshakeV10{ProtocolVersion: 10}
	if handshake.ServerVersion, err = buf.ReadStringNUL(); err != nil {
		return nil, err
	}
	if handshake.ConnectionId, err = buf.ReadU32(); err != nil {
		return nil, err
	}
	if handshake.AuthPluginPart1, err = buf.ReadBytes(8); err != nil {
		return nil, err
	}
	if handshake.Filter, err = buf.ReadU8(); err != nil {
		return nil, err
	}
	if handshake.CapabilityFlags1, err = buf.ReadU16(); err != nil {
		return nil, err
	}
	if handshake.CharacterSet, err = buf.ReadU8(); err != nil {
		return nil, err
	}
	if handshake.StatusFlags, err = buf.ReadU16(); err != nil {
		return nil, err
	}
	if handshake.CapabilityFlags2, err = buf.ReadU16(); err != nil {
		return nil, err
	}
	// TODO capabilities & CLIENT_PLUGIN_AUTH or 00
	//if handshake.AuthPluginDataLen, err = buf.ReadU8(); err != nil {
	//	return nil, err
	//}
	if handshake.Reserved, err = buf.ReadBytes(10); err != nil {
		return nil, err
	}
	if handshake.AuthPluginPart2, err = buf.ReadLenEncodeString(); err != nil {
		return nil, err
	}
	// TODO capabilities & CLIENT_PLUGIN_AUTH
	//if handshake.AuthPluginName, err = buf.ReadBytesNUL(); err != nil {
	//	return nil, err
	//}
	return handshake, nil
}

func (h *HandshakeV9) String() string {
	return fmt.Sprintf(`HandshakeV9(
  ServerVerison=%s
  ConnectionId=%d
  Authentication=%s
`, h.ServerVersion, h.ConnectionId, h.Authentication)
}

func (h *HandshakeV10) String() string {
	return fmt.Sprintf(`HandshakeV10(
  ServerVerison=%s
  ConnectionId=%d
  AuthPluginPart1=%s
  CapabilityFlags1=%b
  CharacterSet=%d
  StatusFlags=%d
  CapabilityFlags2=%b
  AuthPluginDataLen=%d
  Reserved=%s
  AuthPluginPart2=%s
  AuthPluginName=%s
`, h.ServerVersion, h.ConnectionId, h.AuthPluginPart1, h.CapabilityFlags1, h.CharacterSet,
		h.StatusFlags, h.CapabilityFlags2, h.AuthPluginDataLen, h.Reserved, h.AuthPluginPart2, h.AuthPluginName)
}
