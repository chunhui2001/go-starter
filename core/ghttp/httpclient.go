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

type HttpConf struct {
	Timeout             int `mapstructure:"HTTP_CLIENT_TIMEOUT"`
	IdleConnTimeout     int `mapstructure:"HTTP_CLIENT_IDLE_CONN_TIMEOUT"`
	MaxIdleConns        int `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS"`
	MaxIdleConnsPerHost int `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS_PERHOST"`
	MaxConnsPerHost     int `mapstructure:"HTTP_CLIENT_MAX_CONNS_PERHOST"`
}

type HttpClient struct {
	Method      string
	Url         string
	TimeOut     int // 30 * time.Second
	QueryParams map[string]interface{}
	RequestBody string
	Headers     map[string]string
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

func defaultTransport(conf *HttpConf) http.RoundTripper {
	return &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(conf.Timeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(conf.Timeout) * time.Second,
		MaxIdleConns:        conf.MaxIdleConns,
		IdleConnTimeout:     time.Duration(conf.IdleConnTimeout) * time.Second,
		MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
		MaxConnsPerHost:     conf.MaxConnsPerHost,
		DisableCompression:  true,
		DisableKeepAlives:   false, // 默认选项
	}
}

var (
	logger           *logrus.Entry
	myHttpClient     *http.Client
	defaultTimeOut   int = 150
	DefaultTransport http.RoundTripper
)

func Init(conf *HttpConf, log *logrus.Entry) {

	logger = log
	DefaultTransport = defaultTransport(conf)

	myHttpClient = &http.Client{
		Transport: defaultTransport(conf),
		Timeout:   time.Duration(conf.Timeout) * time.Second,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
	}

	logger.Info("Initialization-a-HttpClient: " +
		"TimeOut=" + utils.ToString(conf.Timeout) + "s, " +
		"maxIdleConns=" + utils.ToString(conf.MaxIdleConns) + ", " +
		"idleConnTimeout=" + utils.ToString(conf.IdleConnTimeout) + "s, " +
		"maxIdleConnsPerHost=" + utils.ToString(conf.MaxIdleConnsPerHost) + ", " +
		"maxConnsPerHost=" + utils.ToString(conf.MaxConnsPerHost))

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

func (c *HttpClient) SetHeaders(headers map[string]string) *HttpClient {
	c.Headers = headers
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
	}

	for k, v := range httpClient.Headers {
		req.Header.Set(k, v)
	}

	// req.Header.Set("Accept-Encoding", "gzip")

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Could-Not-Create-HttpRequest: Method=%s, contentType=%s, Url=%s, ErrorMessage=%s",
				httpClient.Method, httpClient.Headers["Content-Type"], httpClient.Url, err))
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
					"Send-HttpRequest-Failed: StatusCode=%d, Curl=%s, ErrorMessage=%s", res.StatusCode, command, err))
		} else {
			logger.Error(
				fmt.Sprintf(
					"Send-HttpRequest-Failed: Curl=%s, ErrorMessage=%s", command, err))
		}
		return &HttpResult{
			Error: err,
		}
	}

	// keys := make([]string, 0, len(res.Header))

	// for k := range res.Header {
	// 	keys = append(keys, k)
	// }

	// logger.Info(
	// 	fmt.Sprintf(
	// 		"HttpRequest-Response-Headers: Keys=%s", utils.ToJsonString(keys)))

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
