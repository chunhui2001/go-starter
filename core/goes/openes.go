package goes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchutil"
	"github.com/sirupsen/logrus"
)

var (
	esConf   *OpenESConf
	logger   *logrus.Entry
	esClient *opensearch.Client
	newLine  string = `
`
)

type OpenESConf struct {
	Enable       bool   `mapstructure:"OPENES_ENABLE"`
	Servers      string `mapstructure:"OPENES_SERVERS"`
	DslFolder    string `mapstructure:"OPENES_DSL_TEMPLATE_FOLDER"`
	PrettyPrint  bool   `mapstructure:"OPENES_PRETTY_PRINT"`
	DisablePrint bool   `mapstructure:"OPENES_PRETTY_DISABLE"`
}

func Init(conf *OpenESConf, log *logrus.Entry) {

	logger = log
	esConf = conf

	retryBackoff := backoff.NewExponentialBackOff()

	cfg := opensearch.Config{
		MaxRetries:    5,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		Transport: ghttp.DefaultTransport,
		Addresses: strings.Split(esConf.Servers, ","),
	}

	es, err := opensearch.NewClient(cfg)

	if err != nil {
		logger.Error(fmt.Sprintf("OpenSearch-Client-Connect-Failed: server=%s, errorMessage=%s", esConf.Servers, utils.ErrorToString(err)))
		return
	}

	esClient = Ping(es) // print server info

	ges.InitDSL(conf.DslFolder, conf.PrettyPrint, conf.DisablePrint, log)

}

func Ping(es *opensearch.Client) *opensearch.Client {

	res, err := es.Info()

	if err != nil {
		logger.Errorf("OpenSearch-Could-Not-Connected: server=%s, ErrorMessage=%s", esConf.Servers, err.Error())
		return nil
	}

	defer res.Body.Close()

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Error("OpenSearch-Error-Parsing-the-Response-Body: errorMessage={}", utils.ErrorToString(err))
	} else {

		clusterName := r["cluster_name"].(string)
		serverInfo := r["version"].(map[string]interface{})
		serverVersion := serverInfo["number"]
		luceneVersion := serverInfo["lucene_version"]

		logger.Info(fmt.Sprintf(
			"OpenSearch-Connected-Successful: Servers=%s, ClusterName=%s, ServerVersion=%s, LuceneVersion=%s, ClientDriverVersion=%s",
			esConf.Servers, clusterName, serverVersion, luceneVersion, opensearch.Version),
		)
	}

	return es

}

// 查询所有索引
// func CatIndices(indexNamePattern ...string) ([]map[string]interface{}, error) {

// 	res, err := esapi.CatIndicesRequest{Format: "json", FilterPath: indexNamePattern}.Do(context.Background(), esClient)

// 	if err != nil {
// 		logger.Errorf("OpenSearch-CatIndices-Error-1: ErrorMessage=%s", err.Error())
// 		return nil, err
// 	}

// 	defer res.Body.Close()
// 	var resMap []map[string]interface{}

// 	body, _ := io.ReadAll(res.Body)

// 	if err2 := json.Unmarshal(body, &resMap); err2 != nil {
// 		logger.Errorf("OpenSearch-CatIndices-Error-2: ErrorMessage=%s, Indices=%s", err.Error(), string(body))
// 		return nil, err
// 	}

// 	return resMap, nil

// }

// 查询索引是否存在
func CatIndices(indexName string) string {

	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/_cat/indices/%s`, serverUri, indexName)

	httpResult := ghttp.SendRequest(
		ghttp.GET(requestUrl),
	)

	if !httpResult.Success() {
		return ""
	}

	return string(httpResult.ResponseBody)

}

func CountApi(indexName string) (int64, error) {

	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/%s/_count`, serverUri, indexName)

	httpResult := ghttp.SendRequest(
		ghttp.GET(requestUrl),
	)

	if !httpResult.Success() {
		return 0, httpResult.Error
	}

	responseMap := utils.AsMap(httpResult.ResponseBody)

	return int64(responseMap["count"].(float64)), nil

}

// 查询索引是否存在
func IndexExists(indexName string) bool {

	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/%s`, serverUri, indexName)

	httpResult := ghttp.SendRequest(
		ghttp.GET(requestUrl),
	)

	if !httpResult.Success() {
		return false
	}

	responseMap := utils.AsMap(httpResult.ResponseBody)

	if responseMap["status"] != nil && responseMap["status"].(int) == 404 {
		return false
	}

	return responseMap[indexName] != nil

}

// 删除索引
func DeleteIndex(indexName string) bool {

	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/%s`, serverUri, indexName)

	httpResult := ghttp.SendRequest(
		ghttp.DELETE(requestUrl),
	)

	return httpResult.Status == 200 || httpResult.Status == 404

}

