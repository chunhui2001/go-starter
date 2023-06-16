package actions

import (
	"io/ioutil"
	"net/http"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/gmongo"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

type AlbumBook struct {
	Id     string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name" binding:"required"`
	Title  string `json:"title" bson:"title" binding:"required"`
	Author string `json:"author" bson:"author" binding:"required"`
}

type Body struct {
	Product string `json:"product" binding:"required,alpha"`
	Price   uint   `json:"price" binding:"required,gte=10,lte=1000"`
}

func AlbumCreateRouter(c *gin.Context) {

	var albumBook = &AlbumBook{}

	if err := c.ShouldBindJSON(&albumBook); err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	data := structs.Map(albumBook)
	gmongo.InsertOne("c_albums", &data)
	bsonBytes, err := bson.Marshal(data)

	if err != nil {
		panic(err)
	}

	bson.Unmarshal(bsonBytes, albumBook)
	c.JSON(http.StatusAccepted, (&R{Data: albumBook}).Success())

}

func AlbumGetRouter(c *gin.Context) {

	id := c.Query("id")
	var data map[string]interface{}
	gmongo.FindOne("c_albums", id, &data)
	c.JSON(200, (&R{Data: data}).Success())
}

func BodyBindHandler(c *gin.Context) {

	body := Body{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(500, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(http.StatusAccepted, (&R{Data: &body}).Success())

}

func ElsCatIndicesRouter(c *gin.Context) {

	indices, err := ges.CatIndices()

	if err != nil {
		c.JSON(http.StatusOK, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(http.StatusOK, (&R{Data: indices}).Success())

}

func ElsSearcherRouter(c *gin.Context) {

	indexName := c.Query("indexName")
	payload, _ := ioutil.ReadAll(c.Request.Body)

	reaults, total, err := ges.Search(indexName, string(payload))

	if err != nil {
		c.JSON(http.StatusOK, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(http.StatusOK, (&R{Data: reaults}).Msg("Total-Count: "+utils.ToString(total)).Success())

}

func ElsSearcherAdvanceRouter(c *gin.Context) {

	indexName := c.Query("indexName")
	payload, _ := ioutil.ReadAll(c.Request.Body)

	reaults, total, err := ges.Search(indexName, string(payload))

	if err != nil {
		c.JSON(http.StatusOK, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(http.StatusOK, (&R{Data: reaults}).Msg("Total-Count: "+utils.ToString(total)).Success())

}

func ElsDslTemplateRouter(c *gin.Context) {

	tplfile := c.Query("tplfile")
	tplname := c.Query("tplname")

	payload, _ := ioutil.ReadAll(c.Request.Body)

	dslJsonString, _ := ges.DSLQuery(tplfile, tplname, utils.AsMap(payload))
	c.JSON(http.StatusOK, (&R{Data: dslJsonString}).Success())

}

func ElsCreateOrUpdateRouter(c *gin.Context) {

	var albumBook = &AlbumBook{}

	if err := c.ShouldBindJSON(&albumBook); err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	_id, err := ges.SaveOrUpdate("go-simple-index", albumBook.Id, utils.ToMap(albumBook))

	if err != nil {
		c.JSON(http.StatusOK, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(http.StatusOK, (&R{Data: _id}).Success())

}
