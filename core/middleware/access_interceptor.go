package middleware

import (
	"net/url"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/gaws"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

type AccessClients struct {
	ClientName      string   `yaml:"clientName"`
	AccessKeyID     string   `yaml:"accessKeyID"`
	SecretAccessKey string   `yaml:"secretAccessKey"`
	Enabled         bool     `yaml:"enabled"`
	Scope           []string `yaml:"scope"`
}

var (
	accessClientsMap map[string]*AccessClients
)

func init() {
	if err := config.ReadConfig("Access-Clients", &accessClientsMap); err != nil {
		logger.Errorf("Access-Clients-Configuration-Error: Key=%s, ErrorMessage=%s", "Access-Clients", err.Error())
	}
}

func AccessInterceptor(enable bool) gin.HandlerFunc {

	return func(c *gin.Context) {

		if !enable || accessClientsMap == nil || len(accessClientsMap) == 0 {
			c.Next()
			return
		}

		requestUrl := utils.RequestURL(c.Request)
		var accessQuery url.Values = requestUrl.Query()

		if !accessQuery.Has(gaws.AWSAccessKeyIdFieldKey) {
			AbortAccess(errors.New("UN_AUTH_MISS_ACCESS_KEY_ID"), c)
			c.Next()
			return
		}

		accessKeyId := accessQuery.Get(gaws.AWSAccessKeyIdFieldKey)

		var accessClient *AccessClients = accessClientsMap[accessKeyId]

		if accessClient == nil {
			AbortAccess(errors.New("UN_AUTH_INVALID_ACCESS_KEY_ID"), c)
			c.Next()
			return
		}

		if accessClient.Enabled {

			if _, err := gaws.CheckSign(accessKeyId, accessClient.SecretAccessKey, c.Request.Method, requestUrl); err != nil {
				AbortAccess(err, c)
				c.Next()
				return
			}

		}

		c.Next()

	}

}

func AbortAccess(err error, c *gin.Context) {
	c.String(401, err.Error())
	if e := c.Error(err); e != nil {
		c.Abort()
		return
	}
	c.Abort()
}
