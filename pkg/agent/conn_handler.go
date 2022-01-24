package agent

import (
	"bufio"
	"cocoon/pkg/model/common"
	"cocoon/pkg/model/rpc"
	"cocoon/pkg/proto"
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
	"net"
)

type ConnHandler struct {
	server *Server
	ctx    context.Context
	Close  context.CancelFunc

	inboundConn    *net.TCPConn
	outboundConn   *net.TCPConn
	inboundReader  *bufio.Reader
	outboundReader *bufio.Reader
	inboundAddr    string
	outboundAddr   string

	pc        *proto.ProtoClassifier
	proto     *common.Protocol
	requestC  chan *common.GenericMessage
	responseC chan *common.GenericMessage
}

func NewConnHandler(server *Server, inboundConn, outboundConn *net.TCPConn) *ConnHandler {
	innerCtx, close := context.WithCancel(server.ctx)
	pc := proto.NewProtoClassifier()
	requestC := make(chan *common.GenericMessage)
	responseC := make(chan *common.GenericMessage)
	return &ConnHandler{
		server:         server,
		inboundConn:    inboundConn,
		outboundConn:   outboundConn,
		inboundReader:  bufio.NewReader(inboundConn),
		outboundReader: bufio.NewReader(outboundConn),
		inboundAddr:    inboundConn.RemoteAddr().String(),
		outboundAddr:   outboundConn.RemoteAddr().String(),
		ctx:            innerCtx,
		Close:          close,

		pc:        pc,
		requestC:  requestC,
		responseC: responseC,
	}
}

func (c *ConnHandler) Start() {
	defer func() {
		if err := c.inboundConn.Close(); err != nil {
			c.server.logger.WithOptions(zap.AddCaller()).Error("proxy inbound conn Close error")
		}
		if err := c.outboundConn.Close(); err != nil {
			c.server.logger.WithOptions(zap.AddCaller()).Error("proxy outbound conn Close error")
		}
		c.handleClose()
	}()

	err := c.classifyProto()
	if err != nil {
		// TODO
	}

	c.startProxy()
	select {
	case <-c.ctx.Done():
		return
	}
}

func (c *ConnHandler) startProxy() {
	if c.proto.Pass {
		c.server.logger.Info("Conn pass",
			zap.String("src", c.inboundAddr),
			zap.String("dst", c.outboundAddr),
			zap.String("proto", c.proto.Name))
		c.normalProxy()
		return
	}
	c.server.logger.Info("Conn proxy",
		zap.String("src", c.inboundAddr),
		zap.String("dst", c.outboundAddr),
		zap.String("proto", c.proto.Name))
	c.proxy()
}

func (c *ConnHandler) classifyProto() error {
	c.proto = c.pc.Classify(c.inboundReader)
	return nil
}

func (c *ConnHandler) proxy() {
	requestDissector := proto.NewRequestDissector(c.proto, c.requestC, c.responseC)
	go c.handleRequest()
	//go c.handleResponse()
	go requestDissector.StartRequestDissect(c.inboundReader)
	go requestDissector.StartResponseDissect(c.outboundReader)
}

func (c *ConnHandler) handleRequest() {
	for {
		select {
		case request, more := <-c.requestC:
			if !more {
				fmt.Printf("no more request. %v\n", request)
				return
			}
			c.server.logger.Debug("Conn request",
				zap.String("src", c.inboundAddr),
				zap.String("dst", c.outboundAddr),
				zap.String("req", request.String()))

			err := c.tryToMock(request)
			if err != nil {
				c.server.logger.Debug("Call mock server failed", zap.Error(err))
				continue
			}
		}
	}
}

//func (c *ConnHandler) handleResponse() {
//	for {
//		select {
//		case response, more := <-c.responseC:
//			if !more {
//				fmt.Printf("no more response. %v\n", response)
//				return
//			}
//			c.server.logger.Debug("Conn response",
//				zap.String("src", c.inboundAddr),
//				zap.String("dst", c.outboundAddr),
//				zap.String("resp", response.String()))
//
//			err := c.tryToRecord(nil, response)
//			if err != nil {
//				c.server.logger.Debug("Call record server failed", zap.Error(err))
//				continue
//			}
//		}
//	}
//}

