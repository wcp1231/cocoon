package http

import (
	"bufio"
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

	httpRequest, err := d.parseRequest(d.reqReader)
	if err != nil {
		if err == io.EOF {
			// conn close
			return err
		}
		fmt.Println("Http request dissect error", err.Error())
		return err
	}

	message.SetRequest(httpRequest)
	message.CaptureNow()
	d.requestC <- message
	return nil
}

func (d *Dissector) dissectResponse() error {
	message := NewHTTPGenericMessage()

	httpResponse, err := d.parseResponse(d.respReader)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			// conn close
			return err
		}
		fmt.Printf("Http response dissect error %v\n", err)
		return err
	}

	message.SetResponse(httpResponse)
	message.CaptureNow()
	d.responseC <- message
	return nil
}

func (d *Dissector) parseRequest(r *bufio.Reader) (*HttpReuqest, error) {
	request, err := http.ReadRequest(r)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	result := &HttpReuqest{
		Header: request.Header,
		Host:   request.Host,
		Method: request.Method,
		URL:    request.URL.String(),
		Body:   body,
	}
	return result, nil
}

func (d *Dissector) parseResponse(r *bufio.Reader) (*HttpResponse, error) {
	response, err := http.ReadResponse(r, nil)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := &HttpResponse{
		StatusCode: response.StatusCode,
		Proto:      response.Proto,
		ProtoMajor: response.ProtoMajor,
		ProtoMinor: response.ProtoMinor,
		Header:     response.Header,
		Body:       body,
	}
	return result, nil
}
