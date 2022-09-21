package http

import (
	"bufio"
	"bytes"
	"cocoon/pkg/model/common"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Dissector struct {
	reqReader  *bufio.Reader
	respReader *bufio.Reader
	requestC   chan common.Message
	responseC  chan common.Message
}

func NewRequestDissector(reqC, respC chan common.Message) *Dissector {
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
	message := NewHTTPGenericMessage()

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

	message.SetHttpHeader(request.Header)
	message.SetHost(request.Host)
	message.SetMethod(request.Method)
	message.SetUrl(request.URL.String())

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Http request read body error", err.Error())
	}
	message.SetBody(body)
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
	message := NewHTTPGenericMessage()

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

	message.SetHttpHeader(response.Header)
	message.SetStatusCode(response.StatusCode)
	message.SetProto(response.Proto)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Http request read body error", err.Error())
	}
	message.SetBody(body)
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
