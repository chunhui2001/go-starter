package gwss

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/gtask"
	"github.com/chunhui2001/go-starter/core/gwss/model"
)

// ulimit -n == 65535
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

var (
	WSS_Conf *config.Wss = config.WssSetting
	Server               = model.NewServer()
)

func init() {

	if WSS_Conf.Enable {
		gtask.AddTask("WebSocket定时任务每秒1次(Ping-and-Pong)", "* * * * * *", func() {
			Server.ServerPing()
			Server.DetectedClientPong()
		})

		gtask.AddTask("WebSocket定时任务每秒15次(打印当前连接数)", "0/15 * * * * *", func() {
			config.Log.Infof(`当前WebSocket维护的所有连接数: count=%d`, len(Server.AllClients()))
		})
	}

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
	Server.NewClient(&client)

	for {

		// Read Message from client
		mt, message, err := ws.ReadMessage()

		if err != nil {
			Server.RemoveClient(client)
			return
		}

		// process messages
		Server.ProcessMessage(client, mt, message)

	}
}
