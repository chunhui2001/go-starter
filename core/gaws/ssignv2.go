package gaws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
)

func SignV2Request(req *http.Request, accessKeyID string, secretAccessKey string) {
	req.URL.RawQuery = SignV2(accessKeyID, secretAccessKey, req.Method, req.URL.RequestURI(), nil)
}

func SignV2(accessKeyID string, secretAccessKey string, method string, rawUrl string, queryParams *map[string]interface{}) string {

	currurl, err := url.Parse(rawUrl)

	if err != nil {
		panic(err)
	}

	var Query url.Values = currurl.Query()

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

	host := currurl.Host
	path := currurl.Path

	if path == "" {
		path = "/"
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

	return Query.Encode()

}
