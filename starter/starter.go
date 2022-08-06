package starter

import (
	"net/http"

	"go-starter/config"

	"github.com/gin-gonic/gin"

	"go-starter/logger"
	"go-starter/middleware"
)

func Setup() *gin.Engine {

	APP_PORT := config.GetEnv("APP_PORT", ":8080")

	// new engine
	engine := gin.New()

	// apply middleware
	middleware.AccessFormat(engine)
	engine.Use(gin.Recovery())
	engine.Use(middleware.CORS(middleware.CORSOptions{}))

	// info router
	engine.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello world", "message": "Ok"})
	})

	logger.Log.Info("Listening and serving HTTP on " + APP_PORT)

	return engine

}
