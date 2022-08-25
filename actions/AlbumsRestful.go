package actions

import (
	"net/http"

	// "encoding/json"
	_ "github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/gmongo"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"

	// "github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	// "github.com/mitchellh/mapstructure"
)

type AlbumBook struct {
	Id     string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name" binding:"required"`
	Title  string `json:"title" bson:"title" binding:"required"`
	Author string `json:"author" bson:"author" binding:"required"`
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
