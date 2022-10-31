package actions

import (
	"net/http"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/gaws"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"
)

func AwsV2SignSimpleRouter(c *gin.Context) {

	accessKeyID := c.Query("accessKeyID")
	secretAccessKey := c.Query("secretAccessKey")

	c.Header("Content-Type", "text/plain")

	newUrl, err := gaws.SignV2(accessKeyID, secretAccessKey, c.Request.Method, utils.RequestURL(c.Request), nil)

	if err != nil {
		c.JSON(http.StatusOK, (&R{Error: err}).Fail(ILLEGAL_ACCESS))
		return
	}

	c.Writer.Write([]byte(newUrl.String()))

}
