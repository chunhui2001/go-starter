package gwss

import (
	"bytes"
	"compress/gzip"
	"context"
	"io/ioutil"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/alitto/pond"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type MessageHandler func(client *Client, opcode string, messageBuf []byte)
type SuccessHandler func(client *Client)

var (
	logger      = config.Log
	ALL_CLIENTS = map[string]*Client{}
)

type Client struct {
	ConnectId        string
	ServerAddr       string
	Connection       net.Conn
	ReCount          int32
	OnSuccessHandler SuccessHandler
	AlreadyClosed    bool
	WillClosed       bool
}

func NewClient(connectId string, serverAddress string) *Client {
	ALL_CLIENTS[connectId] = &Client{ConnectId: connectId, ServerAddr: serverAddress}
	return ALL_CLIENTS[connectId]
}

func Shutdown(connectId string) bool {
	if ALL_CLIENTS[connectId] != nil {
		ALL_CLIENTS[connectId].WillClosed = true
		ALL_CLIENTS[connectId].Connection.Close()
		ALL_CLIENTS[connectId].AlreadyClosed = true
		delete(ALL_CLIENTS, connectId)
		return true
	}
	return false
}

func (c *Client) OnSuccess(successHandler SuccessHandler) *Client {
	c.OnSuccessHandler = successHandler
	return c
}

func (c *Client) Connect(messageHandler MessageHandler) (*Client, error) {

	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), c.ServerAddr)

	if err != nil {
		logger.Errorf(`Connect-WebSocker-Server-Error: ConnectId=%s, ServerAddress=%s, ErrorMessage=%s`, c.ConnectId, c.ServerAddr, err.Error())
		return nil, err
	}

	c.Connection = conn

	go c.ListenMessage(messageHandler)

	logger.Infof(`WebSockerClient-Upgrade-Success: ConnectId=%s, ReCount=%d, ServerAddress=%s`, c.ConnectId, c.ReCount, c.ServerAddr)

	if c.OnSuccessHandler != nil {
		c.OnSuccessHandler(c)
	}

	return c, nil

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

	numCPUs := runtime.NumCPU()
	pool := pond.New(numCPUs, 1000)
	defer pool.StopAndWait()

	logger.Infof(`WebSocket-Listener-Message-As-a-Pool: ConnectId=%s, SeverAddress=%s`, c.ConnectId, c.ServerAddr)

	for {

		msg, opcode, err := wsutil.ReadServerData(c.Connection)

		if err != nil {

			if strings.Contains(err.Error(), "connection reset by peer") {

				c.ReCount = c.ReCount + 1
				// 重建创建连接
				logger.Errorf(`WebSocket-Connection-Has-Been-Closed: opcode=%x, ConnectId=%s, ReCount=%d, memo=%s, SeverAddress=%s, errorMessage=%s`,
					opcode, c.ConnectId, c.ReCount, "Will-be-Reconnect-in-2-sec", c.ServerAddr, err.Error())

			} else {
				if !c.WillClosed {
					logger.Errorf(`WebSocket-Connection-Error-Will-Be-Closed: opcode=%x, ConnectId=%s, ReCount=%d, memo=%s,, SeverAddress=%s, errorMessage=%s`,
						opcode, c.ConnectId, c.ReCount, "Will-be-Reconnect-in-2-sec", c.ServerAddr, err.Error())
				}
			}

			if c.WillClosed || c.AlreadyClosed {
				if c.AlreadyClosed {
					logger.Infof(`WebSocket-Connection-Already-Closed: ConnectId=%s, ServerAddress=%s`, c.ConnectId, c.ServerAddr)
				} else if c.WillClosed {
					logger.Infof(`WebSocket-Connection-Will-Closed: ConnectId=%s, ServerAddress=%s`, c.ConnectId, c.ServerAddr)
					continue
				}
			} else {
				c.Connection.Close()
				time.Sleep(2 * time.Second) // reconnect in 2 seconds
				for {
					if _, err := c.Connect(messageHandler); err == nil {
						// 重连成功
						break
					} else {
						// 链接失败再次重新创建链接
						time.Sleep(2 * time.Second) // reconnect in 2 seconds
					}
				}
			}

			break

		} else {
			if messageHandler != nil {
				pool.Submit(func() {
					defer func() {
						if r := recover(); r != nil {
							if msg[0] == 31 && msg[1] == 139 {
								// gzip
								buf := bytes.NewBuffer(msg)
								reader, err1 := gzip.NewReader(buf)
								if err1 == nil {
									s, err2 := ioutil.ReadAll(reader)
									if err2 == nil {
										logger.Errorf("WebSocker-MessageHandler-Recovered: msg=%s, Error=%v", s, r)
									} else {
										logger.Errorf("WebSocker-MessageHandler-Recovered: msg=%s, Error=%v", string(msg), r)
									}
								} else {
									logger.Errorf("WebSocker-MessageHandler-Recovered: msg=%s, Error=%v", string(msg), r)
								}
							} else {
								logger.Errorf("WebSocker-MessageHandler-Recovered: msg=%s, Error=%v", string(msg), r)
							}
						}
					}()
					messageHandler(c, utils.ToString(opcode), msg)
				})
			} else {
				logger.Warnf(`WebSocker-Message-Received-Not-Processed: ConnectId=%s, message=%s`, c.ConnectId, msg)
			}
		}

	}

}
