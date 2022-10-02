package actions

import (
	"net/http"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/ghttp"

	"github.com/gin-gonic/gin"
)

func GraphClientRouter(c *gin.Context) {

	// httpResult := ghttp.SendRequest(
	// 	ghttp.GET("http://localhost:4002/scan-api/transaction/txns-list").Query(
	// 		utils.MapOf("address", "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", "chainId", "1"),
	// 	),
	// )

	httpResult := ghttp.SendRequest(
		ghttp.GET("http://localhost:8080/index"),
	)

	if httpResult.Success() {
		c.JSON(http.StatusOK, (&R{Data: string(httpResult.ResponseBody)}).Success())
		return
	}

	c.JSON(http.StatusOK, (&R{Error: httpResult.Error}).Fail(400))

}
