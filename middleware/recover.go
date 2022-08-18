package middleware

import (
    "io"

    "github.com/go-errors/errors"

    "github.com/chunhui2001/go-starter/logger"
    "github.com/gin-gonic/gin"
)

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
    return RecoveryWithWriter(f, gin.DefaultErrorWriter)
}

func RecoveryWithWriter(f func(c *gin.Context, err interface{}), out io.Writer) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                //httprequest, _ := httputil.DumpRequest(c.Request, false)
                goErr := errors.Wrap(err, 3)
                logger.Log.Error("requestUri=", c.Request.RequestURI, ", errorMessage=", goErr.Error(), ", errorStack=", string(goErr.Stack()))

                f(c, err)
            }
        }()
        c.Next() // execute all the handlers
    }
}
