package main

import (
	"net/http"
	"strings"

	"github.com/chunhui2001/go-starter/actions"
	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/controller"
	"github.com/chunhui2001/go-starter/gredis"
	"github.com/chunhui2001/go-starter/starter"
	"github.com/gin-gonic/gin"
)

var (
	starterServer *starter.Server
	APP_COOKIE    *config.Cookie
	WEB_PAGE_CONF *config.WebPageConf
)

func init() {

	WEB_PAGE_CONF = config.WebPageSettings
	APP_COOKIE = config.CookieSetting

	starterServer = &starter.Server{
		HandlerInfo: func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": "this is info page", "message": "Ok"})
		},
		HandlerIndexPage: controller.IndexRouter,
		Handler404: func(c *gin.Context) {
			if WEB_PAGE_CONF.Enable {
				c.HTML(http.StatusOK, "404", gin.H{
					"requestUrl": c.Request.URL.Path,
					"content":    "Page not found",
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    http.StatusNotFound,
					"message": "Page not found",
				})
			}
		},
	}

}

func main() {

	r := starter.Setup(starterServer)

	Redis_Conf := config.RedisConf

	if Redis_Conf.Mode != gredis.Disabled {
		for _, channel := range strings.Split(Redis_Conf.SubChannels, ",") {
			gredis.Sub(channel, func(channel string, payload string) {
				config.Log.Info("收到了消息1: channel=" + channel + ", payload=" + payload)
			})
		}
	}

	// simples
	r.GET("/httpclient-simple", actions.HttpClientSimpleRouter)
	r.GET("/labs-bigint", actions.BigRouter)
	r.GET("/labs-ytld", actions.YtIdRouter)
	r.GET("/labs-pem", actions.PemRouter)
	r.GET("/labs-leftpad", actions.PadLeftRouter)
	r.POST("/labs-redis-pub", actions.RedisPubRouter)
	r.POST("/labs-upload-file", actions.UploadFileRouterOne)
	r.POST("/demo/album-create", actions.AlbumCreateRouter)
	r.GET("/demo/album-get", actions.AlbumGetRouter)

	r.Run(config.AppSetting.AppPort)

}
