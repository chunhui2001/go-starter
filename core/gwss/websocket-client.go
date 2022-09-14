package gwss

import (
	"context"
	"net"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type MessageHandler func(client *Client, messageBuf []byte)
type OnSuccessHandler func(client *Client)

var (
	logger = config.Log
)

type Client struct {
	ConnectId  string
	ServerAddr string
	Connection net.Conn
	ReCount    int32
}

func NewClient(connectId string, serverAddress string) *Client {
	return &Client{ConnectId: connectId, ServerAddr: serverAddress}
}

func (c *Client) OnSuccess(handler OnSuccessHandler) {
	handler(c)
}

func (c *Client) Connect(ctx context.Context, messageHandler MessageHandler) (net.Conn, error) {

	conn, _, _, err := ws.DefaultDialer.Dial(ctx, c.ServerAddr)

	if err != nil {
		logger.Errorf(`Connect-WebSocker-Server-Error: ConnectId=%s, ServerAddress=%s, ErrorMessage=%s`, c.ConnectId, c.ServerAddr, err.Error())
		return nil, err
	}

	c.Connection = conn

	logger.Infof(`WebSockerClient-Upgrade-Success: ServerAddress=%s, ConnectId=%s`, c.ServerAddr, c.ConnectId)

	go c.ListenMessage(ctx, messageHandler)

	return conn, nil

}

func (c *Client) ReConnect(ctx context.Context, messageHandler MessageHandler) (net.Conn, error) {

	conn, _, _, err := ws.DefaultDialer.Dial(ctx, c.ServerAddr)

	if err != nil {
		return nil, err
	}

	c.Connection = conn
	go c.ListenMessage(ctx, messageHandler)

	logger.Infof(`WebSockerClient-Upgrade-Success: ServerAddress=%s, ConnectId=%s, ReCount=%d`, c.ServerAddr, c.ConnectId, c.ReCount)

	return conn, nil

}

func (c *Client) WriteMessage(message string) {

	err := wsutil.WriteClientMessage(c.Connection, ws.OpText, []byte(message))

	if err != nil {
		logger.Errorf(`Write-WebSocker-Message-Error: errorMessage=%s`, err.Error())
	}

}

func (c *Client) ListenMessage(ctx context.Context, messageHandler MessageHandler) {

	for {

		msg, _, err := wsutil.ReadServerData(c.Connection)

		if err != nil {
			for {
				c.ReCount = c.ReCount + 1
				logger.Errorf(`Write-WebSocker-Message-Error: ConnectId=%s, ReCount=%d, memo=%s, SeverAddress=%s, errorMessage=%s`,
					c.ConnectId, c.ReCount, "Will-be-Reconnect-in-5-sec", c.ServerAddr, err.Error())
				time.Sleep(5 * time.Second) // reconnect in 5 seconds
				c.ReConnect(ctx, messageHandler)
				c.Connection.Close()
			}
		} else {
			if messageHandler != nil {
				go messageHandler(c, msg)
			} else {
				logger.Warnf(`WebSocker-Message-Received-Not-Processed: ConnectId=%s, message=%s`, c.ConnectId, msg)
			}
		}

	}

}
