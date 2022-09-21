package redis

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

	reqFile, err := os.Open("./samples/redis_request.raw")
	if err != nil {
		t.Error(err)
		return
	}
	reqData, err := io.ReadAll(reqFile)
	if err != nil {
		t.Error(err)
		return
	}
	respFile, err := os.Open("./samples/redis_response.raw")
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
	go dissector.StartRequestDissect(bufio.NewReader(bytes.NewReader(reqData)))

	reqResult := readAllFromChan(reqC)
	c := bytes.Compare(reqData, reqResult)
	if c != 0 {
		t.Fail()
	}

	go dissector.StartResponseDissect(bufio.NewReader(bytes.NewReader(respData)))
	respResult := readAllFromChan(respC)
	c = bytes.Compare(respData, respResult)
	if c != 0 {
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
