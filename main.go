package main

import (
	"fmt"
	"net/http"

	"github.com/chunhui2001/go-starter/controller"
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/starter"
	"github.com/gin-gonic/gin"
)

// go version devel +475d92ba4d Thu Oct 5 10:50:18 2017 +0000 darwin/amd64
var Author string

var (
	starterServer *starter.Server
	WEB_PAGE_CONF *config.WebPageConf = config.WebPageSettings
)

func init() {

	fmt.Printf("MainVar: %s\n", Author)

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
	r := starter.Bootstrap(starterServer)
	r.Run(config.AppSetting.AppPort)
}
