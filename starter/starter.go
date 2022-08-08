package starter

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/wss"
	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/logger"
	"github.com/chunhui2001/go-starter/middleware"

	"github.com/thinkerou/favicon"
)

func Setup() *gin.Engine {

	APP_PORT := config.GetEnv("APP_PORT", ":8080")
	WEBSOCKET_ROUTER := config.GetEnv("WEBSOCKET_ROUTER", "")

	// new engine
	engine := gin.New()

	// init html template
	engine.SetFuncMap(template.FuncMap{"upper": strings.ToUpper})
	engine.LoadHTMLGlob("templates/*.html")

	// apply middleware
	engine.Use(gin.Recovery())
	engine.Use(static.Serve("/static", static.LocalFile("./static", false)))
	engine.Use(favicon.New("./static/favicon.ico")) // set favicon middleware
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessFormat())

	// info router
	engine.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello world", "message": "Ok"})
	})

	// index page
	engine.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an Home page.",
		})
	})

	// about page
	engine.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", gin.H{
			"content": "This is an about page...",
		})
	})

	if WEBSOCKET_ROUTER != "" {
		engine.GET(WEBSOCKET_ROUTER, wss.WebsocketUpgrade)
	}

	logger.Log.Info("Listening and serving HTTP on " + APP_PORT + ", websocket=" + WEBSOCKET_ROUTER)

	return engine

}
