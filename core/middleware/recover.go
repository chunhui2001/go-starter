package middleware

import (
	"io"
	"time"

	"github.com/go-errors/errors"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"
)

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return RecoveryWithWriter(f, gin.DefaultErrorWriter)
}

func RecoveryWithWriter(f func(c *gin.Context, err interface{}), out io.Writer) gin.HandlerFunc {
	start := time.Now()
	return func(c *gin.Context) {

		logLine := DefaultLogFormatter(LogParam(c, 500, c.Request.URL.Path, start))

		defer func() {
			if err := recover(); err != nil {
				// httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := errors.Wrap(err, 3)
				config.Log.Errorf(`%s ErrorStack=%s`, logLine, utils.ErrorToString(goErr))
				//f(c, err)
				AbortMsg(500, goErr, c) // Instead of c.AbortWithError(500, err)
				return
			}
		}()
		c.Next() // execute all the handlers
	}
}

func AbortMsg(code int, err error, c *gin.Context) {
	c.String(code, "Oops! Please retry.")
	if e := c.Error(err); e != nil {
		c.Abort()
		return
	}
	c.Abort()
}
