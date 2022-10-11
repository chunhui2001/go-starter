package goes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/chunhui2001/go-starter/core/ges"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/opensearch-project/opensearch-go"
	"github.com/sirupsen/logrus"
)

var (
	esConf   *OpenESConf
	logger   *logrus.Entry
	esClient *opensearch.Client
)

type OpenESConf struct {
	Enable      bool   `mapstructure:"OPENES_ENABLE"`
	Servers     string `mapstructure:"OPENES_SERVERS"`
	DslFolder   string `mapstructure:"OPENES_DSL_TEMPLATE_FOLDER"`
	PrettyPrint bool   `mapstructure:"OPENES_PRETTY_PRINT"`
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

	ges.InitDSL(conf.DslFolder, conf.PrettyPrint, log)

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

func ConstructQuery(q string, size int) *strings.Reader {

	var queryJsonString = fmt.Sprintf(`{"query": { %s }, "size": %d}`, q, size)

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		fmt.Println("constructQuery() ERROR: query string not valid:", queryJsonString)
		fmt.Println("Using default match_all query")
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

// 查询所有索引
func CatIndices(indexNamePattern ...string) ([]map[string]interface{}, error) {

	res, err := esapi.CatIndicesRequest{Format: "json", FilterPath: indexNamePattern}.Do(context.Background(), esClient)

	if err != nil {
		logger.Errorf("OpenSearch-CatIndices-Error-1: ErrorMessage=%s", err.Error())
		return nil, err
	}

	defer res.Body.Close()
	var resMap []map[string]interface{}

	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &resMap); err != nil {
		logger.Errorf("OpenSearch-CatIndices-Error-2: ErrorMessage=%s, Indices=%s", err.Error(), string(body))
		return nil, err
	}

	return resMap, nil

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

func Search(indexName string, queryJsonString string) ([]map[string]interface{}, int64, error) {

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("OpenSearch-Search-Failed: ErrorMessage=%s, queryJsonString=%s", "Not a valid json query string", queryJsonString)
		return nil, 0, errors.New("Not a valid json query string")
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

func Collapse(indexName string, queryJsonString string) ([]map[string]interface{}, int64, error) {

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("OpenSearch-Collapse-Failed: ErrorMessage=%s, queryJsonString=%s", "Not a valid json query string", queryJsonString)
		return nil, 0, errors.New("Not a valid json query string")
	}

	// Pass the JSON query to the Golang client's Search() method
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(queryJsonString)),
		esClient.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		logger.Errorf("OpenSearch-Collapse-Error-1: queryJsonString=%s, ErrorMessage=%s", queryJsonString, err.Error())
		return nil, 0, err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("OpenSearch-Collapse-Error-2: ErrorMessage=%s", err.Error())
		return nil, 0, err
	}

	if resMap["error"] != nil {
		if resMap["error"].(map[string]interface{})["type"].(string) == "index_not_found_exception" {
			return nil, 0, nil
		}
		logger.Errorf("OpenSearch-Collapse-Error-3: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
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
			object := _map["fields"].(map[string]interface{})
			interfaceSlice = append(interfaceSlice, object)
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
		return nil, 0, errors.New("Not a valid json query string")
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
		logger.Errorf("OpenSearch-AggsQuery-Error-3: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
		return nil, 0, errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	if resMap["hits"] == nil {
		return nil, 0, nil
	}

	if resMap["_shards"] != nil {
		shardsMap := resMap["_shards"].(map[string]interface{})
		failuresArray := shardsMap["failures"].([]map[string]interface{})
		if len(failuresArray) > 0 {
			for _, failItem := range failuresArray {
				logger.Warnf(`OpenSearch-AggsQuery-Fail-3: indexName=%s, queryJsonString=%s, ErrorMessage=%s`, indexName, queryJsonString, utils.ToJsonString(failItem))
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
