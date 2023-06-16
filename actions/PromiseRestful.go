package actions

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/core/cache"
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/promise"
)

func WaitGroup(c *gin.Context) {

	cacheKey := "__WaitGroup"

	theResult := cache.CacheGet(cacheKey)

	if theResult == nil {

		// result := make(map[string]interface{})
		theMap := promise.NewConcurrencyMap()

		success := promise.WaitGroup(3, func() {
			// result["1"] = 1
			theMap.Put("1", 1)
		}, func() {
			// result["2"] = 1
			theMap.Put("2", 2)
			// time.Sleep(20 * time.Second)
		}, func() {
			// result["3"] = 1
			theMap.Put("3", 3)
		}, func() {
			// val := 0
			// result["4"] = 5 / val
			// theMap.Put("4", 5/val)
			theMap.Put("4", 4)
		})

		result := theMap.ToMap()

		if success {
			cache.PutCacheWithTTL(cacheKey, result, time.Duration(5)*time.Second)
			c.JSON(http.StatusOK, (&R{Data: result}).Ok())
			return
		}
	} else {
		logger.Infof(`命中缓存`)
		c.JSON(http.StatusOK, (&R{Data: theResult}).Ok())
	}

}
