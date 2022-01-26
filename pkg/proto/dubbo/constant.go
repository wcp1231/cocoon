package dubbo

import (
	"errors"
	"regexp"
)

const (
	GroupKey               = "group"
	VersionKey             = "version"
	InterfaceKey           = "interface"
	MessageSizeKey         = "message_size"
	PathKey                = "path"
	ServiceKey             = "service"
	MethodsKey             = "methods"
	TimeoutKey             = "timeout"
	CategoryKey            = "category"
	CheckKey               = "check"
	EnabledKey             = "enabled"
	SideKey                = "side"
	OverrideProvidersKey   = "providerAddresses"
	BeanNameKey            = "bean.name"
	GenericKey             = "generic"
	ClassifierKey          = "classifier"
	TokenKey               = "token"
	LocalAddr              = "local-addr"
	RemoteAddr             = "remote-addr"
	DefaultRemotingTimeout = 3000
	ReleaseKey             = "release"
	AnyhostKey             = "anyhost"
	PortKey                = "port"
	ProtocolKey            = "protocol"
	PathSeparator          = "/"
	CommaSeparator         = ","
	SslEnabledKey          = "ssl-enabled"
	// ParamsTypeKey key used in pass through invoker factory, to define param type
	ParamsTypeKey   = "parameter-type-names"
	MetadataTypeKey = "metadata-type"

	// Body map keys
	DubboVersionKey = "dubboVersion"
	ArgsTypesKey    = "argsTypes"
	ArgsKey         = "args"
	AttachmentsKey  = "attachments"
)

const (
	// header length.
	HEADER_LENGTH = 16

	// magic header
	MAGIC      = uint16(0xdabb)
	MAGIC_HIGH = byte(0xda)
	MAGIC_LOW  = byte(0xbb)

	// message flag.
	FLAG_REQUEST = byte(0x80)
	FLAG_TWOWAY  = byte(0x40)
	FLAG_EVENT   = byte(0x20) // for heartbeat
	SERIAL_MASK  = 0x1f

	DUBBO_VERSION                          = "2.5.4"
	DUBBO_VERSION_KEY                      = "dubbo"
	DEFAULT_DUBBO_PROTOCOL_VERSION         = "2.0.2" // Dubbo RPC protocol version, for compatibility, it must not be between 2.0.10 ~ 2.6.2
	LOWEST_VERSION_FOR_RESPONSE_ATTACHMENT = 2000200
	DEFAULT_LEN                            = 8388608 // 8 * 1024 * 1024 default body max length

	Zero             = byte(0x00)
)

const (
	SHessian2 byte = 2
	SProto    byte = 21
)

const (
	Hessian2Serialization = "hessian2"
	ProtobufSerialization = "protobuf"
	MsgpackSerialization  = "msgpack"
)

// Dubbo request response related consts
var (
	DubboRequestHeaderBytesTwoWay = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST | FLAG_TWOWAY}
	DubboRequestHeaderBytes       = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST}
	DubboResponseHeaderBytes      = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, Zero, Response_OK}
	DubboRequestHeartbeatHeader   = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST | FLAG_TWOWAY | FLAG_EVENT}
	DubboResponseHeartbeatHeader  = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_EVENT, Response_OK}
)

// Error part
var (
	ErrHeaderNotEnough = errors.New("header buffer too short")
	ErrBodyNotEnough   = errors.New("body buffer too short")
	ErrJavaException   = errors.New("got java exception")
	ErrIllegalPackage  = errors.New("illegal package!")
)

const (
	DefaultKey   = "default"
	Generic      = "$invoke"
	GenericAsync = "$invokeAsync"
	Echo         = "$echo"
)

// ResponsePayload related consts
const (
	Response_OK                byte = 20
	Response_CLIENT_TIMEOUT    byte = 30
	Response_SERVER_TIMEOUT    byte = 31
	Response_BAD_REQUEST       byte = 40
	Response_BAD_RESPONSE      byte = 50
	Response_SERVICE_NOT_FOUND byte = 60
	Response_SERVICE_ERROR     byte = 70
	Response_SERVER_ERROR      byte = 80
	Response_CLIENT_ERROR      byte = 90

	// According to "java dubbo" There are two cases of response:
	// 		1. with attachments
	// 		2. no attachments
	RESPONSE_WITH_EXCEPTION                  int32 = 0
	RESPONSE_VALUE                           int32 = 1
	RESPONSE_NULL_VALUE                      int32 = 2
	RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS int32 = 3
	RESPONSE_VALUE_WITH_ATTACHMENTS          int32 = 4
	RESPONSE_NULL_VALUE_WITH_ATTACHMENTS     int32 = 5
)

// regular
const (
	JAVA_IDENT_REGEX = "(?:[_$a-zA-Z][_$a-zA-Z0-9]*)"
	CLASS_DESC = "(?:L" + JAVA_IDENT_REGEX + "(?:\\/" + JAVA_IDENT_REGEX + ")*;)"
	ARRAY_DESC = "(?:\\[+(?:(?:[VZBCDFIJS])|" + CLASS_DESC + "))"
	DESC_REGEX = "(?:(?:[VZBCDFIJS])|" + CLASS_DESC + "|" + ARRAY_DESC + ")"
)

var DescRegex, _ = regexp.Compile(DESC_REGEX)