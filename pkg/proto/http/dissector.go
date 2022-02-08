package http

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/common"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Dissector struct {
	reqReader  *bufio.Reader
	respReader *bufio.Reader
	requestC   chan *common.GenericMessage
	responseC  chan *common.GenericMessage
}

func NewRequestDissector(reqC, respC chan *common.GenericMessage) *Dissector {
	return &Dissector{
		requestC:  reqC,
		responseC: respC,
	}
}

func (d *Dissector) StartRequestDissect(reader *bufio.Reader) {
	d.reqReader = reader
	for {
		err := d.dissectRequest()
		if err != nil {
			break
		}
	}
	fmt.Println("Http request dissect finish")
	close(d.requestC)
}

func (d *Dissector) StartResponseDissect(reader *bufio.Reader) {
	d.respReader = reader
	for {
		err := d.dissectResponse()
		if err != nil {
			break
		}
	}
	fmt.Println("Http response dissect finish")
	close(d.responseC)
}

func (d *Dissector) dissectRequest() error {
	message := common.NewHTTPGenericMessage()

	request, err := http.ReadRequest(d.reqReader)
	if err != nil {
		if err == io.EOF {
			// conn close
			return err
		}
		fmt.Println("Http request dissect error", err.Error())
		return err
		// TODO response 500?
	}

	for k, vv := range request.Header {
		message.Header[k] = strings.Join(vv, ";")
	}

	message.Meta["HOST"] = request.Host
	message.Meta["METHOD"] = request.Method
	message.Meta["URL"] = request.URL.String()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Http request read body error", err.Error())
	}
	message.Body = &body
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	buf := new(bytes.Buffer)
	err = request.Write(buf)
	if err != nil {
		fmt.Println("Http request read raw error", err.Error())
	}
	bs := buf.Bytes()
	message.Raw = &bs

	message.CaptureNow()
	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	message := common.NewHTTPGenericMessage()

	response, err := http.ReadResponse(d.respReader, nil)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			// conn close
			return err
		}
		fmt.Printf("Http response dissect error %v\n", err)
		return err
		// TODO response 500?
	}

	for k, vv := range response.Header {
		message.Header[k] = strings.Join(vv, ";;")
	}

	message.Meta["STATUS"] = response.Status
	message.Meta["PROTO"] = response.Proto

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Http request read body error", err.Error())
	}
	message.Body = &body
	response.Body = io.NopCloser(bytes.NewBuffer(body))

	buf := new(bytes.Buffer)
	err = response.Write(buf)
	bs := buf.Bytes()
	if err != nil {
		fmt.Println("Http response read raw error", err.Error())
	}
	message.Raw = &bs

	message.CaptureNow()
	d.responseC <- message
	return nil
}
