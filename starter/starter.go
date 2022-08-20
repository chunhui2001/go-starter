package starter

import (
	"html/template"
	"net/http"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/wss"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/static"
	"github.com/go-errors/errors"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/logger"
	"github.com/chunhui2001/go-starter/middleware"

	"github.com/chunhui2001/go-starter/actions"
	"github.com/chunhui2001/go-starter/controller"
	"github.com/thinkerou/favicon"
)

func Setup() *gin.Engine {

	APP_PORT := config.AppSetting.AppPort
	WSS_PREFIX := config.WssSetting.Prefix

	// new engine
	engine := gin.New()

	// init html template
	engine.HTMLRender = ginview.New(goview.Config{
		Root:      "views",
		Extension: ".html",
		Master:    "layouts/master",
		//Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
			// more funcs
		},
		DisableCache: true,
	})

	// apply middleware
	engine.Use(middleware.Recovery(recoveryHandler)) // error nice handle
	engine.Use(static.Serve("/static", static.LocalFile("./static", false)))
	engine.Use(favicon.New("./static/favicon.ico")) // set favicon middleware
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessFormat())

	// info router
	engine.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello world", "message": "Ok"})
	})

	// index page
	engine.GET("", controller.IndexRouter)
	engine.GET("/index", controller.IndexRouter)
	engine.GET("/home", controller.IndexRouter)

	// about page
	engine.GET("/about", controller.AboutRouter)

	// labs
	engine.GET("/labs-bigint", actions.BigRouter)
	engine.GET("/labs-ytld", actions.YtIdRouter)
	engine.GET("/labs-pem", actions.PemRouter)

	if WSS_PREFIX != "" {
		engine.GET(WSS_PREFIX, wss.WebsocketUpgrade)
	}

	engine.NoRoute(func(c *gin.Context) {
		if c.Request.RequestURI == "/favicon.ico" {
			c.Next()
		} else {
			logger.Log.Warn(c.Request.RequestURI)

			c.HTML(http.StatusNotFound, "404", gin.H{
				"content": "Page not found",
			})
		}
	})

	logger.Log.Info("Listening and serving HTTP on " + APP_PORT + ", websocket=" + config.WssSetting.Wss())

	return engine

}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": errors.Wrap(err, 3).Error()})
}
