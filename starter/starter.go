package starter

import (
	"net/http"

	"github.com/chunhui2001/go-starter/config"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/logger"
	"github.com/chunhui2001/go-starter/middleware"

	"github.com/thinkerou/favicon"
)

func Setup() *gin.Engine {

	APP_PORT := config.GetEnv("APP_PORT", ":8080")

	// new engine
	engine := gin.New()

	// apply middleware
	engine.Use(gin.Recovery())
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessFormat())
	engine.Use(favicon.New("./static/favicon.ico")) // set favicon middleware

	// info router
	engine.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello world", "message": "Ok"})
	})

	logger.Log.Info("Listening and serving HTTP on " + APP_PORT)

	return engine

}
