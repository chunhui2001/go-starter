package ghttp

import (
	"bytes"
	_ "context"
	"crypto/tls"
	"fmt"
	"io"
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
	Timeout             int  `mapstructure:"HTTP_CLIENT_TIMEOUT"`
	IdleConnTimeout     int  `mapstructure:"HTTP_CLIENT_IDLE_CONN_TIMEOUT"`
	MaxIdleConns        int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS"`
	MaxIdleConnsPerHost int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS_PERHOST"`
	MaxConnsPerHost     int  `mapstructure:"HTTP_CLIENT_MAX_CONNS_PERHOST"`
	PrintCurl           bool `mapstructure:"HTTP_CLIENT_PRINT_CURL"`
	PrintDebug          bool `mapstructure:"HTTP_CLIENT_PRINT_DEBUG"`
}

type HttpClient struct {
	Method      string
	Url         string
	TimeOut     int // 30 * time.Second
	QueryParams map[string]interface{}
	RequestBody string
	Headers     map[string]string
	Ellipsis    bool
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
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
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
	printCurl        bool
	printDebug       bool
)

func Init(conf *HttpConf, log *logrus.Entry) {

	logger = log
	DefaultTransport = defaultTransport(conf)
	printCurl = conf.PrintCurl
	printDebug = conf.PrintDebug

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

func DELETE(url string) *HttpClient {
	return &HttpClient{
		Method:  "DELETE",
		TimeOut: defaultTimeOut,
		Url:     url,
	}
}

func PUT(url string, reqBody string) *HttpClient {
	return &HttpClient{
		Method:      "PUT",
		TimeOut:     defaultTimeOut,
		Url:         url,
		RequestBody: reqBody,
		Ellipsis:    printDebug,
	}
}

func POST(url string, reqBody string) *HttpClient {
	return &HttpClient{
		Method:      "POST",
		TimeOut:     defaultTimeOut,
		Url:         url,
		RequestBody: reqBody,
		Ellipsis:    printDebug,
	}
}

func (c *HttpClient) SetHeaders(headers map[string]string) *HttpClient {
	c.Headers = headers
	return c
}

func (c *HttpClient) AddHeader(key string, val string) *HttpClient {
	c.Headers = utils.OfMap(key, val)
	return c
}

func (c *HttpClient) Query(queryParams map[string]interface{}) *HttpClient {

	if len(queryParams) == 0 {
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
	}

	for k, v := range httpClient.Headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Could-Not-Create-HttpRequest: Method=%s, contentType=%s, Url=%s, ErrorMessage=%s",
				httpClient.Method, httpClient.Headers["Content-Type"], httpClient.Url, err))

		return &HttpResult{
			Error: err,
		}
	}

	if len(httpClient.QueryParams) > 0 {

		q := req.URL.Query()

		for k, v := range httpClient.QueryParams {
			q.Add(k, utils.ToString(v))
		}

		req.URL.RawQuery = q.Encode()
	}

	command, _ := http2curl.GetCurlCommand(req)
	commandCurl := command.String()

	// req = req.WithContext(ctx)

	start := time.Now()
	res, err = myHttpClient.Do(req)
	latency := time.Since(start)

	if httpClient.RequestBody != "" {
		if !httpClient.Ellipsis && len(httpClient.RequestBody) > 165 {
			commandCurl = strings.Replace(command.String(), "-d ''", "-d '"+httpClient.RequestBody[0:160]+"....'", 1)
		} else {
			commandCurl = strings.Replace(command.String(), "-d ''", "-d '"+httpClient.RequestBody+"'", 1)
		}
	}

	if err != nil {
		if res != nil {
			logger.Error(
				fmt.Sprintf(
					"Send-HttpRequest-Failed: StatusCode=%d, Curl=%s, ErrorMessage=%s", res.StatusCode, commandCurl, err))
		} else {
			logger.Error(
				fmt.Sprintf(
					"Send-HttpRequest-Failed: Curl=%s, ErrorMessage=%s", commandCurl, err))
		}

		return &HttpResult{
			Error: err,
		}
	}

	contentLength := res.Header.Get("Content-Length")
	keepAlived := res.Header.Get("Connection")
	resBody, err := io.ReadAll(res.Body)
	contentLengthValue, err2 := strconv.Atoi(contentLength)

	if err2 != nil {
		contentLengthValue = -1
	}

	if err != nil {

		logger.Error(
			fmt.Sprintf(
				"HttpRequest-Could-Not-Read-Response-Body: StatusCode=%d, ContentLength=%s, Connection=%s, Curl=%s, ErrorMessage=%s",
				res.StatusCode, utils.HumanFileSizeWithInt(contentLengthValue), keepAlived, commandCurl, err))

		return &HttpResult{
			Error: err,
		}
	}

	if printCurl {
		logger.Info(
			fmt.Sprintf(
				"HttpRequest-Successful: Latency=%s, StatusCode=%d, ContentLength=%s, Connection=%s, Curl=%s",
				latency, res.StatusCode, utils.HumanFileSizeWithInt(contentLengthValue), keepAlived, commandCurl))
	}

	return &HttpResult{
		Status:       res.StatusCode,
		ResponseBody: resBody,
	}
}
