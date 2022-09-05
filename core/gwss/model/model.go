package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gorilla/websocket"
)

var (
	logger = config.Log
)

// contant for 4 type actions
const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
	pong        = "pong"
)

const (
	server_ping                 = "server_ping"
	action_connected_successful = "connected_successful"
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
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func NewMessage(topic string, action string, message string) *Message {
	return &Message{
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
	s.Publish(server_ping, utils.DateTimeUTCString())
}

func (s *Server) DetectedClientPong() {

	var clients []Client = s.AllClients()

	d1 := 45 * time.Second
	d2 := 15 * time.Second

	for _, client := range clients {

		var t time.Duration

		if client.LastPong.IsZero() {
			t = time.Now().Sub(client.CreatedAt)
		} else {
			t = time.Now().Sub(client.LastPong)
		}

		if t >= d1 {
			s.Send(&client, NewMessage("sys", "connection_closed", `Your connection has been closed, Bye.`).Bytes())
			s.RemoveClient(client)
		} else if t >= d2 {
			s.Send(&client, NewMessage("sys", "connection_alert", fmt.Sprintf(`Your connection will be closed, Please send Pong in '%s'`, d1-t)).Bytes())
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
	s.Send(client, NewMessage("sys", action_connected_successful, fmt.Sprintf(`Welcome! Your ID is: '%s'`, client.ID)).Bytes())
}

func (s *Server) Send(client *Client, messageBytes []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	client.Connection.WriteMessage(1, messageBytes)
}

func (s *Server) ProcessMessage(client Client, messageType int, payload []byte) *Server {

	m := Message{}

	if err := json.Unmarshal(payload, &m); err != nil {
		s.Send(&client, []byte("Server: Invalid payload"))
		return s
	}

	switch m.Action {
	case publish:
		s.Publish(m.Topic, m.Message)
		break
	case subscribe:
		s.Subscribe(&client, m.Topic)
		break
	case unsubscribe:
		s.Unsubscribe(&client, m.Topic)
		break
	case pong:
		s.ReceiveClientPong(&client, m.Message)
		break
	default:
		s.Send(&client, []byte("Server: Action unrecognized"))
		break
	}

	return s
}

func (s *Server) Publish(topic string, message string) {

	var clients []Client

	// get list of clients subscribed to topic
	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			clients = append(clients, *sub.Clients...)
		}
	}

	if len(clients) != 0 {
		// send to clients
		for _, client := range clients {
			m := utils.MapOf("topic", topic, "message", message)
			s.Send(&client, []byte(utils.ToJsonString(m)))
			// logger.Log.Info(topic + ": " + message + ", clientId: " + client.ID)
		}
	} else {
		//logger.Log.Info("no-have-clients-to-be-subscribe: topic=" + topic)
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
