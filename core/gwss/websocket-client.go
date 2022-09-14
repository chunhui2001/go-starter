package gwss

import (
	"context"
	"net"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type MessageHandler func(ctx context.Context, client *Client, messageBuf []byte)

var (
	logger = config.Log
)

type Client struct {
	ConnectId  string
	ServerAddr string
	CTX        context.Context
	Connection net.Conn
	ReCount    int32
}

func New(connectId string, serverAddress string) *Client {
	return &Client{ConnectId: connectId, ServerAddr: serverAddress}
}

func (c *Client) Connect(messageHandler MessageHandler) (context.Context, net.Conn, error) {

	ctx := context.Background()

	conn, _, _, err := ws.DefaultDialer.Dial(ctx, c.ServerAddr)

	if err != nil {
		logger.Errorf(`Connect-WebSocker-Server-Error: ConnectId=%s, ServerAddress=%s, ErrorMessage=%s`, c.ConnectId, c.ServerAddr, err.Error())
		return nil, nil, err
	}

	c.CTX = ctx
	c.Connection = conn
	// WriteMessage(ctx, conn, "asdf")

	go c.ListenMessage(messageHandler)

	return ctx, conn, nil

}

func (c *Client) ReConnect(messageHandler MessageHandler) (net.Conn, error) {

	conn, _, _, err := ws.DefaultDialer.Dial(c.CTX, c.ServerAddr)

	if err != nil {
		return nil, err
	}

	c.Connection = conn
	go c.ListenMessage(messageHandler)

	return conn, nil

}

func (c *Client) WriteMessage(message string) {

	err := wsutil.WriteClientMessage(c.Connection, ws.OpText, []byte(message))

	if err != nil {
		logger.Errorf(`Write-WebSocker-Message-Error: errorMessage=%s`, err.Error())
	}

}

func (c *Client) ListenMessage(messageHandler MessageHandler) {

	for {

		msg, _, err := wsutil.ReadServerData(c.Connection)

		if err != nil {
			for {
				c.ReCount = c.ReCount + 1
				logger.Errorf(`Write-WebSocker-Message-Error: ConnectId=%s, ReCount=%d, memo=%s, SeverAddress=%s, errorMessage=%s`,
					c.ConnectId, c.ReCount, "Will-be-Reconnect-in-5-sec", c.ServerAddr, err.Error())
				time.Sleep(5 * time.Second) // reconnect in 5 seconds
				c.ReConnect(messageHandler)
			}
		} else {
			if messageHandler != nil {
				go messageHandler(c.CTX, c, msg)
			} else {
				logger.Warnf(`WebSocker-Message-Received-Not-Processed: ConnectId=%s, message=%s`, c.ConnectId, msg)
			}
		}

	}

}
