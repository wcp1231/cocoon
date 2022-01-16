package api

import (
	"bufio"
	"cocoon/pkg/model/common"
	"io"
)

type TcpReader interface {
	io.Reader

	BufferReader() *bufio.Reader
	ReadCurrent() []byte
	Connection() *common.ConnectionInfo
}
