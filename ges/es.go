package ges

import (
	"encoding/json"
	"fmt"
	"github.com/chunhui2001/go-starter/utils"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var (
	esConf   *ESConf
	logger   *logrus.Entry
	esClient *elasticsearch.Client
)

type ESConf struct {
	Enable  bool   `mapstructure:"ES_Enable"`
	Servers string `mapstructure:"ES_Servers"`
}

func Init(conf *ESConf, log *logrus.Entry) {

	logger = log
	esConf = conf

	cfg := elasticsearch.Config{
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
