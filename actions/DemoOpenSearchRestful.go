package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/goes"
	"github.com/chunhui2001/go-starter/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func OpenSearchIndicesRouter(c *gin.Context) {
	reault := goes.CatIndices("")
	c.JSON(200, (&R{Data: reault}).IfErr(400))
}

func NdJsonHandler(c *gin.Context) {

	indexName := c.Query("indexName")
	serverUri := c.Query("serverUri")
	var data = new([]map[string]interface{})

	if err := c.ShouldBindWith(data, binding.JSON); err != nil {
		c.JSON(200, (&R{Error: err}).Msg(err.Error()).IfErr(413))
		return
	}

	nsJsonString := goes.GetNdJson(indexName, "_doc", data)

	// c.Header("Content-Type", "application/octet-stream")
	// c.Writer.Write([]byte(nsJsonString))

	requestUrl := fmt.Sprintf(`%s/%s/_bulk?pretty=`, serverUri, indexName)

	httpResult := ghttp.SendRequest(
		ghttp.POST(requestUrl, nsJsonString).AddHeader("Content-Type", "application/x-ndjson"),
	)

	c.JSON(200, (&R{Data: httpResult.Success(), Error: httpResult.Error}).IfErr(400))

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

// {
//     "groupByFieldName": "bdUserId",
//     "aggs": [
//         {"inst": "cardinality", "aggsFieldName" : "fUid", "alias": "Accounts_Assigned"},
//         {"inst": "sum", "aggsFieldName": "netIncomeSubPoint", "alias": "Profit_QTD"},
//         {"inst": "sum", "aggsFieldName": "incomeSubPoint", "alias": "Revenue_QTD"},
//         {"inst": "sum", "aggsFieldName": "spotFeeSubPoint", "alias": "QTD_Revenue_Spot"},
//         {"inst": "sum", "aggsFieldName": "futureFee", "alias": "QTD_Revenue_Futures"},
//         {"inst": "sum", "aggsFieldName": "warrantFee", "alias": "QTD_Revenue_Options"},
//         {"inst": "sum", "aggsFieldName": "interestUsdt", "alias": "QTD_Margin_Interest_Revenue"},
//         {"inst": "sum", "aggsFieldName": "brokerUsdt", "alias": "Referral_Cost"},
//         {"inst": "sum", "aggsFieldName": "kolInviteCost", "alias": "Kol_Inviter_Cost"},
//         {"inst": "sum", "aggsFieldName": "kolRewardCost", "alias": "Kol_Reward_Cost"},
//         {"inst": "sum", "aggsFieldName": "allDeal", "alias": "Volume_QTD"},
//         {"inst": "sum", "aggsFieldName": "spotDeal", "alias": "QTD_Volume_Spot"},
//         {"inst": "sum", "aggsFieldName": "futureDeal", "alias": "QTD_Volume_Futures"}
//     ],
//     "where": {
//         "bdUsers": [147],
//         "labels": ["机构"],
//         "flags": [0,1]
//     }
// }
func OpenSearchMultipleAggsQueryRouter(c *gin.Context) {

	indexName := c.Query("indexName")
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

	groupByFieldName := params["groupByFieldName"].(string)

	result, _, err := goes.AggsQuery(indexName, "group_by_"+groupByFieldName, dslJsonString)

	if err != nil {
		c.JSON(200, (&R{Error: err}).Fail(400))
		return
	}

	c.JSON(200, (&R{Data: result}).Success())

}

func BulkQuery(c *gin.Context) {

	indexName := c.Query("indexName")
	payload, _ := ioutil.ReadAll(c.Request.Body)

	m := make([]string, 0)

	if err := json.Unmarshal(payload, &m); err != nil {
		fmt.Println(fmt.Sprintf("字节转json异常: jsonString=%s", string(payload)))
		panic(err)
	}

	result, err := goes.BulkQuery(indexName, &m)

	c.JSON(http.StatusOK, (&R{Data: result, Error: err}).IfErr(400))

}
