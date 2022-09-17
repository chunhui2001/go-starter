package controller

import (
	"net/http"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/gid"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	WEB_PAGE_CONF *config.WebPageConf = config.WebPageSettings
)

func IndexRouter(c *gin.Context) {

	session := sessions.Default(c)

	if session.Get("yourRoomId") == nil {
		session.Set("yourRoomId", gid.ID())
		session.Save()
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"wssEndpoint": config.WssSetting.Wss(),
		"yourRoomId":  session.Get("yourRoomId"),
		"content":     "This is an Home page...",
	})
}

func AboutRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"content": "This is an about page...",
	})
}
