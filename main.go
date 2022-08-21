package main

import (
	"net/http"
	"strings"

	"github.com/chunhui2001/go-starter/starter"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/controller"
	_ "github.com/chunhui2001/go-starter/cron"
	_ "github.com/chunhui2001/go-starter/gredis"
	"github.com/chunhui2001/go-starter/mycache"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var APP_COOKIE *config.Cookie = config.CookieSetting

var starterServer = &starter.Server{
	HandlerInfo: func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "this is info page", "message": "Ok"})
	},
	HandlerIndexPage: controller.IndexRouter,
	Handler404: func(c *gin.Context) {
		shortId := strings.Trim(strings.TrimSpace(c.Request.RequestURI), "/")
		if mycache.ShortIdExists(shortId) {
			controller.IndexRouter(c)
		} else {
			if APP_COOKIE.Enable {
				session := sessions.Default(c)
				if session.Get("yourRoomId") != nil && session.Get("yourRoomId") == shortId {
					controller.IndexRouter(c)
				} else {
					c.HTML(http.StatusNotFound, "404", gin.H{
						"content": "Page not found",
					})
				}
			} else {
				c.HTML(http.StatusNotFound, "404", gin.H{
					"content": "Page not found",
				})
			}
		}
	},
}

func main() {
	r := starter.Setup(starterServer)
	r.Run(config.AppSetting.AppPort)
}
