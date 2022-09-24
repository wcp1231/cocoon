package mongo

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

	reqFile, err := os.Open("./samples/mongo_request.raw")
	if err != nil {
		t.Error(err)
		return
	}
	reqData, err := io.ReadAll(reqFile)
	if err != nil {
		t.Error(err)
		return
	}
	respFile, err := os.Open("./samples/mongo_response.raw")
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
	resultChan := make(chan []byte)
	go readAllFromChan(reqC, resultChan)
	go dissector.StartRequestDissect(bufio.NewReader(bytes.NewReader(reqData)))
	reqResult := <-resultChan

	if bytes.Compare(reqData, reqResult) != 0 {
		t.Errorf("Mongo Request dissect dismatch")
	}

	go readAllFromChan(respC, resultChan)
	go dissector.StartResponseDissect(bufio.NewReader(bytes.NewReader(respData)))
	respResult := <-resultChan
	if bytes.Compare(respData, respResult) != 0 {
		t.Errorf("Mongo Response dissect dismatch.")
	}
}

func readAllFromChan(c chan common.Message, result chan []byte) {
	buf := bytes.Buffer{}
	for {
		select {
		case data, more := <-c:
			if !more {
				result <- buf.Bytes()
				return
			}
			buf.Write(*data.GetRaw())
		}
	}

}
