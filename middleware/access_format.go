package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var timeStampFormat = "2006-01-02T15:04:05.000Z07:00"

// CORS middleware from https://github.com/gin-gonic/gin/issues/29#issuecomment-89132826
func AccessFormat(c *gin.Engine) {
	c.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("[%s] %s - Access %s \"%s %s %s %d %s\"\n",
			"INFO",
			param.TimeStamp.Format(timeStampFormat),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			//param.Request.UserAgent(),
		)
	}))
}
