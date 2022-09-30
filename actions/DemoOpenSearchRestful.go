package actions

import (
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/goes"
	"github.com/chunhui2001/go-starter/core/utils"

	"github.com/gin-gonic/gin"
)

// var (
// 	logger = config.Log
// )

func init() {

}

func OpenSearchIndicesRouter(c *gin.Context) {
	reault, err := goes.CatIndices()
	c.JSON(200, (&R{Data: reault, Error: err}).IfErr(400))
}

func OpenSearchRouter(c *gin.Context) {

	indexName := c.Query("indexName")

	params := utils.MapOf("a", "b")

	dslJsonString, err := ges.DSLQuery("dsl2.yaml", "QUERY_ALL", params)

	logger.Infof(`dslJsonString=%s`, dslJsonString)

	reault, _, err := goes.Search(indexName, dslJsonString)

	c.JSON(200, (&R{Data: reault, Error: err}).IfErr(400))

}
