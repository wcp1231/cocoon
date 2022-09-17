package agent

import (
	"bufio"
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto"
	"context"
	"go.uber.org/zap"
	"net"
)

type conn struct {
	c      *net.TCPConn
	reader *bufio.Reader
	addr   string
}

func newConn(c *net.TCPConn) *conn {
	reader := bufio.NewReader(c)
	addr := c.RemoteAddr().String()
	return &conn{
		c:      c,
		reader: reader,
		addr:   addr,
	}
}

type ConnHandler struct {
	server *Agent
	ctx    context.Context

	inboundConn  *conn
	outboundConn *conn

	pc    *proto.ProtoClassifier
	proto *common.Protocol
}

func NewConnHandler(server *Agent, ctx context.Context, inboundConn, outboundConn *net.TCPConn) *ConnHandler {
	pc := proto.NewProtoClassifier()
	return &ConnHandler{
		server:       server,
		inboundConn:  newConn(inboundConn),
		outboundConn: newConn(outboundConn),
		ctx:          ctx,

		pc: pc,
	}
}

func (c *ConnHandler) Start() {
	defer func() {
		if err := c.inboundConn.c.Close(); err != nil {
			c.server.logger.WithOptions(zap.AddCaller()).Error("proxy inbound conn Close error")
		}
		if err := c.outboundConn.c.Close(); err != nil {
			c.server.logger.WithOptions(zap.AddCaller()).Error("proxy outbound conn Close error")
		}
		c.handleClose()
	}()

	err := c.classifyProto()
	if err != nil {
		c.server.logger.Error("classify protocol error", zap.Error(err))
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
		proxy := newPassHandler(c.server.logger, c.proto, c.inboundConn, c.outboundConn)
		proxy.start()
		return
	}
	if c.proto == common.PROTOCOL_MYSQL {
		proxy := newMysqlHandler(c.ctx, c.server, c.inboundConn, c.outboundConn)
		proxy.start()
		return
	}
	proxy := newRequestResponseHandler(c.ctx, c.server, c.proto, c.inboundConn, c.outboundConn)
	proxy.start()
}

func (c *ConnHandler) classifyProto() error {
	c.proto = c.pc.Classify(c.outboundConn.addr, c.inboundConn.reader)
	return nil
}

//func (c *ConnHandler) handleRequest() {
//	for {
//		select {
//		case request, more := <-c.requestC:
//			if !more {
//				fmt.Printf("%s no more request. %v\n", c.proto.String(), request)
//				return
//			}
//			c.server.logger.Debug("Conn request",
//				zap.String("src", c.inboundConn.addr),
//				zap.String("dst", c.outboundConn.addr),
//				zap.String("proto", c.proto.String()),
//				zap.String("req", request.String()))
//
//			err := c.tryToMock(request)
//			if err != nil {
//				c.server.logger.Debug("Call mock server failed", zap.Error(err))
//				continue
//			}
//		}
//	}
//}

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

//// 请求 server 进行 mock
//func (c *ConnHandler) tryToMock(request *common.GenericMessage) error {
//	req := &rpc.OutboundReq{
//		Session: c.server.session,
//		Proto:   c.proto,
//		Request: request,
//		// TODO agent mode, 如果只是 record 那 server 端就不判断了
//	}
//	call, err := c.server.rpcClient.RequestOutbound(c.ctx, req)
//	if err != nil {
//		c.server.logger.Debug("Call mock server failed", zap.Error(err))
//		return err
//	}
//
//	if c.proto.Mock /* TODO && agent mock mode */ {
//		return c.handleMockResponse(call, request)
//	}
//
//	// record mode
//	return c.sendRequestToOriginAndWait(request)
//}

//func (c *ConnHandler) tryToRecord(request, response *common.GenericMessage) error {
//	rpcReq := &rpc.RecordReq{
//		Session:    c.server.session,
//		IsOutgoing: true, // TODO
//		Proto:      c.proto,
//		ReqHeader:  request.Header,
//		RespHeader: response.Header,
//		ReqBody:    request.Body,
//		RespBody:   response.Body,
//	}
//	c.server.rpcClient.RecordRequestResponse(c.ctx, rpcReq)
//	return c.responseToInbound(response.Raw)
//}

//func (c *ConnHandler) handleMockResponse(call *client.Call, request *common.GenericMessage) error {
//	<-call.Done
//	resp := call.Reply.(*rpc.OutboundResp)
//
//	if resp.OpType == rpc.OP_MOCK {
//		return c.responseToInbound(resp.Data)
//	}
//
//	c.server.logger.Debug("Send request to origin")
//	return c.sendRequestToOriginAndWait(request)
//}
//
//func (c *ConnHandler) responseToInbound(data *[]byte) error {
//	_, err := c.inboundConn.c.Write(*data)
//	if err != nil {
//		c.server.logger.Debug("Write back failed", zap.Error(err))
//		return err
//	}
//	return nil
//}
//
//func (c *ConnHandler) sendRequestToOrigin(request *common.GenericMessage) error {
//	_, err := c.outboundConn.c.Write(*request.Raw)
//	if err != nil {
//		c.server.logger.Debug("Send request to origin failed", zap.Error(err))
//	}
//	return nil
//}
//
//// sendRequestToOriginAndWait 处理 request-response 类型的情况
//// 不支持 steam 或者双向通信类型的情况
//func (c *ConnHandler) sendRequestToOriginAndWait(request *common.GenericMessage) error {
//	_, err := c.outboundConn.c.Write(*request.Raw)
//	if err != nil {
//		c.server.logger.Debug("Send request to origin failed", zap.Error(err))
//	}
//
//	// 等待 response 并解析
//	response, more := <-c.responseC
//	if !more {
//		return nil
//	}
//	c.server.logger.Debug("Conn response",
//		zap.String("src", c.inboundConn.addr),
//		zap.String("dst", c.outboundConn.addr),
//		zap.String("resp", response.String()))
//
//	err = c.tryToRecord(request, response)
//	if err != nil {
//		c.server.logger.Debug("Call record server failed", zap.Error(err))
//	}
//	return err
//}

func (c *ConnHandler) handleClose() {
	//req := &rpc.ConnCloseReq{
	//	Source:      c.inboundAddr,
	//	Destination: c.outboundAddr,
	//	Direction:   &direction,
	//}
	//c.server.rpcClient.ConnClose(d.ctx, req)
	c.server.logger.Info("Conn close",
		zap.String("src", c.inboundConn.addr),
		zap.String("dst", c.outboundConn.addr))
}
