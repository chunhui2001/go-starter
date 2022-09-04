package actions

import (
	"net/http"

	. "github.com/chunhui2001/go-starter/commons"
	_ "github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/gmongo"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

var (
// logger = config.Log
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	data := structs.Map(albumBook)
	gmongo.InsertOne("c_albums", &data)
	bsonBytes, err := bson.Marshal(data)

	if err != nil {
		panic(err)
	}

	bson.Unmarshal(bsonBytes, albumBook)

	c.JSON(http.StatusBadRequest, gin.H{
		"code":    200,
		"data":    albumBook,
		"message": "Ok",
	})

}

func AlbumGetRouter(c *gin.Context) {

	id := c.Query("id")

	var data map[string]interface{}
	gmongo.FindOne("c_albums", id, &data)

	c.JSON(http.StatusBadRequest, gin.H{
		"code":    200,
		"data":    data,
		"message": "Ok",
	})

}

func BodyBindHandler(c *gin.Context) {

	body := Body{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(500, R{Error: err}.Fail(400))
		return
	}

	c.JSON(http.StatusAccepted, R{Data: &body}.Success())

}
