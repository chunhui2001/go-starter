package ghttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/config"
	_ "github.com/chunhui2001/go-starter/utils"
	"moul.io/http2curl"
)

type ContentType string

const (
	JSONBody ContentType = "application/json"
)

var (
	logger             = config.Log
	defaultTimeOut int = 5 // * time.Second
)

var DefaultTransport http.RoundTripper = &http.Transport{
	MaxIdleConns:        100,
	IdleConnTimeout:     90 * time.Second,
	DisableCompression:  true,
	MaxIdleConnsPerHost: 100,
	MaxConnsPerHost:     100,
}

type HttpClient struct {
	Method      string
	Url         string
	QueryParams map[string]interface{}
	RequestBody string
	ContentType ContentType // ghttp.JSONBody
	TimeOut     int         // 30 * time.Second
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

func (c *HttpClient) SetQueryParams(queryParams map[string]interface{}) *HttpClient {
	c.QueryParams = queryParams
	return c
}

type HttpResult struct {
	Status       int
	Message      string
	ResponseBody string
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

/*
	ghttp.Post(&ghttp.HttpClient{
		Method: http.MethodPost,
		Url: "https://www.google.com"
	})
*/
func SendRequest(httpClient *HttpClient) *HttpResult {

	client := &http.Client{
		Transport: DefaultTransport,
		Timeout:   time.Duration(httpClient.TimeOut) * time.Second,
	}

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

	res, err = client.Do(req)
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
			"HttpRequest-Successful: StatusCode=%d, Curl=%s", res.StatusCode, command))

	return &HttpResult{
		Status:       res.StatusCode,
		ResponseBody: string(resBody),
	}

}
