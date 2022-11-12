package main

import (
	"net/http"

	"github.com/chunhui2001/go-starter/controller"
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/starter"
	"github.com/gin-gonic/gin"
)

var (
	starterServer *starter.Server
	WEB_PAGE_CONF *config.WebPageConf = config.WebPageSettings
)

func hander404(c *gin.Context) {
	if WEB_PAGE_CONF.Enable {
		c.HTML(http.StatusOK, "404", gin.H{
			"requestUrl": c.Request.URL.Path,
			"content":    "Page not found",
		})
	} else {
		c.JSON(http.StatusOK, (&R{}).Msg("Page-Not-Found").Fail(404))
	}
}

func init() {

	starterServer = &starter.Server{
		HandlerIndexPage: controller.IndexRouter,
		Handler404:       hander404,
	}

}

func main() {

	starterServer.Bootstrap(func(c *gin.Engine) {

	}).Running()
	// starterServer.Bootstrap().RunningTLS()

}
