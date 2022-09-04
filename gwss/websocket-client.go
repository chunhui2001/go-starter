package gwss

import (
	"context"
	"net"

	"github.com/chunhui2001/go-starter/config"
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
}

func New(connectId string, serverAddress string) *Client {
	return &Client{ConnectId: connectId, ServerAddr: serverAddress}
}

func (c *Client) Connect(messageHandler MessageHandler) (context.Context, net.Conn, error) {

	ctx := context.Background()

	conn, _, _, err := ws.DefaultDialer.Dial(ctx, c.ServerAddr)

	if err != nil {
		logger.Errorf(`Connect-WebSocker-Server-Error: serverAddress=%s, errorMessage=%s`, c.ServerAddr, err.Error())
		return nil, nil, err
	}

	c.CTX = ctx
	c.Connection = conn
	// WriteMessage(ctx, conn, "asdf")

	go c.ListenMessage(messageHandler)

	return ctx, conn, nil

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
			logger.Errorf(`Write-WebSocker-Message-Error: errorMessage=%s`, err.Error())
		} else {
			// logger.Debugf(`WebSocker-Message-Received: message=%s`, msg)
			if messageHandler != nil {
				go messageHandler(c.CTX, c, msg)
			}
		}

	}

}
