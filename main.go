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

func init() {

	starterServer = &starter.Server{
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

	// srv := &http.Server{
	// 	Addr:    config.AppSetting.AppPort,
	// 	Handler: r,
	// }

	// config.Log.Infof("listen: %s", "11")

	// if err := srv.ListenAndServe(); err != nil {
	// 	config.Log.Infof("listen: %s", "11")
	// 	config.Log.Infof("listen: %s\n", err)
	// } else {
	// 	config.Log.Infof("listen: %s", "11")
	// 	config.Log.Info("Congratulations! Your server startup successfully, Listening and serving HTTP on " + config.AppSetting.AppPort)

	// }

	r.Run(config.AppSetting.AppPort)
}
