package actions

import (
	"io/ioutil"
	"net/http"

	_ "github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/ghttp"
	"github.com/chunhui2001/go-starter/gras"
	"github.com/chunhui2001/go-starter/gredis"
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

	_, publicKey := gras.GenerateRSAKey(2048)

	data := utils.StringToBytes(publicKey)

	c.Header("Content-Type", "application/octet-stream")
	// Force browser download
	c.Header("Content-Disposition", "attachment; filename=public.pem")
	// Browser download or preview
	c.Header("Content-Disposition", "inline;filename=public.pem")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.Writer.Write(data)

}

func PadLeftRouter(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    utils.PadLeft("chui", "..", 3),
		"message": "Ok",
	})
}

func RedisPubRouter(c *gin.Context) {

	channel := c.Query("channel")
	payload, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		panic(err)
	}

	gredis.Pub(channel, string(payload))
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    true,
		"message": "Ok",
	})

}

func HttpClientSimpleRouter(c *gin.Context) {

	httpResult := ghttp.SendRequest(ghttp.GET("https://www.google.com"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    httpResult.ResponseBody,
		"message": "Ok",
	})

}
