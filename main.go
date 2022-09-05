package main

import (
	"net/http"
	"strings"

	"github.com/chunhui2001/go-starter/actions"
	. "github.com/chunhui2001/go-starter/commons"
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
			c.JSON(http.StatusOK, R{Data: "this is info page"}.Success())
		},
		HandlerIndexPage: controller.IndexRouter,
		Handler404: func(c *gin.Context) {
			if WEB_PAGE_CONF.Enable {
				c.HTML(http.StatusOK, "404", gin.H{
					"requestUrl": c.Request.URL.Path,
					"content":    "Page not found",
				})
			} else {
				c.JSON(http.StatusOK, R{}.Msg("Page-Not-Found").Fail(404))
			}
		},
	}

}

func main() {

	// simples
	starter.AppendRouter("GET", []string{"/httpclient-simple"}, actions.HttpClientSimpleRouter)
	starter.AppendRouter("GET", []string{"/labs-bigint"}, actions.BigRouter)
	starter.AppendRouter("GET", []string{"/labs-ytld"}, actions.YtIdRouter)
	starter.AppendRouter("GET", []string{"/labs-pem"}, actions.PemRouter)
	starter.AppendRouter("GET", []string{"/labs-leftpad"}, actions.PadLeftRouter)
	starter.AppendRouter("POST", []string{"/labs-redis-pub"}, actions.RedisPubRouter)
	starter.AppendRouter("POST", []string{"/labs-upload-file"}, actions.UploadFileRouterOne)
	starter.AppendRouter("GET", []string{"/labs-update-struct-pointer"}, actions.UpdateStructPointer)

	starter.AppendRouter("GET", []string{"/labs-redis-get"}, actions.RedisGetRouter)
	starter.AppendRouter("GET", []string{"/labs-redis-set"}, actions.RedisSetRouter)
	starter.AppendRouter("GET", []string{"/labs-redis-lpush"}, actions.RedisLpushRouter)
	starter.AppendRouter("GET", []string{"/labs-redis-del"}, actions.RedisDelRouter)
	starter.AppendRouter("GET", []string{"/labs-redis-hset"}, actions.RedisHsetRouter)
	starter.AppendRouter("GET", []string{"/labs-redis-hsetnx"}, actions.RedisDelRouter)

	starter.AppendRouter("POST", []string{"/websocket-client-simple"}, actions.WsClientSimple)
	starter.AppendRouter("POST", []string{"/demo/album-create"}, actions.AlbumCreateRouter)
	starter.AppendRouter("GET", []string{"/demo/album-get"}, actions.AlbumGetRouter)
	starter.AppendRouter("POST", []string{"/demo/binding-body"}, actions.BodyBindHandler)

	r := starter.Setup(starterServer)

	Redis_Conf := config.RedisConf

	if Redis_Conf.Mode != gredis.Disabled {
		for _, channel := range strings.Split(Redis_Conf.SubChannels, ",") {
			gredis.Sub(channel, func(channel string, payload string) {
				config.Log.Info("收到了消息1: channel=" + channel + ", payload=" + payload)
			})
		}
	}

	r.Run(config.AppSetting.AppPort)

}
