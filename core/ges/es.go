package ges

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/sirupsen/logrus"
)

// ## Elasticsearch Knowledge Base
// # How To Insert Elasticsearch Documents Into An Index Using Golang
// https://kb.objectrocket.com/elasticsearch/how-to-insert-elasticsearch-documents-into-an-index-using-golang-451
// # How To Construct Elasticsearch Queries From A String Using Golang
// https://kb.objectrocket.com/elasticsearch/how-to-construct-elasticsearch-queries-from-a-string-using-golang-550

var (
	esConf   *ESConf
	logger   *logrus.Entry
	esClient *elasticsearch.Client
)

type ESConf struct {
	Enable       bool   `mapstructure:"ES_ENABLE"`
	Servers      string `mapstructure:"ES_SERVERS"`
	DslFolder    string `mapstructure:"ES_DSL_TEMPLATE_FOLDER"`
	PrettyPrint  bool   `mapstructure:"OPENES_PRETTY_PRINT"`
	DisablePrint bool   `mapstructure:"OPENES_PRETTY_DISABLE"`
}

func Init(conf *ESConf, log *logrus.Entry) {

	logger = log
	esConf = conf

	retryBackoff := backoff.NewExponentialBackOff()

	cfg := elasticsearch.Config{
		MaxRetries:    5,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 5 * time.Second,
		},
		Addresses: strings.Split(esConf.Servers, ","),
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		logger.Error(fmt.Sprintf("ElasticSearch-Client-Connect-Failed: server=%s, errorMessage=%s", esConf.Servers, utils.ErrorToString(err)))
		return
	}

	esClient = Ping(es) // print server info

	InitDSL(conf.DslFolder, conf.PrettyPrint, conf.DisablePrint, log)

}

func Ping(es *elasticsearch.Client) *elasticsearch.Client {

	res, err := es.Info()

	if err != nil {
		logger.Errorf("ElasticSearch-Could-Not-Connected: server=%s, ErrorMessage=%s", esConf.Servers, err.Error())
		return nil
	}

	defer res.Body.Close()

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Error("ElasticSearch-Error-Parsing-the-Response-Body: errorMessage={}", utils.ErrorToString(err))
	} else {
		clusterName := r["cluster_name"].(string)
		serverInfo := r["version"].(map[string]interface{})
		serverVersion := serverInfo["number"]
		luceneVersion := serverInfo["lucene_version"]

		logger.Info(fmt.Sprintf(
			"Elastic-Client-Connected-Successful: Servers=%s, ClusterName=%s, ServerVersion=%s, LuceneVersion=%s, ClientDriverVersion=%s",
			esConf.Servers, clusterName, serverVersion, luceneVersion, elasticsearch.Version),
		)
	}

	return es

}

// 查询所有索引
func CatIndices() ([]map[string]interface{}, error) {

	res, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), esClient)

	if err != nil {
		logger.Errorf("Es-CatIndices-Error-1: ErrorMessage=%s", err.Error())
		return nil, err
	}

	defer res.Body.Close()
	var resMap []map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("Es-CatIndices-Error-2: ErrorMessage=%s", err.Error())
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
		logger.Errorf("Es-SaveOrUpdate-Error-1: ErrorMessage=%s", err.Error())
		return "", err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("Es-SaveOrUpdate-Error-1: ErrorMessage=%s", err.Error())
		return "", err
	}

	if resMap["error"] != nil {
		logger.Errorf("Es-SaveOrUpdate-Failed: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
		return "", errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	return _id, nil

}

// 批量处理
func Bulk(indexName string, dataMap *[]map[string]interface{}) (bool, error) {

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        esClient,
		NumWorkers:    4,
		FlushBytes:    1024 * 1024, // bytes
		FlushInterval: 1 * time.Second,
	})

	if err != nil {
		logger.Errorf("Es-Bulk-Error-1: ErrorMessage=%s", err.Error())
		return false, err
	}

	var countSuccessful uint64

	for _, item := range *dataMap {

		err = bi.Add(context.Background(), getBulkIndexerItem(&item, &countSuccessful))

		if err != nil {
			panic(err)
		}

	}

	if err := bi.Close(context.Background()); err != nil {
		panic(err)
	}

	biStatus := bi.Stats()

	if biStatus.NumFailed > 0 {
		return false, nil
	}

	return true, nil

}

func getBulkIndexerItem(item *map[string]interface{}, countSuccessful *uint64) esutil.BulkIndexerItem {

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

	return esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: (*item)["_id"].(string),
		Body:       bytes.NewReader(data),
		OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
			atomic.AddUint64(countSuccessful, 1)
		},
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				logger.Errorf("Es-Bulk-ERROR: ErrorMessage=%s", err.Error())
			} else {
				logger.Errorf("Es-Bulk-ERROR: ErrorType=%s, ErrorMessage=%s", res.Error.Type, res.Error.Reason)
			}
		},
	}
}

func Search(indexName string, queryJsonString string) ([]map[string]interface{}, int64, error) {

	// Check for JSON errors
	isValid := json.Valid([]byte(queryJsonString)) // returns bool

	// Default query is "{}" if JSON is invalid
	if !isValid {
		logger.Errorf("Es-Search-Failed: ErrorMessage=%s, queryJsonString=%s", "Not a valid json query string", queryJsonString)
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
		logger.Errorf("Es-Search-Error-1: queryJsonString=%s, ErrorMessage=%s", queryJsonString, err.Error())
		return nil, 0, err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		logger.Errorf("Es-Search-Error-2: ErrorMessage=%s", err.Error())
		return nil, 0, err
	}

	if resMap["error"] != nil {
		if resMap["error"].(map[string]interface{})["type"].(string) == "index_not_found_exception" {
			return nil, 0, nil
		}
		logger.Errorf("Es-Search-Error-3: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
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
			logger.Infof(`Id=%s, len=%d`, id, len(object))
			object["id"] = id

			interfaceSlice = append(interfaceSlice, object)
		}
	}

	return interfaceSlice, int64(total), nil

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
