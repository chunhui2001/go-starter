package actions

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/chunhui2001/go-starter/commons"
	"github.com/chunhui2001/go-starter/config"
	_ "github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/ghttp"
	"github.com/chunhui2001/go-starter/gras"
	"github.com/chunhui2001/go-starter/gredis"
	"github.com/chunhui2001/go-starter/utils"
	"github.com/gin-gonic/gin"
)

var (
	logger = config.Log
)

func BigRouter(c *gin.Context) {
	b := utils.BigIntRandom()
	c.JSON(http.StatusOK, R{Data: gin.H{
		"a": b,
		"b": utils.BigIntHexString(b),
		"c": utils.BigIntFromHexString(utils.BigIntHexString(b)),
		"d": b.String(),
		"e": utils.BigIntFromString(b.String()),
	}}.Success())
}

func YtIdRouter(c *gin.Context) {
	c.JSON(http.StatusOK, R{Data: utils.ShortId()}.Success())
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
	c.JSON(http.StatusOK, R{Data: utils.PadLeft("chui", "..", 3)}.Success())
}

func RedisPubRouter(c *gin.Context) {

	channel := c.Query("channel")
	payload, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		panic(err)
	}

	gredis.Pub(channel, string(payload))

	c.JSON(http.StatusOK, R{Data: true}.Success())

}

func HttpClientSimpleRouter(c *gin.Context) {
	httpResult := ghttp.SendRequest(ghttp.GET("https://www.google.com?fff=gg").Query(utils.MapOf("a", "b", "v", "你好")))
	c.JSON(http.StatusOK, R{Data: string(httpResult.ResponseBody)}.Success())
}

func UploadFileRouterOne(c *gin.Context) {

	// single file
	formFile, err := c.FormFile("file1")

	if err != nil {
		logger.Error("Upload-a-File-Error: errorMessage=" + err.Error())
		c.JSON(http.StatusOK, R{Error: err, Message: "Upload one file failed."}.Fail(400))
		return
	}

	openedFile, openerr := formFile.Open()

	if openerr != nil {
		logger.Error("Upload-File-Open-Error: errorMessage=" + openerr.Error())
		c.JSON(http.StatusOK, R{Error: openerr, Message: "Upload-File-Open-Error"}.Fail(400))
		return
	}

	uploadFileBytes, readerr := ioutil.ReadAll(openedFile)

	if readerr != nil {
		logger.Error("Upload-File-Read-Error: errorMessage=" + readerr.Error())
		c.JSON(http.StatusOK, R{Error: readerr, Message: "Upload-File-Read-Error."}.Fail(400))
		return
	}

	var fileBytes bytes.Buffer
	fileWriter := bufio.NewWriter(&fileBytes)

	zipWriter := zip.NewWriter(fileWriter)
	w1, ziperr := zipWriter.Create(formFile.Filename)

	if ziperr != nil {
		logger.Error("Upload-File-Create-ZipWriter-Error: errorMessage=" + ziperr.Error())
		c.JSON(http.StatusOK, R{Error: ziperr, Message: "Upload-File-Create-ZipWriter-Error."}.Fail(400))
		return
	}

	fileReader := bytes.NewReader(uploadFileBytes)

	if _, err := io.Copy(w1, fileReader); err != nil {
		logger.Error("Upload-File-Copy-ZipStream-Error: errorMessage=" + err.Error())
		c.JSON(http.StatusOK, R{Error: err, Message: "Upload-File-Copy-ZipStream-Error."}.Fail(400))
		return
	}

	zipfilename := filepath.Join(utils.TempDir(), strings.TrimSuffix(formFile.Filename, filepath.Ext(formFile.Filename))+".zip")

	if wrierr := os.WriteFile(zipfilename, fileBytes.Bytes(), 0644); wrierr != nil {
		logger.Error("Upload-File-Write-ZipFile-Error: errorMessage=" + err.Error())
		c.JSON(http.StatusOK, R{Error: wrierr, Message: "Upload-File-Write-ZipFile-Error."}.Fail(400))
		return
	}

	zipWriter.Close()

	logger.Info("Upload-a-File: FileName=" + formFile.Filename + ", Size=" + utils.ToString(formFile.Size))

	c.JSON(http.StatusOK, R{Data: zipfilename}.Success())

}

func UploadFileRouterMany(c *gin.Context) {

	// channel := c.Query("channel")
	// payload, err := ioutil.ReadAll(c.Request.Body)

	// if err != nil {
	// 	panic(err)
	// }

	c.JSON(http.StatusOK, R{Data: true}.Success())

}
