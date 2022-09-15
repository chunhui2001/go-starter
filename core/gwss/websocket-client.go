package gwss

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type MessageHandler func(client *Client, opcode string, messageBuf []byte)
type SuccessHandler func(client *Client)

var (
	logger = config.Log
)

type Client struct {
	ConnectId        string
	ServerAddr       string
	Connection       net.Conn
	ReCount          int32
	OnSuccessHandler SuccessHandler
}

func NewClient(connectId string, serverAddress string) *Client {
	return &Client{ConnectId: connectId, ServerAddr: serverAddress}
}

func (c *Client) OnSuccess(successHandler SuccessHandler) *Client {
	c.OnSuccessHandler = successHandler
	return c
}

func (c *Client) Connect(messageHandler MessageHandler) (net.Conn, error) {

	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), c.ServerAddr)

	if err != nil {
		logger.Errorf(`Connect-WebSocker-Server-Error: ConnectId=%s, ServerAddress=%s, ErrorMessage=%s`, c.ConnectId, c.ServerAddr, err.Error())
		return nil, err
	}

	c.Connection = conn

	go c.ListenMessage(messageHandler)

	logger.Infof(`WebSockerClient-Upgrade-Success: ServerAddress=%s, ConnectId=%s`, c.ServerAddr, c.ConnectId)

	if c.OnSuccessHandler != nil {
		c.OnSuccessHandler(c)
	}

	return conn, nil

}

func (c *Client) ReConnect(messageHandler MessageHandler) (net.Conn, error) {

	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), c.ServerAddr)

	if err != nil {
		return nil, err
	}

	c.Connection = conn

	logger.Infof(`WebSockerClient-ReConnect-Success: ServerAddress=%s, ConnectId=%s, ReCount=%d`, c.ServerAddr, c.ConnectId, c.ReCount)

	return conn, nil

}

func (c *Client) WriteMessage(message string) {

	err := wsutil.WriteClientMessage(c.Connection, ws.OpText, []byte(message))

	if err != nil {
		logger.Errorf(`Write-WebSocker-Message-Error: errorMessage=%s`, err.Error())
	}

}

/*
//	|Opcode  | Meaning                             | Reference |
// -+--------+-------------------------------------+-----------|
//	| 0      | Continuation Frame                  | RFC 6455  |
// -+--------+-------------------------------------+-----------|
//	| 1      | Text Frame                          | RFC 6455  |
// -+--------+-------------------------------------+-----------|
//	| 2      | Binary Frame                        | RFC 6455  |
// -+--------+-------------------------------------+-----------|
//	| 8      | Connection Close Frame              | RFC 6455  |
// -+--------+-------------------------------------+-----------|
//	| 9      | Ping Frame                          | RFC 6455  |
// -+--------+-------------------------------------+-----------|
//	| 10     | Pong Frame                          | RFC 6455  |
// -+--------+-------------------------------------+-----------|
*/
func (c *Client) ListenMessage(messageHandler MessageHandler) {

	for {

		msg, opcode, err := wsutil.ReadServerData(c.Connection)

		if err != nil {

			if strings.Contains(err.Error(), "connection reset by peer") {

				// 重建创建连接
				logger.Errorf(`WebSocket-Connection-Has-Been-Closed: ConnectId=%s, opcode=%x, SeverAddress=%s, errorMessage=%s`, c.ConnectId, opcode, c.ServerAddr, err.Error())

				if err2 := c.Connection.Close(); err2 != nil {
					logger.Errorf(`WebSocket-Closed-Error: ConnectId=%s, ReCount=%d, SeverAddress=%s, errorMessage=%s`,
						c.ConnectId, c.ReCount, c.ServerAddr, err2.Error())
				}

				break
			}

			for {

				c.ReCount = c.ReCount + 1

				logger.Errorf(`Read-WebSocket-Message-Error: ConnectId=%s, ReCount=%d, memo=%s, SeverAddress=%s, errorMessage=%s`,
					c.ConnectId, c.ReCount, "Will-be-Reconnect-in-5-sec", c.ServerAddr, err.Error())

				time.Sleep(5 * time.Second) // reconnect in 5 seconds
				c.ReConnect(messageHandler)

			}

		} else {
			if messageHandler != nil {

				if opcode == 0x9 {
					logger.Infof(`WebSocker-Receive-Ping: ConnectId=%s, opcode=%s, message=%s`, c.ConnectId, utils.ToString(opcode), msg)
				} else if opcode == 0x8 {
					logger.Infof(`WebSocker-Receive-Closed: ConnectId=%s, opcode=%s, message=%s`, c.ConnectId, utils.ToString(opcode), msg)
				} else if opcode == 0xa {
					logger.Infof(`WebSocker-Receive-Pong: ConnectId=%s, opcode=%s, message=%s`, c.ConnectId, utils.ToString(opcode), msg)
				}

				go messageHandler(c, utils.ToString(opcode), msg)

			} else {
				logger.Warnf(`WebSocker-Message-Received-Not-Processed: ConnectId=%s, message=%s`, c.ConnectId, msg)
			}
		}

	}

}
