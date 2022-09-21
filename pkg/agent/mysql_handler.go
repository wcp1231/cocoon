package agent

import (
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto"
	"cocoon/pkg/proto/mysql"
	"context"
	"go.uber.org/zap"
)

// mysqlHandler 处理 mysql 的连接
type mysqlHandler struct {
	server       *Agent
	ctx          context.Context
	proto        *common.Protocol
	inboundConn  *conn
	outboundConn *conn

	dissector *mysql.Dissector
	requestC  chan common.Message
	responseC chan common.Message
}

func newMysqlHandler(ctx context.Context, server *Agent, inbound, outbound *conn) *mysqlHandler {
	requestC := make(chan common.Message)
	responseC := make(chan common.Message)
	return &mysqlHandler{
		server:       server,
		ctx:          ctx,
		proto:        common.PROTOCOL_MYSQL,
		inboundConn:  inbound,
		outboundConn: outbound,

		dissector: proto.NewMysqlDissector(requestC, responseC),
		requestC:  requestC,
		responseC: responseC,
	}
}

func (c *mysqlHandler) start() {
	c.server.logger.Info("Mysql try hand shake",
		zap.String("src", c.inboundConn.addr),
		zap.String("dst", c.outboundConn.addr))
	c.dissector.Init(c.inboundConn.reader, c.outboundConn.reader)

	err := c.handleHandShake()
	if err != nil {
		// TODO
		c.server.logger.Error("Mysql hand shake failed",
			zap.String("src", c.inboundConn.addr),
			zap.String("dst", c.outboundConn.addr),
			zap.Error(err))
		return
	}

	c.server.logger.Info("Mysql Conn proxy",
		zap.String("src", c.inboundConn.addr),
		zap.String("dst", c.outboundConn.addr))
	go c.handleRequest()
	go c.handleRawResponse()
	go c.dissector.StartRequestDissect(c.inboundConn.reader)
	go c.dissector.StartResponseDissect(c.outboundConn.reader)
}

func (c *mysqlHandler) handleHandShake() error {
	// parse server greeting packet
	data, err := c.dissector.ReadServerHandshake()
	if err != nil {
		return err
	}
	err = c.sendToClient(data)
	if err != nil {
		return err
	}

	// parse client auth packet
	data, err = c.dissector.ReadClientHandshakeResponse()
	if err != nil {
		return err
	}
	err = c.sendToServer(data)
	if err != nil {
		return err
	}

	// parse server auth response
	// ok packet to client handshake response
	data, err = c.dissector.ReadPacketFromServer()
	if err != nil {
		return err
	}

	return c.sendToClient(data)
}

func (c *mysqlHandler) handleRequest() {
	for {
		select {
		case request, more := <-c.requestC:
			if !more {
				return
			}
			request.SetId(c.server.nextId())
			c.server.logger.Debug("Mysql Conn request",
				zap.String("src", c.inboundConn.addr),
				zap.String("dst", c.outboundConn.addr),
				zap.String("req", request.String()))
			err := c.tryToMock(request)
			if err != nil {
				c.server.logger.Debug("Call mock server failed", zap.Error(err))
				continue
			}
		}
	}
}

func (c *mysqlHandler) handleRawResponse() {
	for {
		select {
		case response, more := <-c.responseC:
			if !more {
				return
			}
			c.server.logger.Debug("Mysql Conn request",
				zap.String("src", c.inboundConn.addr),
				zap.String("dst", c.outboundConn.addr),
				zap.String("resp", response.String()))

			mysqlResp := response.(*mysql.MysqlMessage)
			err := c.handleResponse(mysqlResp.GetRequest(), mysqlResp)
			if err != nil {
				c.server.logger.Warn("Mysql send to client failed", zap.Error(err))
				continue
			}
		}
	}
}

// tryToMock 由 agent 进行 mock
// 根据协议和 agent 配置判断是否进行 mock
func (c *mysqlHandler) tryToMock(request common.Message) error {
	// record request
	c.server.recordServer.RecordRequest(request)

	if c.proto.Mock /* TODO && agent mock mode */ {
		return c.requestMockServer(request)
	}

	return c.sendToServer(*request.GetRaw())
}

// requestMockServer 获取 mock 结果
func (c *mysqlHandler) requestMockServer(request common.Message) error {
	result := c.server.mockServer.Mock(c.proto.Name, request)

	if result.Pass {
		c.server.logger.Debug("Send request to origin")
		return c.sendToServer(*request.GetRaw())
	}

	return c.handleResponse(request, result.Data)
}

// handleResponse 处理 response
// 无论是 mock 还是真实数据
// 主要功能暂时只有记录
func (c *mysqlHandler) handleResponse(request, response common.Message) error {
	// record response
	c.server.recordServer.RecordResponse(request, response)
	return c.sendToClient(*response.GetRaw())
}

func (c *mysqlHandler) sendToClient(data []byte) error {
	_, err := c.inboundConn.c.Write(data)
	if err != nil {
		c.server.logger.Debug("Write to client failed", zap.Error(err))
		return err
	}
	return nil
}

func (c *mysqlHandler) sendToServer(data []byte) error {
	_, err := c.outboundConn.c.Write(data)
	if err != nil {
		c.server.logger.Debug("Write to server failed", zap.Error(err))
		return err
	}

	return nil
}
