package agent

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/model/rpc"
	"cocoon/pkg/proto"
	"context"
	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
)

// requestResponseHandler 处理 request-response 类型的连接
type requestResponseHandler struct {
	server       *Agent
	ctx          context.Context
	proto        *common.Protocol
	inboundConn  *conn
	outboundConn *conn

	requestC  chan *common.GenericMessage
	responseC chan *common.GenericMessage
}

func newRequestResponseHandler(ctx context.Context, server *Agent, proto *common.Protocol, inbound, outbound *conn) *requestResponseHandler {
	requestC := make(chan *common.GenericMessage)
	responseC := make(chan *common.GenericMessage)
	return &requestResponseHandler{
		server:       server,
		ctx:          ctx,
		proto:        proto,
		inboundConn:  inbound,
		outboundConn: outbound,

		requestC:  requestC,
		responseC: responseC,
	}
}

func (c *requestResponseHandler) start() {
	c.server.logger.Info("Conn proxy",
		zap.String("src", c.inboundConn.addr),
		zap.String("dst", c.outboundConn.addr),
		zap.String("proto", c.proto.Name))

	requestDissector := proto.NewRequestDissector(c.proto, c.requestC, c.responseC)
	go c.handleRequest()
	go requestDissector.StartRequestDissect(c.inboundConn.reader)
	go requestDissector.StartResponseDissect(c.outboundConn.reader)
}

func (c *requestResponseHandler) handleRequest() {
	for {
		select {
		case request, more := <-c.requestC:
			if !more {
				return
			}
			c.server.logger.Debug("Conn request",
				zap.String("src", c.inboundConn.addr),
				zap.String("dst", c.outboundConn.addr),
				zap.String("proto", c.proto.String()),
				zap.String("req", request.String()))

			err := c.tryToMock(request)
			if err != nil {
				c.server.logger.Debug("Call mock server failed", zap.Error(err))
				continue
			}
		}
	}
}

// 由 agent 进行 mock
func (c *requestResponseHandler) tryToMock(request *common.GenericMessage) error {
	if c.proto.Mock /* TODO && agent mock mode */ {
		return c.handleMock(request)
	}

	// record mode
	return c.sendRequestToOriginAndWait(request)
}

func (c *requestResponseHandler) tryToRecord(request, response *common.GenericMessage) error {
	//rpcReq := &rpc.RecordReq{
	//	Session:    c.server.session,
	//	IsOutgoing: true, // TODO
	//	Proto:      c.proto,
	//	ReqHeader:  request.Header,
	//	RespHeader: response.Header,
	//	ReqBody:    request.Body,
	//	RespBody:   response.Body,
	//}
	//c.server.rpcClient.RecordRequestResponse(c.ctx, rpcReq)
	return c.responseToInbound(response.Raw)
}

func (c *requestResponseHandler) handleMock(request *common.GenericMessage) error {
	result := c.server.mockServer.Mock(c.proto.Name, request)

	if result.Pass {
		c.server.logger.Debug("Send request to origin")
		return c.sendRequestToOriginAndWait(request)

	}

	return c.responseToInbound(result.Data)
}

func (c *requestResponseHandler) handleMockResponse(call *client.Call, request *common.GenericMessage) error {
	<-call.Done
	resp := call.Reply.(*rpc.OutboundResp)

	if resp.OpType == rpc.OP_MOCK {
		return c.responseToInbound(resp.Data)
	}

	c.server.logger.Debug("Send request to origin")
	return c.sendRequestToOriginAndWait(request)
}

func (c *requestResponseHandler) responseToInbound(data *[]byte) error {
	_, err := c.inboundConn.c.Write(*data)
	if err != nil {
		c.server.logger.Debug("Write back failed", zap.Error(err))
		return err
	}
	return nil
}

// sendRequestToOriginAndWait 处理 request-response 类型的情况
// 不支持 steam 或者双向通信类型的情况
func (c *requestResponseHandler) sendRequestToOriginAndWait(request *common.GenericMessage) error {
	_, err := c.outboundConn.c.Write(*request.Raw)
	if err != nil {
		c.server.logger.Debug("Send request to origin failed", zap.Error(err))
	}

	// 等待 response 并解析
	response, more := <-c.responseC
	if !more {
		return nil
	}
	c.server.logger.Debug("Conn response",
		zap.String("src", c.inboundConn.addr),
		zap.String("dst", c.outboundConn.addr),
		zap.String("resp", response.String()))

	err = c.tryToRecord(request, response)
	if err != nil {
		c.server.logger.Debug("Call record server failed", zap.Error(err))
	}
	return err
}
