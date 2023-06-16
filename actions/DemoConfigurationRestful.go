package actions

import (
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/config"

	"github.com/gin-gonic/gin"
)

func ReadCacheKey(c *gin.Context) {

	key := c.Query("key")
	var dataValue interface{}

	if err := config.ReadConfig(key, &dataValue); err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
	}

	c.JSON(200, (&R{Data: dataValue}).Success())

}
