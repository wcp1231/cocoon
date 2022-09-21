package mysql

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/common"
	"io"
	"os"
	"testing"
)

func TestDissector(t *testing.T) {
	reqC := make(chan common.Message, 20)
	respC := make(chan common.Message, 20)

	reqFile, err := os.Open("./samples/mysql_request.raw")
	if err != nil {
		t.Error(err)
		return
	}
	reqData, err := io.ReadAll(reqFile)
	if err != nil {
		t.Error(err)
		return
	}
	respFile, err := os.Open("./samples/mysql_response.raw")
	if err != nil {
		t.Error(err)
		return
	}
	respData, err := io.ReadAll(respFile)
	if err != nil {
		t.Error(err)
		return
	}

	dissector := NewRequestDissector(reqC, respC)
	dissector.Init(bufio.NewReader(bytes.NewReader(reqData)), bufio.NewReader(bytes.NewReader(respData)))

	reqBuf := bytes.Buffer{}
	respBuf := bytes.Buffer{}

	bs, err := dissector.ReadServerHandshake()
	if err != nil {
		t.Error(err)
		return
	}
	respBuf.Write(bs)
	bs, err = dissector.ReadClientHandshakeResponse()
	if err != nil {
		t.Error(err)
		return
	}
	reqBuf.Write(bs)
	bs, err = dissector.ReadPacketFromServer()
	if err != nil {
		t.Error(err)
	}
	respBuf.Write(bs)

	go dissector.StartRequestDissect(nil)

	reqResult := readAllFromChan(reqC)
	reqBuf.Write(reqResult)
	c := bytes.Compare(reqData, reqBuf.Bytes())
	if c != 0 {
		t.Fail()
		return
	}

	go dissector.StartResponseDissect(nil)
	respResult := readAllFromChan(respC)
	respBuf.Write(respResult)
	// TODO 当前 float 类型的二进制格式可能会有误差
	if len(respData) != len(respBuf.Bytes()) {
		t.Fail()
	}
}

func readAllFromChan(c chan common.Message) []byte {
	buf := bytes.Buffer{}
	for {
		select {
		case data, more := <-c:
			if !more {
				return buf.Bytes()
			}
			buf.Write(*data.GetRaw())
		}
	}

}