// 重命名索引
func RenameIndex(source string, dest string) bool {

	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/_reindex`, serverUri)

	httpResult := ghttp.SendRequest(
		ghttp.POST(requestUrl,
			fmt.Sprintf(`{ "source": { "index": "%s" }, "dest": { "index": "%s" } }`, source, dest)).AddHeader("Content-Type", "application/json"),
	)

	logger.Infof(`RenameIndex: source=%s, dest=%s, ResponseBody=%s`, source, dest, string(httpResult.ResponseBody))

	return httpResult.Status == 200

}

// 创建mapping
func PutMapping(mappingName string, jsonTemplate string) bool {

	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/_template/%s`, serverUri, mappingName)

	httpResult := ghttp.SendRequest(
		ghttp.PUT(requestUrl, jsonTemplate).AddHeader("Content-Type", "application/json"),
	)

	return httpResult.Success()

}

func ConstructQuery(q string, size int) *strings.Reader {

	var queryJsonString = fmt.Sprintf(`{"query": { %s }, "size": %d}`, q, size)

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("constructQuery() ERROR: query string not valid: %s", queryJsonString)
		logger.Errorf("Using default match_all query")
		queryJsonString = "{}"
	}

	// Build a new string from JSON query
	var b strings.Builder
	b.WriteString(queryJsonString)

	// Instantiate a *strings.Reader object from string
	read := strings.NewReader(b.String())

	// Return a *strings.Reader object
	return read

}

// 新增
func Save(indexName string, dataMap map[string]interface{}) (string, error) {
	return SaveOrUpdate(indexName, "", dataMap)
}

