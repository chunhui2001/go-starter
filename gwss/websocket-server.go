package gwss

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	"github.com/chunhui2001/go-starter/cron"
	"github.com/chunhui2001/go-starter/gwss/model"
)

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	//check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
		//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		//return origin == "chrome-extension://cbcbkhdmedgianpaifchdaddpnmgnknn"
		return true
	},
}

var server = &model.Server{}

func init() {

	cron.Add("* * * * * *", func() {
		server.ServerPing()
	})

	cron.Add("* * * * * *", func() {
		server.DetectedClientPong()
	})

}

func WebsocketUpgrade(c *gin.Context) {

	// upgrade get request to websocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	defer ws.Close()

	// create new client & add to client list
	client := model.Client{
		ID:         uuid.Must(uuid.NewRandom()).String(),
		Connection: ws,
		CreatedAt:  time.Now(),
	}

	// greet the new client
	server.NewClient(&client)

	for {

		// Read Message from client
		mt, message, err := ws.ReadMessage()

		if err != nil {
			server.RemoveClient(client)
			return
		}

		// process messages
		server.ProcessMessage(client, mt, message)

	}
}
