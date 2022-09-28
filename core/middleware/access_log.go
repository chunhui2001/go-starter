package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/core/config"
)

var defaultLogFormatter = func(param gin.LogFormatterParams) string {

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("Access %s \"%s %s %s %d %s\"",
		param.ClientIP,
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		//param.Request.UserAgent(),
	)
}

func init() {

}

func AccessLog() gin.HandlerFunc {
	return Print(gin.LoggerConfig{
		SkipPaths: []string{"/favicon.ico", "/static"},
	})
}

func Print(conf gin.LoggerConfig) gin.HandlerFunc {

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {

			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.Request.RemoteAddr
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			config.Log.Info(defaultLogFormatter(param))

		}
	}

}
