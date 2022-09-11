package actions

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/gras"
	"github.com/chunhui2001/go-starter/core/gredis"
	"github.com/chunhui2001/go-starter/core/gwss"
	"github.com/chunhui2001/go-starter/core/utils"
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

func RedisDelRouter(c *gin.Context) {
	key := c.Query("key")
	gredis.Del(key)
	c.JSON(http.StatusOK, R{Data: gredis.Llen(key)}.Success())
}

func RedisGetRouter(c *gin.Context) {
	key := c.Query("key")
	c.JSON(http.StatusOK, R{Data: gredis.Get(key)}.Success())
}

func RedisSetRouter(c *gin.Context) {
	key := c.Query("key")
	val := c.Query("val")
	expir := c.Query("expir")
	gredis.Set(key, val, utils.StrToInt(expir))
	c.JSON(http.StatusOK, R{Data: utils.StrToInt(expir)}.Success())
}

func RedisLpushRouter(c *gin.Context) {
	key := c.Query("key")
	val := c.Query("val")
	gredis.Lpush(key, val)
	c.JSON(http.StatusOK, R{Data: gredis.Llen(key)}.Success())
}

func RedisLrangeRouter(c *gin.Context) {
	key := c.Query("key")
	gredis.Lrange(key, 0, -1)
	c.JSON(http.StatusOK, R{Data: gredis.Lrange(key, 0, -1)}.Success())
}

func RedisHsetRouter(c *gin.Context) {
	key := c.Query("key")
	f1 := c.Query("f1")
	v1 := c.Query("v1")
	gredis.Hset(key, f1, v1)
	result := []interface{}{
		gredis.Hgetall(key),
		gredis.Hget(key, f1),
	}
	c.JSON(http.StatusOK, R{Data: result}.Success())
}

func RedisIncrRouter(c *gin.Context) {
	key := c.Query("key")
	if result, err := gredis.Zincr(key); err == nil {
		c.JSON(http.StatusOK, R{Data: result}.Success())
		return
	}
	c.JSON(http.StatusOK, R{Data: false}.Success())
}

func RedisExpireRouter(c *gin.Context) {
	key := c.Query("key")
	gredis.Expire(key, 1)
	c.JSON(http.StatusOK, R{Data: true}.Success())
}

func RedisSetNxRouter(c *gin.Context) {
	key := c.Query("key")
	c.JSON(http.StatusOK, R{Data: gredis.SetNX(key, "asdfasd", 5)}.Success())
}

func RedisTtlRouter(c *gin.Context) {
	key := c.Query("key")
	if ttl, err := gredis.Ttl(key); err != nil {
		fmt.Println("1" + utils.ToString(ttl))
		c.JSON(http.StatusOK, R{Error: err}.Fail(400))
	} else {
		fmt.Println("2" + utils.ToString(ttl))
		c.JSON(http.StatusOK, R{Data: ttl}.Success())
	}
}

func HttpClientSimpleRouter(c *gin.Context) {

	httpResult := ghttp.SendRequest(ghttp.GET("https://www.google.com?fff=gg").Query(utils.MapOf("a", "b", "v", "你好")))

	if httpResult.Success() {
		c.JSON(http.StatusOK, R{Data: string(httpResult.ResponseBody)}.Success())
		return
	}

	c.JSON(http.StatusOK, R{Error: httpResult.Error}.Fail(400))

}

func HttpClientTimeOutRouter(c *gin.Context) {

	queryParams := utils.MapOf(
		"module", "account",
		"action", "txlist",
		"address", "0xf915a84bbcaff206d7745c448c0808d7c863731d",
		"startblock", "0",
		"sort", "asc",
	)

	httpResult := ghttp.SendRequest(ghttp.GET("https://api.etherscan.io/api").Query(queryParams))

	if httpResult.Success() {
		c.JSON(http.StatusOK, R{Data: string(httpResult.ResponseBody)}.Success())
		return
	}

	c.JSON(http.StatusOK, R{Error: httpResult.Error}.Fail(400))
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

func WsClientSimple(c *gin.Context) {

	connectId := utils.ShortId()
	serverAddress := c.Query("serverAddress")

	_, _, err := gwss.New(connectId, serverAddress).Connect(func(ctx context.Context, client *gwss.Client, messageBuf []byte) {
		message := utils.AsMap(messageBuf)
		if message != nil && message["topic"] != nil && message["topic"] == "server_ping" {
			msg := fmt.Sprintf(`{"message":"%s","action": "pong"}`, utils.DateTimeUTCString())
			client.WriteMessage(msg)
		} else {
			logger.Info(fmt.Sprintf(`WebSocket-Receive-a-Message: connectId=%s, message=%s`, client.ConnectId, string(messageBuf)))
		}
	})

	if err != nil {
		c.JSON(http.StatusOK, R{Error: err}.Fail(400))
		return
	}

	c.JSON(http.StatusOK, R{Data: connectId}.Msg("Connect Websocket successful").Success())

}

// // https://github.com/wamuir/go-xslt
// func XsltDemoRouter(c *gin.Context) {
// 	style := `<?xml version="1.0" encoding="UTF-8"?>
// <xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform" version="1.0">
//   <xsl:output method="xml" indent="yes"/>

//   <xsl:template match="/persons">
//     <root>
//       <xsl:apply-templates select="person"/>
//     </root>
//   </xsl:template>

//   <xsl:template match="person">
//   	<xsl:variable name="family">
//       <xsl:value-of select="family-name"/>
//     </xsl:variable>
//     <name username="{@username}" family-name26="{$family}">
//       <xsl:value-of select="name" />
//     </name>
//   </xsl:template>

// </xsl:stylesheet>`

// 	doc := `<?xml version="1.0" ?>
// <persons>
//   <person username="JS1">
//     <name>John春辉张</name>
//     <family-name>Smith</family-name>
//   </person>
//   <person username="MI1">
//     <name>Morka</name>
//     <family-name>Ismincius</family-name>
//   </person>
// </persons>`

// 	xs, err := xslt.NewStylesheet([]byte(style))

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer xs.Close()

// 	// doc is an XML document to be transformed and res is the result of
// 	// the XSL transformation, both as []byte.
// 	data, err := xs.Transform([]byte(doc))

// 	if err != nil {
// 		panic(err)
// 	}

// 	c.Header("Content-Type", "application/xml")
// 	// Browser download or preview
// 	c.Header("Content-Disposition", "inline;filename=public.pem")
// 	c.Header("Content-Transfer-Encoding", "binary")
// 	c.Header("Cache-Control", "no-cache")

// 	c.Writer.Write(data)

// }

type Client struct {
	ID        string
	LastPong  time.Time
	CreatedAt time.Time
}

type Subscription struct {
	Topic   string
	Clients *[]Client
}

// a server type to store all subscriptions
type Server struct {
	Subscriptions []Subscription
}

var s *Server = &Server{}

func UpdateStructPointer(c *gin.Context) {

	var client *Client = &Client{}
	var newClient *[]Client = &[]Client{*client}

	newTopic := &Subscription{
		Topic:   "topic1",
		Clients: newClient,
	}

	s.Subscriptions = append(s.Subscriptions, *newTopic)

	for i := range s.Subscriptions {
		var sub Subscription = s.Subscriptions[i]
		sub.Topic = "topic2"
		s.Subscriptions[i] = sub
		clients := *sub.Clients
		for j := range clients {
			clients[j].LastPong = time.Now()
		}
	}

	c.JSON(http.StatusOK, R{Data: s}.Success())

}
