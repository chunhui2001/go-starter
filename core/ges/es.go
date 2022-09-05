package ges

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

var (
	esConf   *ESConf
	logger   *logrus.Entry
	esClient *elasticsearch.Client
	ctx      context.Context
)

type ESConf struct {
	Enable  bool   `mapstructure:"ES_ENABLE"`
	Servers string `mapstructure:"ES_SERVERS"`
}

func Init(conf *ESConf, log *logrus.Entry) {

	logger = log
	esConf = conf
	ctx = context.Background()

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
		logger.Error(fmt.Sprintf("ElasticSearch-Client-Connect-Failed: %s, errorMessage=%s", esConf.Servers, utils.ErrorToString(err)))
		return
	}

	esClient = Ping(es) // print server info

}

func Ping(es *elasticsearch.Client) *elasticsearch.Client {

	res, err := es.Info()

	if err != nil {
		logger.Error("Error getting response: " + utils.ErrorToString(err))
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

func Search() {

}

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
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: _id,
		Body:       strings.NewReader(utils.ToJsonString(dataMap)),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, esClient)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	// Deserialize the response into a map.
	var resMap map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		return "", err
	}

	if resMap["error"] != nil {
		logger.Errorf("Es-SaveOrUpdate-Save-Failed: ErrorMessage=%s", utils.ToJsonString(resMap["error"]))
		return "", errors.New(resMap["error"].(map[string]interface{})["reason"].(string))
	}

	return _id, nil

}
