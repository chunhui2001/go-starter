package ghttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/utils"
	_ "github.com/chunhui2001/go-starter/utils"
	"github.com/sirupsen/logrus"
	"moul.io/http2curl"
)

type ContentType string

const (
	JSONBody ContentType = "application/json"
)

var DefaultTransport http.RoundTripper = &http.Transport{
	MaxIdleConns:        maxIdleConns,
	IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
	DisableCompression:  true,
	MaxIdleConnsPerHost: maxIdleConnsPerHost,
	MaxConnsPerHost:     maxConnsPerHost,
}

type HttpClient struct {
	Method      string
	Url         string
	QueryParams map[string]interface{}
	RequestBody string
	ContentType ContentType // ghttp.JSONBody
	TimeOut     int         // 30 * time.Second
}

type HttpResult struct {
	Status       int
	Message      string
	ResponseBody []byte
	Error        error
}

func (r *HttpResult) Success() bool {
	if r.Error != nil {
		return false
	}
	if r.Status < 200 || r.Status > 300 {
		return false
	}
	return true
}

var (
	logger              *logrus.Entry
	myHttpClient        *http.Client
	defaultTimeOut      int = 5 // * time.Second
	maxIdleConns        int = 100
	idleConnTimeout     int = 90
	maxIdleConnsPerHost int = 100
	maxConnsPerHost     int = 100
)

func Init(log *logrus.Entry) {

	logger = log

	myHttpClient = &http.Client{
		Transport: DefaultTransport,
		Timeout:   time.Duration(defaultTimeOut) * time.Second,
	}

	logger.Info("Initialization-a-HttpClient: " +
		"TimeOut=" + utils.ToString(defaultTimeOut) + "s, " +
		"maxIdleConns=" + utils.ToString(maxIdleConns) + ", " +
		"idleConnTimeout=" + utils.ToString(idleConnTimeout) + "s, " +
		"maxIdleConnsPerHost=" + utils.ToString(maxIdleConnsPerHost) + ", " +
		"maxConnsPerHost=" + utils.ToString(maxConnsPerHost))

}

func NEW(method string, url string) *HttpClient {
	return &HttpClient{
		Method:  method,
		TimeOut: defaultTimeOut,
		Url:     url,
	}
}

func GET(url string) *HttpClient {
	return &HttpClient{
		Method:  "GET",
		TimeOut: defaultTimeOut,
		Url:     url,
	}
}

func POST(url string, reqBody string) *HttpClient {
	return &HttpClient{
		Method:      "POST",
		TimeOut:     defaultTimeOut,
		Url:         url,
		RequestBody: reqBody,
	}
}

func (c *HttpClient) SetContentType(contentType ContentType) *HttpClient {
	c.ContentType = contentType
	return c
}

func (c *HttpClient) Query(queryParams map[string]interface{}) *HttpClient {

	if queryParams == nil || len(queryParams) == 0 {
		return c
	}

	if c.QueryParams == nil {
		c.QueryParams = queryParams
	} else {
		for k, v := range queryParams {
			c.QueryParams[k] = v
		}
	}

	return c
}

/*
	ghttp.Post(&ghttp.HttpClient{
		Method: http.MethodPost,
		Url: "https://www.google.com"
	})
*/
func SendRequest(httpClient *HttpClient) *HttpResult {

	var req *http.Request
	var res *http.Response
	var err error

	if strings.EqualFold(http.MethodGet, strings.TrimSpace(httpClient.Method)) {
		req, err = http.NewRequest(httpClient.Method, httpClient.Url, nil)
	} else {
		// postData := bytes.NewReader([]byte(httpClient.RequestBody))
		postData := bytes.NewBufferString(httpClient.RequestBody)
		req, err = http.NewRequest(httpClient.Method, httpClient.Url, postData)
		req.Header.Set("Content-Type", string(httpClient.ContentType))
	}

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Could-Not-Create-HttpRequest: Method=%s, contentType=%s, Url=%s, ErrorMessage=%s",
				httpClient.Method, httpClient.ContentType, httpClient.Url, err))
		return &HttpResult{
			Error: err,
		}
	}

	if httpClient.QueryParams != nil && len(httpClient.QueryParams) > 0 {

		q := req.URL.Query()

		for k, v := range httpClient.QueryParams {
			q.Add(k, utils.ToString(v))
		}

		req.URL.RawQuery = q.Encode()

	}

	start := time.Now()
	res, err = myHttpClient.Do(req)
	latency := time.Now().Sub(start)
	command, _ := http2curl.GetCurlCommand(req)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Send-HttpRequest-Failed: Curl=%s, ErrorMessage=%s", command, err))
		return &HttpResult{
			Error: err,
		}
	}

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"HttpRequest-Could-Not-Read-Response-Body: Curl=%s, ErrorMessage=%s", command, err))
		return &HttpResult{
			Error: err,
		}

	}

	logger.Info(
		fmt.Sprintf(
			"HttpRequest-Successful: Latency=%s, StatusCode=%d, Curl=%s", latency, res.StatusCode, command))

	return &HttpResult{
		Status:       res.StatusCode,
		ResponseBody: resBody,
	}

}
