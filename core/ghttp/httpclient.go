package ghttp

import (
	"bytes"
	_ "context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/sirupsen/logrus"
	"moul.io/http2curl"
)

type ContentType string

const (
	JSONBody ContentType = "application/json"
)

var DefaultTransport http.RoundTripper = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: time.Duration(defaultTimeOut) * time.Second,
	}).Dial,
	TLSHandshakeTimeout: time.Duration(defaultTimeOut) * time.Second,
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
	defaultTimeOut      int = 150 // * time.Second
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
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 避免 ioutil.ReadAll(res.Body) 超时
	// defer cancel()

	// go func() {
	// 	select {
	// 	case <-time.After(10 * time.Second):
	// 		logger.Error(
	// 			fmt.Sprintf(
	// 				"Send-HttpRequest-Context-TimeOut-After: Url=%s, Message=%s", httpClient.Url, ctx.Err()))
	// 	case <-ctx.Done():
	// 		logger.Infof(
	// 			fmt.Sprintf(
	// 				"Send-HttpRequest-Context-TimeOut-Done: Url=%s, Message=%s", httpClient.Url, ctx.Err()))
	// 	}
	// }()

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

	// req = req.WithContext(ctx)

	start := time.Now()
	res, err = myHttpClient.Do(req)
	latency := time.Since(start)
	command, _ := http2curl.GetCurlCommand(req)

	if err != nil {
		if res != nil {
			logger.Error(
				fmt.Sprintf(
					"Send-HttpRequest-Failed: Curl=%s, ErrorMessage=%s", command, err))
		} else {
			logger.Error(
				fmt.Sprintf(
					"Send-HttpRequest-Failed: StatusCode=%d, Curl=%s, ErrorMessage=%s", res.StatusCode, command, err))
		}
		return &HttpResult{
			Error: err,
		}
	}

	keys := make([]string, 0, len(res.Header))

	for k := range res.Header {
		keys = append(keys, k)
	}

	contentLength := res.Header.Get("Content-Length")
	keepAlived := res.Header.Get("Connection")
	resBody, err := ioutil.ReadAll(res.Body)
	contentLengthValue, err2 := strconv.Atoi(contentLength)

	if err2 != nil {
		contentLengthValue = -1
	}

	if err != nil {

		logger.Error(
			fmt.Sprintf(
				"HttpRequest-Could-Not-Read-Response-Body: StatusCode=%d, ContentLength=%s, Connection=%s, Curl=%s, ErrorMessage=%s",
				res.StatusCode, utils.HumanFileSizeWithInt(contentLengthValue), keepAlived, command, err))
		return &HttpResult{
			Error: err,
		}

	}

	logger.Info(
		fmt.Sprintf(
			"HttpRequest-Successful: Latency=%s, StatusCode=%d, ContentLength=%s, Connection=%s, Curl=%s",
			latency, res.StatusCode, utils.HumanFileSizeWithInt(contentLengthValue), keepAlived, command))

	return &HttpResult{
		Status:       res.StatusCode,
		ResponseBody: resBody,
	}

}
