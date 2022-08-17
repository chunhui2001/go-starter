package actions

import (
	"net/http"

	"github.com/chunhui2001/go-starter/gras"
	_ "github.com/chunhui2001/go-starter/logger"
	"github.com/chunhui2001/go-starter/utils"
	"github.com/gin-gonic/gin"
)

func BigRouter(c *gin.Context) {
	b := utils.BigIntRandom()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"a": b,
			"b": utils.BigIntHexString(b),
			"c": utils.BigIntFromHexString(utils.BigIntHexString(b)),
			"d": b.String(),
			"e": utils.BigIntFromString(b.String()),
		},
		"message": "Ok",
	})
}

func YtIdRouter(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    utils.ShortId(),
		"message": "Ok",
	})
}

func PemRouter(c *gin.Context) {
	privateKey, publicKey := gras.GenerateRSAKey(2048)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"privateKey": privateKey,
			"publicKey":  publicKey,
		},
		"message": "Ok",
	})
}
