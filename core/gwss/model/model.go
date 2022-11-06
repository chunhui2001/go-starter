package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/gid"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gorilla/websocket"
)

var (
	logger              = config.Log
	WSSConf *config.Wss = config.WssSetting
)

// contant for 4 type actions
const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
	pong        = "pong"
)

const (
	server_ping = "server_ping"
)

// a server type to store all subscriptions
type Server struct {
	Subscriptions []Subscription
	lock          sync.Mutex
}

func NewServer() *Server {
	return &Server{}
}

// each subscription consists of topic-name & client
type Subscription struct {
	Topic   string
	Clients *[]Client
}

// each client consists of auto-generated ID & connection
type Client struct {
	ID         string
	Connection *websocket.Conn
	LastPong   time.Time
	CreatedAt  time.Time
}

// type for a valid message.
type Message struct {
	Id      string      `json:"id"`
	Action  string      `json:"action"`
	Topic   string      `json:"topic"`
	Message interface{} `json:"message"`
	Time    string      `json:"time"`
}

func NewMessage(topic string, action string, message interface{}) *Message {
	return &Message{
		Id:      gid.ID(),
		Action:  action,
		Topic:   topic,
		Message: message,
		Time:    utils.DateTimeUTCString(),
	}
}

func (m *Message) Bytes() []byte {
	return utils.ToJsonBytes(m)
}

func (s *Server) ServerPing() {
	s.Publish(NewMessage(server_ping, "ping", utils.DateTimeUTCString()))
}

func (s *Server) DetectedClientPong() {

	var clients []Client = s.AllClients()

	d1 := 45 * time.Second
	d2 := 15 * time.Second

	for _, client := range clients {

		var t time.Duration

		if client.LastPong.IsZero() {
			t = time.Since(client.CreatedAt)
		} else {
			t = time.Since(client.LastPong)
		}

		if t >= d1 {
			s.Send(&client, NewMessage("sys", "connection_closed", `Your connection has been closed, Bye.`).Bytes())
			s.RemoveClient(client)
		} else if t >= d2 {
			s.Send(&client, NewMessage("sys", "connection_warnning", fmt.Sprintf(`Your connection will be closed, Please send Pong in '%s'`, d1-t)).Bytes())
		}

	}

}

func (s *Server) AllClients() []Client {

	var clients []Client

	for _, sub := range s.Subscriptions {
		if sub.Topic == server_ping {
			clients = append(clients, *sub.Clients...)
		}
	}

	return clients
}

func (s *Server) ReceiveClientPong(client *Client, message string) {
	for _, sub := range s.Subscriptions {
		if sub.Topic == server_ping {
			clients := *sub.Clients
			for j, c := range clients {
				if c.ID == client.ID {
					clients[j].LastPong = utils.DateTimeParse(message)
					break
				}
			}
			break
		}
	}
}

func (s *Server) NewClient(client *Client) {
	s.Subscribe(client, server_ping)
	s.Send(client, NewMessage("sys", "connected_successful", fmt.Sprintf(`Welcome! Your ID is: '%s'`, client.ID)).Bytes())
}

func (s *Server) Send(client *Client, messageBytes []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	err := client.Connection.WriteMessage(1, messageBytes)
	if err != nil {
		logger.Errorf("WebSocket-Send-Error: ErrorMessage=%V", err)
	}
}

func (s *Server) ProcessMessage(client Client, messageType int, payload []byte) *Server {

	m := Message{}

	if err := json.Unmarshal(payload, &m); err != nil {
		s.Send(&client, []byte("Server-Invalid-Message: Message="+string(payload)+"ErrorMessage="+err.Error()))
		return s
	}

	switch m.Action {
	case publish:
		// s.Publish(m.Topic, "publish", m.Message)
		s.Publish(NewMessage(m.Topic, "publish", m.Message))
		return s
	case subscribe:
		s.Subscribe(&client, m.Topic)
		return s
	case unsubscribe:
		s.Unsubscribe(&client, m.Topic)
		return s
	case pong:
		s.ReceiveClientPong(&client, m.Message.(string))
		return s
	default:
		s.Send(&client, []byte("Server: Action unrecognized"))
		return s
	}

}

func (s *Server) Publish(message *Message) {

	if message == nil {
		return
	}

	var clients []Client

	// get list of clients subscribed to topic
	for _, sub := range s.Subscriptions {
		if sub.Topic == message.Topic {
			clients = append(clients, *sub.Clients...)
		}
	}

	if len(clients) != 0 {
		// send to clients
		for _, client := range clients {
			s.Send(&client, message.Bytes())
			if WSSConf.PrintMessage {
				logger.Debugf(`Wss广播了一条消息: topic=%s, messageSize=%d, clientId=%s`, message.Topic, len(message.Bytes()), client.ID)
			}
		}
	} else {
		if WSSConf.PrintMessage {
			logger.Debugf("no-have-clients-to-be-subscribe: topic=" + message.Topic)
		}
	}

}

func (s *Server) Subscribe(client *Client, topic string) {

	exist := false

	// find existing topics
	for _, sub := range s.Subscriptions {
		// if found, add client
		if sub.Topic == topic {
			exist = true
			*sub.Clients = append(*sub.Clients, *client)
		}
	}

	// else, add new topic & add client to that topic
	if !exist {

		newClient := &[]Client{*client}

		newTopic := &Subscription{
			Topic:   topic,
			Clients: newClient,
		}

		s.Subscriptions = append(s.Subscriptions, *newTopic)
	}
}

func (s *Server) Unsubscribe(client *Client, topic string) {
	// Read all topics
	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			// Read all topics' client
			for i := 0; i < len(*sub.Clients); i++ {
				if client.ID == (*sub.Clients)[i].ID {
					// If found, remove client
					if i == len(*sub.Clients)-1 {
						// if it's stored as the last element, crop the array length
						*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
					} else {
						// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
						*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
						i--
					}
				}
			}
		}
	}
}

func (s *Server) RemoveClient(client Client) {
	// Read all subs
	for _, sub := range s.Subscriptions {
		// Read all client
		for i := 0; i < len(*sub.Clients); i++ {
			if client.ID == (*sub.Clients)[i].ID {
				// If found, remove client
				if i == len(*sub.Clients)-1 {
					// if it's stored as the last element, crop the array length
					*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
				} else {
					// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
					*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
					i--
				}
			}
		}
	}
}
