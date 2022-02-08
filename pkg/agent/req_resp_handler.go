package agent

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto"
	"context"
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
			request.Id = c.server.nextId()
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

// tryToMock 由 agent 进行 mock
// 根据协议和 agent 配置判断是否进行 mock
func (c *requestResponseHandler) tryToMock(request *common.GenericMessage) error {
	// record request
	c.server.recordServer.RecordRequest(request)

	if c.proto.Mock /* TODO && agent mock mode */ {
		return c.requestMockServer(request)
	}

	return c.sendRequestToOriginAndWait(request)
}

// requestMockServer 获取 mock 结果
func (c *requestResponseHandler) requestMockServer(request *common.GenericMessage) error {
	result := c.server.mockServer.Mock(c.proto.Name, request)

	if result.Pass {
		c.server.logger.Debug("Send request to origin")
		return c.sendRequestToOriginAndWait(request)
	}

	return c.handleResponse(request, result.Data)
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
	err = c.handleResponse(request, response)
	if err != nil {
		c.server.logger.Debug("Call record server failed", zap.Error(err))
	}
	return err
}

// handleResponse 处理 response
// 无论是 mock 还是真实数据
// 主要功能暂时只有记录
func (c *requestResponseHandler) handleResponse(request, response *common.GenericMessage) error {
	// record response
	c.server.recordServer.RecordResponse(request, response)
	return c.sendResponseToInbound(response.Raw)
}

func (c *requestResponseHandler) sendResponseToInbound(data *[]byte) error {
	_, err := c.inboundConn.c.Write(*data)
	if err != nil {
		c.server.logger.Debug("Write back failed", zap.Error(err))
		return err
	}
	return nil
}
