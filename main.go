package main

import (
	"net/http"
	"strings"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/controller"
	"github.com/chunhui2001/go-starter/gredis"
	"github.com/chunhui2001/go-starter/starter"
	"github.com/gin-gonic/gin"
)

var (
	starterServer *starter.Server
	APP_COOKIE    *config.Cookie
)

func init() {

	APP_COOKIE = config.CookieSetting

	starterServer = &starter.Server{
		HandlerInfo: func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": "this is info page", "message": "Ok"})
		},
		HandlerIndexPage: controller.IndexRouter,
		Handler404: func(c *gin.Context) {
			c.HTML(http.StatusNotFound, "404", gin.H{
				"content": "Page not found",
			})
		},
	}

}

func main() {

	r := starter.Setup(starterServer)

	Redis_Conf := config.RedisConf

	for _, channel := range strings.Split(Redis_Conf.SubChannels, ",") {
		gredis.Sub(channel, func(channel string, payload string) {
			config.Log.Info("收到了消息: channel=" + channel + ", payload=" + payload)
		})
	}

	r.Run(config.AppSetting.AppPort)

}
