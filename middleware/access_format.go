package middleware

import (
	"fmt"

	"github.com/chunhui2001/go-starter/utils"

	"github.com/gin-gonic/gin"
)

// CORS middleware from https://github.com/gin-gonic/gin/issues/29#issuecomment-89132826
func AccessFormat() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/favicon.ico", "/static"}, Formatter: func(param gin.LogFormatterParams) string {
			// your custom format
			return fmt.Sprintf("%s [%s] - Access %s \"%s %s %s %d %s\"\n",
				param.TimeStamp.Format(utils.TimeStampFormat),
				"INFO",
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				//param.Request.UserAgent(),
			)
		}})
}
