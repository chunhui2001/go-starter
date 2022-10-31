package gaws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
)

const (
	signatureVersion = "2"
	signatureMethod  = "HmacSHA256"
	timeFormat       = "2006-01-02T15:04:05Z"
)

func SignV2Request(req *http.Request, accessKeyID string, secretAccessKey string) {

	newUrl, err1 := SignV2(accessKeyID, secretAccessKey, req.Method, req.URL, nil)

	if err1 != nil {
		panic(err1)
	}

	req.URL = newUrl

}

func CheckSign(accessKeyID string, secretAccessKey string, method string, reqUrl *url.URL) (bool, error) {

	newUrl, err1 := SignV2(accessKeyID, secretAccessKey, method, reqUrl, nil)

	if err1 != nil {
		return false, err1
	}

	var accessQuery url.Values = reqUrl.Query()
	var newQuery url.Values = newUrl.Query()

	// 签名不匹配, 签名无效
	if accessQuery.Get("Signature") != newQuery.Get("Signature") {
		return false, errors.New("UN_AUTH")
	}

	if accessQuery.Has("ExpireSeconds") {

		signTime, err := time.Parse(timeFormat, accessQuery.Get("Timestamp"))

		// 时间格式不对
		if err != nil {
			return false, errors.New("ILLEGAL_ACCESS")
		}

		expireSeconds, err := strconv.Atoi(accessQuery.Get("ExpireSeconds"))

		// 过期时间格式不对
		if err != nil {
			return false, errors.New("ILLEGAL_ACCESS")
		}

		// 过期时间在当前时间之后, 签名有效
		if signTime.Add(time.Duration(expireSeconds) * time.Second).After(time.Now()) {
			return true, nil
		}

		// 签名已过期
		return false, errors.New("ILLEGAL_ACCESS")

	}

	return true, nil

}

func SignV2(accessKeyID string, secretAccessKey string, method string, reqUrl *url.URL, queryParams *map[string]interface{}) (*url.URL, error) {

	var Query url.Values = reqUrl.Query()

	if queryParams != nil {
		for key, val := range *queryParams {
			Query.Set(key, utils.ToString(val))
		}
	}

	// Set new query parameters
	Query.Set("AWSAccessKeyId", accessKeyID)
	Query.Set("SignatureVersion", signatureVersion)
	Query.Set("SignatureMethod", signatureMethod)
	Query.Set("Timestamp", time.Now().Format(timeFormat))

	// in case this is a retry, ensure no signature present
	Query.Del("Signature")

	host := reqUrl.Host
	path := reqUrl.Path

	if path == "" {
		path = "/"
	} else if strings.Contains(path, "../") {
		return nil, errors.New("ILLEGAL_PARAMS")
	}

	// obtain all of the query keys and sort them
	queryKeys := make([]string, 0, len(Query))

	for key := range Query {
		queryKeys = append(queryKeys, key)
	}

	// sort keys
	sort.Strings(queryKeys)

	// build URL-encoded query keys and values
	queryKeysAndValues := make([]string, len(queryKeys))

	for i, key := range queryKeys {
		k := strings.Replace(url.QueryEscape(key), "+", "%20", -1)
		v := strings.Replace(url.QueryEscape(Query.Get(key)), "+", "%20", -1)
		queryKeysAndValues[i] = k + "=" + v
	}

	// join into one query string
	query := strings.Join(queryKeysAndValues, "&")

	// build the canonical string for the V2 signature
	stringToSign := strings.Join([]string{
		method,
		host,
		path,
		query,
	}, "\n")

	hash := hmac.New(sha256.New, []byte(secretAccessKey))
	hash.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	Query.Set("Signature", signature)

	newUrl, err1 := url.Parse(fmt.Sprintf(`%s://%s%s?%s`, reqUrl.Scheme, reqUrl.Host, reqUrl.Path, Query.Encode()))

	if err1 != nil {
		panic(err1)
	}

	return newUrl, nil

}
