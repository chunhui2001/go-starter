package controller

import (
	"net/http"

	"encoding/json"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/gid"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/chunhui2001/go-starter/pb/demo/wallet"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	WEB_PAGE_CONF *config.WebPageConf = config.WebPageSettings
)

func IndexRouter(c *gin.Context) {

	session := sessions.Default(c)

	if session.Get("yourRoomId") == nil {
		session.Set("yourRoomId", gid.ID())
		session.Save()
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"wssEndpoint": config.WssSetting.Wss(),
		"yourRoomId":  session.Get("yourRoomId"),
		"username":    session.Get("username"),
	})

}

func AboutRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "about", gin.H{
		"content": "This is an about page...",
	})
}

func TransactionRouter(c *gin.Context) {

	httpResult := ghttp.SendRequest(
		ghttp.GET("http://localhost:4002/scan-api/transaction/txns-list").Query(
			utils.MapOf("address", "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", "chainId", "1"),
		),
	)

	var transactionList []*wallet.Txns

	if httpResult.Success() {
		var m wallet.TxnsResult
		if err := json.Unmarshal(httpResult.ResponseBody, &m); err != nil {
			panic(err)
		} else {
			transactionList = m.Data
		}
	}

	c.HTML(http.StatusOK, "transactions/txns_index", gin.H{
		"transactionList": transactionList,
	})

}