// 请求 server 进行 mock
func (c *ConnHandler) tryToMock(request *common.GenericMessage) error {
	req := &rpc.OutboundReq{
		Session: c.server.session,
		Proto:   c.proto,
		Request: request,
		// TODO agent mode, 如果只是 record 那 server 端就不判断了
	}
	call, err := c.server.rpcClient.RequestOutbound(c.ctx, req)
	if err != nil {
		c.server.logger.Debug("Call mock server failed", zap.Error(err))
		return err
	}

	if c.proto.Mock /* TODO && agent mock mode */ {
		return c.handleMockResponse(call, request)
	}

	// record mode
	return c.sendRequestToOriginAndWait(request)
}

func (c *ConnHandler) tryToRecord(request, response *common.GenericMessage) error {
	rpcReq := &rpc.RecordReq{
		Session:    c.server.session,
		IsOutgoing: true, // TODO
		Proto:      c.proto,
		ReqHeader:  request.Header,
		RespHeader: response.Header,
		ReqBody:    request.Body,
		RespBody:   response.Body,
	}
	c.server.rpcClient.RecordRequestResponse(c.ctx, rpcReq)
	return c.responseToInbound(response.Raw)
}

func (c *ConnHandler) handleMockResponse(call *client.Call, request *common.GenericMessage) error {
	<-call.Done
	resp := call.Reply.(*rpc.OutboundResp)

	if resp.OpType == rpc.OP_MOCK {
		return c.responseToInbound(resp.Data)
	}

	c.server.logger.Debug("Send request to origin")
	return c.sendRequestToOriginAndWait(request)
}

func (c *ConnHandler) responseToInbound(data *[]byte) error {
	_, err := c.inboundConn.Write(*data)
	if err != nil {
		c.server.logger.Debug("Write back failed", zap.Error(err))
		return err
	}
	return nil
}

func (c *ConnHandler) sendRequestToOrigin(request *common.GenericMessage) error {
	_, err := c.outboundConn.Write(*request.Raw)
	if err != nil {
		c.server.logger.Debug("Send request to origin failed", zap.Error(err))
	}
	return nil
}

// sendRequestToOriginAndWait 处理 request-response 类型的情况
// 不支持 steam 或者双向通信类型的情况
func (c *ConnHandler) sendRequestToOriginAndWait(request *common.GenericMessage) error {
	_, err := c.outboundConn.Write(*request.Raw)
	if err != nil {
		c.server.logger.Debug("Send request to origin failed", zap.Error(err))
	}

	// 等待 response 并解析
	response, more := <-c.responseC
	if !more {
		return nil
	}
	c.server.logger.Debug("Conn response",
		zap.String("src", c.inboundAddr),
		zap.String("dst", c.outboundAddr),
		zap.String("resp", response.String()))

	err = c.tryToRecord(request, response)
	if err != nil {
		c.server.logger.Debug("Call record server failed", zap.Error(err))
	}
	return err
}

func (c *ConnHandler) normalProxy() {
	go c.pipeTo(c.inboundReader, c.outboundConn, common.ClientToRemote)
	go c.pipeTo(c.outboundReader, c.inboundConn, common.RemoteToClient)
}

func (c *ConnHandler) pipeTo(reader *bufio.Reader, dst *net.TCPConn, direction common.Direction) {
	for {
		_, err := reader.WriteTo(dst)
		if err != nil {
			c.server.logger.Warn("Conn write error",
				zap.String("src", c.inboundAddr),
				zap.String("dst", c.outboundAddr),
				zap.String("dir", direction.String()),
				zap.Error(err))
			break
		}
	}
}

func (c *ConnHandler) handleClose() {
	//req := &rpc.ConnCloseReq{
	//	Source:      c.inboundAddr,
	//	Destination: c.outboundAddr,
	//	Direction:   &direction,
	//}
	//c.server.rpcClient.ConnClose(d.ctx, req)
	c.server.logger.Info("Conn close",
		zap.String("src", c.inboundAddr),
		zap.String("dst", c.outboundAddr))
}
