package actions

import (
	"strconv"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/goes"
	"github.com/chunhui2001/go-starter/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	DSL_FILE_NAME = "dsl2.yaml"
)

func init() {

}

func OpenSearchIndicesRouter(c *gin.Context) {
	reault, err := goes.CatIndices()
	c.JSON(200, (&R{Data: reault, Error: err}).IfErr(400))
}

func OpenSearchLastSnapshotDateRouter(c *gin.Context) {

	indexName := c.Query("indexName")

	dslJsonString, err := ges.DSLQuery(DSL_FILE_NAME, "LAST_SNAPSHOT_DATE", nil)

	reault, _, err := goes.Search(indexName, dslJsonString)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	if len(reault) > 0 {
		if date, err := strconv.Atoi(reault[0]["createDate"].(string)); err == nil {
			c.JSON(200, (&R{Data: date}).Success())
			return
		}
	}

	c.JSON(200, (&R{Data: nil}).Success())

}

func OpenSearchQueryDataByDateRouter(c *gin.Context) {

	indexName := c.Query("indexName")
	querySize := c.Query("size")
	snapshotDate := c.Query("createDate")

	var data = new(map[string]interface{})

	if err := c.ShouldBindWith(data, binding.JSON); err != nil {
		c.JSON(200, (&R{Error: err}).Msg(err.Error()).IfErr(413))
		return
	}

	params := utils.MapOf(
		"size", utils.StrToInt(querySize),
		"createDate", utils.StrToInt(snapshotDate),
		"bdUsers", (*data)["bdUsers"],
		"labels", (*data)["labels"],
		"flags", (*data)["flags"],
	)

	dslJsonString, err := ges.DSLQuery(DSL_FILE_NAME, "QUERY_DATA_BY_DATE", params)

	reault, _, err := goes.Search(indexName, dslJsonString)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(200, (&R{Data: reault}).Success())

}

func OpenSearchDynamicQueryRouter(c *gin.Context) {

	indexName := c.Query("indexName")
	querySize := c.Query("size")
	snapshotDate := c.Query("createDate")
	name := c.Query("tplname")

	var dynamicParams = new(map[string]interface{})

	if err := c.ShouldBindWith(dynamicParams, binding.JSON); err != nil {
		c.JSON(200, (&R{Error: err}).Msg(err.Error()).IfErr(413))
		return
	}

	var params map[string]interface{} = utils.MapOf(
		"size", utils.StrToInt(querySize),
		"createDate", utils.StrToInt(snapshotDate),
	)

	for key, val := range *dynamicParams {
		params[key] = val
	}

	dslJsonString, err := ges.DSLQuery(DSL_FILE_NAME, name, params)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	reault, _, err := goes.Search(indexName, dslJsonString)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(200, (&R{Data: reault}).Success())

}

func OpenSearchDistinctQueryRouter(c *gin.Context) {

	indexName := c.Query("indexName")
	querySize := c.Query("size")
	snapshotDate := c.Query("createDate")
	fieldName := c.Query("fieldName")
	tplname := "DISTINCT_QUERY"

	var dynamicParams = new(map[string]interface{})

	if err := c.ShouldBindWith(dynamicParams, binding.JSON); err != nil {
		c.JSON(200, (&R{Error: err}).Msg(err.Error()).IfErr(413))
		return
	}

	var params map[string]interface{} = utils.MapOf(
		"size", utils.StrToInt(querySize),
		"createDate", utils.StrToInt(snapshotDate),
		"fieldName", fieldName,
	)

	for key, val := range *dynamicParams {
		params[key] = val
	}

	dslJsonString, err := ges.DSLQuery(DSL_FILE_NAME, tplname, params)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	result, _, err := goes.Collapse(indexName, dslJsonString)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	var currentResult []map[string]interface{} = result
	var returnResult = make([]any, 0)

	for i, _ := range currentResult {
		var currItem map[string]interface{} = currentResult[i]
		var currentArray []any = currItem[fieldName].([]any)
		for j, _ := range currentArray {
			returnResult = append(returnResult, currentArray[j])
		}
	}

	c.JSON(200, (&R{Data: returnResult}).Success())

}

func OpenSearchAggsSumQueryRouter(c *gin.Context) {

	indexName := c.Query("indexName")
	// querySize := c.Query("size") 0
	snapshotDate := c.Query("createDate")
	byFieldName := c.Query("byFieldName")
	sumFieldName := c.Query("sumFieldName")
	tplname := "AGGS_SUM"

	var dynamicParams = new(map[string]interface{})

	if err := c.ShouldBindWith(dynamicParams, binding.JSON); err != nil {
		c.JSON(200, (&R{Error: err}).Msg(err.Error()).IfErr(413))
		return
	}

	var params map[string]interface{} = utils.MapOf(
		"createDate", utils.StrToInt(snapshotDate),
		"byFieldName", byFieldName,
		"sumFieldName", sumFieldName,
	)

	for key, val := range *dynamicParams {
		params[key] = val
	}

	dslJsonString, err := ges.DSLQuery(DSL_FILE_NAME, tplname, params)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	result, _, err := goes.AggsSum(indexName, "by_"+byFieldName, dslJsonString)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(200, (&R{Data: result}).Success())

}

func OpenSearchMultipleAggsQueryRouter(c *gin.Context) {

	// indexName := c.Query("indexName")
	// querySize := c.Query("size") // default: 10000
	snapshotDate := c.Query("createDate")
	tplname := c.Query("tplname")

	var dynamicParams = new(map[string]interface{})

	if err := c.ShouldBindWith(dynamicParams, binding.JSON); err != nil {
		c.JSON(200, (&R{Error: err}).Msg(err.Error()).IfErr(413))
		return
	}

	var params map[string]interface{} = utils.MapOf(
		"createDate", utils.StrToInt(snapshotDate),
	)

	for key, val := range *dynamicParams {
		params[key] = val
	}

	dslJsonString, err := ges.DSLQuery(DSL_FILE_NAME, tplname, params)

	if err != nil {
		c.JSON(200, (&R{Data: dslJsonString, Error: err}).Fail(400))
		return
	}

	// result, _, err := goes.AggsSum(indexName, "group_by_"+groupByFieldName, dslJsonString)

	// if err != nil {
	// 	c.JSON(200, (&R{Error: err}).Fail(400))
	// 	return
	// }

	c.JSON(200, (&R{Data: dslJsonString}).Success())

}
