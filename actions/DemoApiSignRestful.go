package actions

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/gaws"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"
)

func AwsV2SignSimpleRouter(c *gin.Context) {

	accessKeyID := c.Query("accessKeyID")
	secretAccessKey := c.Query("secretAccessKey")

	c.Header("Content-Type", "text/plain")

	preSignedUrl, err := gaws.PreSignedUrlV2(accessKeyID, secretAccessKey, 10, c.Request.Method, utils.RequestURL(c.Request), nil)

	if err != nil {
		c.JSON(http.StatusOK, (&R{Error: err}).Fail(ILLEGAL_ACCESS))
		return
	}

	time.Sleep(1 * time.Second)

	result, err := gaws.CheckSign(accessKeyID, secretAccessKey, "POST", preSignedUrl)

	fmt.Println(err)
	fmt.Println(result)

	c.Writer.Write([]byte(preSignedUrl.String()))

}
