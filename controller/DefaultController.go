package controller

import (
	"net/http"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/mycache"
	"github.com/chunhui2001/go-starter/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IndexRouter(c *gin.Context) {

	APP_COOKIE := config.CookieSetting
	var yourRoomId string

	if APP_COOKIE.Enable {
		session := sessions.Default(c)
		if session.Get("yourRoomId") == nil {
			session.Set("yourRoomId", utils.ShortId())
			session.Save()
		}
		yourRoomId = utils.ToString(session.Get("yourRoomId"))
	} else {
		yourRoomId = mycache.ShortIdPut(utils.ShortId())
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"wssEndpoint": config.WssSetting.Wss(),
		"yourRoomId":  yourRoomId,
		"content":     "This is an Home page...",
	})
}

func AboutRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"content": "This is an about page...",
	})
}
