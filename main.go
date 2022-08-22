package main

import (
	"net/http"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/controller"
	_ "github.com/chunhui2001/go-starter/gredis"
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
	r.Run(config.AppSetting.AppPort)
}