// 新增或更新
func SaveOrUpdate(indexName string, id string, dataMap map[string]interface{}) (string, error) {

	if dataMap == nil {
		return "", nil
	}

	if dataMap["@timestamp"] == nil {
		dataMap["@timestamp"] = utils.DateTimeUTCString()
	}

	_id := id

	if id == "" {
		_id = utils.Base64UUID()
	}

	// Instantiate a request object
	res, err := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: _id,
		Body:       strings.NewReader(utils.ToJsonString(dataMap)),
		Refresh:    "true",
	}.Do(context.Background(), esClient)

	if err != nil {
		logger.Errorf("OpenSearch-SaveOrUpdate-Error-1: ErrorMessage=%s", err.Error())
		return "", err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("OpenSearch-SaveOrUpdate-Error-1: ErrorMessage=%s", err.Error())
		return "", err
	}

	if resMap["error"] != nil {
		logger.Errorf("OpenSearch-SaveOrUpdate-Failed: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
		return "", errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	return _id, nil

}

// 批量处理
func Bulk(indexName string, dataMap *[]map[string]interface{}) (uint64, error) {

	bi, err := opensearchutil.NewBulkIndexer(opensearchutil.BulkIndexerConfig{
		Client:        esClient,
		NumWorkers:    4,
		FlushBytes:    1024 * 1024, // bytes
		FlushInterval: 1 * time.Second,
	})

	if err != nil {
		logger.Errorf("Es-Bulk-Error-1: ErrorMessage=%s", err.Error())
		return 0, err
	}

	res, err := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}.Do(context.Background(), esClient)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {

		res2, err := esClient.Indices.Create(indexName)

		if err != nil {
			panic(err)
		}

		defer res2.Body.Close()

	}

	start := time.Now().UTC()

	var countSuccessful uint64

	for _, item := range *dataMap {
		if err := bi.Add(context.Background(), getBulkIndexerItem(&item, &countSuccessful)); err != nil {
			panic(err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		panic(err)
	}

	dur := time.Since(start)

	if biStatus := bi.Stats(); biStatus.NumFailed > 0 {
		logger.Errorf(
			"Es-Bulk-Failed: IndexName=%s, Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			indexName,
			humanize.Comma(int64(biStatus.NumFlushed)),
			humanize.Comma(int64(biStatus.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStatus.NumFlushed))),
		)
		return 0, nil
	}

	logger.Infof("Es-Bulk-Successful: IndexName=%s, Count=%d, Duration=%s", indexName, countSuccessful, dur)

	return countSuccessful, nil

}

// // curl -X 'POST' -H 'Content-Type: application/x-ndjson' 'http://localhost:9092/index_name_here/_bulk?pretty' --data-binary "@/Users/keesh/Desktop/ndjson.txt"
func BulkRequest(indexName string, dataMap *[]map[string]interface{}) (bool, error) {

	nsJsonString := GetNdJson(indexName, "_doc", dataMap)
	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/%s/_bulk?pretty=`, serverUri, indexName)

	httpResult := ghttp.SendRequest(
		ghttp.POST(requestUrl, nsJsonString).AddHeader("Content-Type", "application/x-ndjson"),
	)

	// logger.Infof("Es-BulkRequest: IndexName=%s, Count=%d, ResponseBody=%s", indexName, len(*dataMap), string(httpResult.ResponseBody))

	if httpResult.Success() {
		return true, nil
	}

	return false, httpResult.Error

}

// curl -H "Content-Type: application/x-ndjson" -XGET http://127.0.0.1:9200/_msearch/template -d '
// {"index" : "index_name_here_*"}
// {"inline": {"size": 1, "query": {"bool":  {"must" : [{"term": {"type": "正向永续"}}] }}, "sort": [{"ts": {"order": "desc"}}]}}
// {"index" : "index_name_here_*"}
// {"inline": {"size": 1, "query": {"bool":  {"must" : [{"term": {"type": "正向永续"}}] }}, "sort": [{"ts": {"order": "desc"}}]}}
// '
func BulkQuery(indexName string, dslArray *[]string) ([]map[string]interface{}, error) {

	stringArray := make([]string, 0, len(*dslArray)*2)

	for _, dsl := range *dslArray {
		stringArray = append(stringArray, fmt.Sprintf(`{"index": "%s"}`, indexName))
		stringArray = append(stringArray, fmt.Sprintf(`{"inline": %s}`, dsl))
	}

	ndJsonString := strings.Join(stringArray, newLine) + newLine

	serverUri := strings.Split(esConf.Servers, ",")[0]
	requestUrl := fmt.Sprintf(`%s/_msearch/template`, serverUri)

	httpResult := ghttp.SendRequest(
		ghttp.POST(requestUrl, ndJsonString).AddHeader("Content-Type", "application/x-ndjson"),
	)

	if !httpResult.Success() {
		return nil, httpResult.Error
	}

	// Deserialize the response into a map.
	resMapResponses := utils.AsMap(httpResult.ResponseBody)["responses"].([]interface{})

	var interfaceSlice []map[string]interface{}

	for _, resMap1 := range resMapResponses {

		resMap := resMap1.(map[string]interface{})

		if resMap["error"] != nil {
			logger.Errorf("OpenSearch-BulkQuery-Error: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
			return nil, errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
		}

		if resMap["hits"] != nil {
			hitsMap := resMap["hits"].(map[string]interface{})
			var dataMap []interface{} = hitsMap["hits"].([]interface{})

			for _, item := range dataMap {

				_map := item.(map[string]interface{})
				id := _map["_id"].(string)
				object := _map["_source"].(map[string]interface{})
				object["id"] = id

				interfaceSlice = append(interfaceSlice, object)

			}
		}

	}

	return interfaceSlice, nil

}

func GetNdJson(indexName string, docType string, dataMap *[]map[string]interface{}) string {

	stringArray := make([]string, 0, len(*dataMap)*2)

	for _, item := range *dataMap {

		_id := item["_id"]

		if item["_id"] == nil {
			_id = utils.Base64UUID()
		}

		if item["@timestamp"] == nil {
			item["@timestamp"] = utils.DateTimeUTCString()
		}

		stringArray = append(stringArray, fmt.Sprintf(`{"index": {"_index": "%s", "_id": "%s", "_type" : "_doc"}}`, indexName, _id.(string)))
		stringArray = append(stringArray, utils.ToJsonString(item))

		item["_id"] = _id

	}

	return strings.Join(stringArray, newLine) + newLine

}

func getBulkIndexerItem(item *map[string]interface{}, countSuccessful *uint64) opensearchutil.BulkIndexerItem {

	if (*item)["_id"] == nil {
		(*item)["_id"] = utils.Base64UUID()
	}

	if (*item)["@timestamp"] == nil {
		(*item)["@timestamp"] = utils.DateTimeUTCString()
	}

	data, err := json.Marshal(item)

	if err != nil {
		panic(err)
	}

	// index, create, delete, update
	return opensearchutil.BulkIndexerItem{
		Action:     "create",
		DocumentID: (*item)["_id"].(string),
		Body:       bytes.NewReader(data),
		OnSuccess: func(ctx context.Context, item opensearchutil.BulkIndexerItem, res opensearchutil.BulkIndexerResponseItem) {
			atomic.AddUint64(countSuccessful, 1)
		},
		OnFailure: func(ctx context.Context, item opensearchutil.BulkIndexerItem, res opensearchutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				logger.Errorf("Es-Bulk-ERROR: ErrorMessage=%s", err.Error())
			} else {
				logger.Errorf("Es-Bulk-ERROR: ErrorType=%s, Reason=%s", res.Error.Type, res.Error.Reason)
			}
		},
	}
}

func Search(indexName string, queryJsonString string) ([]map[string]interface{}, int64, error) {

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("OpenSearch-Search-Failed: ErrorMessage=%s, queryJsonString=%s", "Not a valid json query string", queryJsonString)
		return nil, 0, errors.New("not a valid json query string")
	}

	// Pass the JSON query to the Golang client's Search() method
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(queryJsonString)),
		esClient.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		logger.Errorf("OpenSearch-Search-Error-1: queryJsonString=%s, ErrorMessage=%s", queryJsonString, err.Error())
		return nil, 0, err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("OpenSearch-Search-Error-2: ErrorMessage=%s", err.Error())
		return nil, 0, err
	}

	if resMap["error"] != nil {
		if resMap["error"].(map[string]interface{})["type"].(string) == "index_not_found_exception" {
			return nil, 0, nil
		}
		logger.Errorf("OpenSearch-Search-Error-3: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
		return nil, 0, errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	if resMap["hits"] == nil {
		return nil, 0, nil
	}

	hitsMap := resMap["hits"].(map[string]interface{})

	if hitsMap["hits"] == nil {
		return nil, 0, nil
	}

	var dataMap []interface{} = hitsMap["hits"].([]interface{})
	total := hitsMap["total"].(map[string]interface{})["value"].(float64)

	var interfaceSlice []map[string]interface{}

	if total > 0 {
		for _, item := range dataMap {

			_map := item.(map[string]interface{})
			id := _map["_id"].(string)
			object := _map["_source"].(map[string]interface{})
			object["id"] = id

			interfaceSlice = append(interfaceSlice, object)
		}
	}

	return interfaceSlice, int64(total), nil

}

// GET shakespeare/_search
// {
//   "size": 10,
//   "query": {
//     "match": {
//       "play_name": "Hamlet"
//     }
//   },
//   "search_after": [ 1, "32635"],
//   "sort": [
//     { "speech_number": "asc" },
//     { "_id": "asc" }
//   ]
// }
func SearchAfter(indexName string, queryJsonString string) ([]map[string]interface{}, int64, []interface{}, error) {

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("OpenSearch-Search-Failed: ErrorMessage=%s, queryJsonString=%s", "Not a valid json query string", queryJsonString)
		return nil, 0, nil, errors.New("not a valid json query string")
	}

	// Pass the JSON query to the Golang client's Search() method
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(queryJsonString)),
		esClient.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		logger.Errorf("OpenSearch-Search-Error-1: queryJsonString=%s, ErrorMessage=%s", queryJsonString, err.Error())
		return nil, 0, nil, err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("OpenSearch-Search-Error-2: ErrorMessage=%s", err.Error())
		return nil, 0, nil, err
	}

	if resMap["error"] != nil {
		if resMap["error"].(map[string]interface{})["type"].(string) == "index_not_found_exception" {
			return nil, 0, nil, nil
		}
		logger.Errorf("OpenSearch-Search-Error-3: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
		return nil, 0, nil, errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	if resMap["hits"] == nil {
		return nil, 0, nil, nil
	}

	hitsMap := resMap["hits"].(map[string]interface{})

	if hitsMap["hits"] == nil {
		return nil, 0, nil, nil
	}

	var dataMap []interface{} = hitsMap["hits"].([]interface{})
	total := int64(hitsMap["total"].(map[string]interface{})["value"].(float64))
	_count := len(dataMap)

	var interfaceSlice []map[string]interface{}
	var lastSort []interface{}

	if total > 0 {

		for index, item := range dataMap {

			_map := item.(map[string]interface{})
			id := _map["_id"].(string)
			object := _map["_source"].(map[string]interface{})
			object["id"] = id

			interfaceSlice = append(interfaceSlice, object)

			if index == _count-1 {
				lastSort = _map["sort"].([]interface{})
			}

		}
	}

	return interfaceSlice, total, lastSort, nil

}

func Collapse(indexName string, queryJsonString string) ([]map[string]interface{}, int64, error) {

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("OpenSearch-Collapse-Failed: IndexName=%s, ErrorMessage=%s, queryJsonString=%s", "Not a valid json query string", indexName, queryJsonString)
		return nil, 0, errors.New("not a valid json query string")
	}

	// Pass the JSON query to the Golang client's Search() method
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(queryJsonString)),
		esClient.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		logger.Errorf("OpenSearch-Collapse-Error-1: IndexName=%s, queryJsonString=%s, ErrorMessage=%s", indexName, queryJsonString, err.Error())
		return nil, 0, err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("OpenSearch-Collapse-Error-2: IndexName=%s, ErrorMessage=%s", indexName, err.Error())
		return nil, 0, err
	}

	if resMap["error"] != nil {
		if resMap["error"].(map[string]interface{})["type"].(string) == "index_not_found_exception" {
			return nil, 0, nil
		}
		logger.Errorf("OpenSearch-Collapse-Error-3: IndexName=%s, ErrorMessage=%s", indexName, utils.ToJsonString(resMap["error"]))
		return nil, 0, errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	if resMap["hits"] == nil {
		return nil, 0, nil
	}

	hitsMap := resMap["hits"].(map[string]interface{})

	if hitsMap["hits"] == nil {
		return nil, 0, nil
	}

	var dataMap []interface{} = hitsMap["hits"].([]interface{})
	total := hitsMap["total"].(map[string]interface{})["value"].(float64)

	var interfaceSlice []map[string]interface{}

	if total > 0 {
		for _, item := range dataMap {
			_map := item.(map[string]interface{})
			if _map["fields"] != nil {
				object := _map["fields"].(map[string]interface{})
				interfaceSlice = append(interfaceSlice, object)
			}
		}
	}

	return interfaceSlice, int64(total), nil

}

func AggsQuery(indexName string, aggsName string, queryJsonString string) ([]map[string]interface{}, int64, error) {

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("OpenSearch-AggsQuery-Failed: ErrorMessage=%s, queryJsonString=%s", "Not a valid json query string", queryJsonString)
		return nil, 0, errors.New("not a valid json query string")
	}

	// Pass the JSON query to the Golang client's Search() method
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(queryJsonString)),
		esClient.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		logger.Errorf("OpenSearch-AggsQuery-Error-1: queryJsonString=%s, ErrorMessage=%s", queryJsonString, err.Error())
		return nil, 0, err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("OpenSearch-AggsQuery-Error-2: ErrorMessage=%s", err.Error())
		return nil, 0, err
	}

	if resMap["error"] != nil {
		if resMap["error"].(map[string]interface{})["type"].(string) == "index_not_found_exception" {
			return nil, 0, nil
		}
		logger.Errorf("OpenSearch-AggsQuery-Error-3: indexName=%s, queryJsonString=%s, ErrorMessage=%s", indexName, queryJsonString, utils.ToJsonString(resMap["error"]))
		return nil, 0, errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	if resMap["hits"] == nil {
		return nil, 0, nil
	}

	if resMap["_shards"] != nil {
		shardsMap := resMap["_shards"].(map[string]interface{})
		if shardsMap["failures"] != nil {
			failuresArray := shardsMap["failures"].([]interface{})
			if len(failuresArray) > 0 {
				for _, failItem := range failuresArray {
					logger.Warnf(`OpenSearch-AggsQuery-Fail-4: indexName=%s, queryJsonString=%s, ErrorMessage=%s`, indexName, queryJsonString, utils.ToJsonString(failItem))
				}
			}
		}
	}

	hitsMap := resMap["hits"].(map[string]interface{})

	if hitsMap["hits"] == nil {
		return nil, 0, nil
	}

	total := hitsMap["total"].(map[string]interface{})["value"].(float64)

	var returnResult []map[string]interface{}

	if total > 0 {

		var aggsMap map[string]interface{} = resMap["aggregations"].(map[string]interface{})
		var aggsResult map[string]interface{} = aggsMap[aggsName].(map[string]interface{})
		var interfaceSlice = aggsResult["buckets"].([]interface{})

		mapSlice := make([]map[string]interface{}, len(interfaceSlice))

		for i, v := range interfaceSlice {
			mapSlice[i] = v.(map[string]interface{})
		}

		returnResult = mapSlice

	}

	return returnResult, int64(total), nil

}
