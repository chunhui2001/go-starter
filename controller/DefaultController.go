package controller

import (
	"net/http"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/mycache"
	"github.com/chunhui2001/go-starter/utils"
	"github.com/gin-gonic/gin"
)

func IndexRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"wssEndpoint": config.WssSetting.Wss(),
		"yourRoomId":  mycache.ShortIdPut(utils.ShortId()),
		"content":     "This is an Home page...",
	})
}

func AboutRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"content": "This is an about page...",
	})
}
