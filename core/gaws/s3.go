package gaws

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/ghttp"
	"github.com/chunhui2001/go-starter/core/utils"
)

var (
	logger = config.Log
)

type CredentialMaps struct {
	Region     string `yaml:"region"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	BucketName string `yaml:"bucketName"`
}

func NewSession(accessKey string, secretKey string, regionName string) *session.Session {

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(regionName),
		MaxRetries:  aws.Int(3),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))

	return sess

}

// // 读取所有文件对象
// func ListObjects(sess *session.Session, bucket string, prefix string) [][]string {

// 	var maxKeys int64 = 10000000
// 	svc := s3.New(sess)

// 	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
// 		Bucket:  aws.String(bucket),
// 		MaxKeys: &maxKeys,
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	keys := [][]string{}

// 	for _, item := range resp.Contents {
// 		if strings.HasPrefix(*item.Key, prefix) {
// 			if *item.Size == 0 {
// 				keys = append(keys, []string{*item.Key, "0"})
// 			} else {
// 				keys = append(keys, []string{*item.Key, utils.HumanFileSize(float64(*item.Size))})
// 			}
// 		}
// 	}

// 	return keys

// }

// 读取所有文件对象
func ListObjects(sess *session.Session, bucket string, prefix string) [][]string {

	svc := s3.New(sess)
	keys := [][]string{}

	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		// Prefix: "",
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {

		for _, item := range p.Contents {
			if strings.HasPrefix(*item.Key, prefix) {
				if *item.Size == 0 {
					keys = append(keys, []string{*item.Key, "0"})
				} else {
					keys = append(keys, []string{*item.Key, utils.HumanFileSize(float64(*item.Size))})
				}
			}
		}

		if last {
			return false
		}

		return true

	})

	if err != nil {
		panic(err)
	}

	return keys

}

func GetS3RequestUrl(sess *session.Session, bucket string, key string) string {

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		panic(err)
	}

	return urlStr

}

func PutObject(sess *session.Session, bucket string, key string, dataBuf []byte) (bool, error) {

	svc := s3.New(sess)
	contentLength := int64(len(dataBuf))

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),                      // Required
		Key:           aws.String(key),                         // 文件名
		ACL:           aws.String("bucket-owner-full-control"), // Optional
		Body:          bytes.NewReader(dataBuf),
		ContentLength: aws.Int64(contentLength),
	}

	_, err := svc.PutObject(params)

	if err != nil {
		logger.Errorf(`PutObject-Failed: FilePath=%s, ErrorMessage=%s`, key, utils.ErrorToString(err))
		return false, err
	}

	logger.Infof(`PutObject-Success: FilePath=%s, Size=%s`, key, utils.HumanFileSizeWithInt(len(dataBuf)))

	return true, nil

}

func DownloadFile(sess *session.Session, bucket string, key string) (*[]byte, error) {

	fileUrl := GetS3RequestUrl(sess, bucket, key)

	httpResult := ghttp.SendRequest(
		ghttp.GET(fileUrl),
	)

	if !httpResult.Success() {
		return nil, httpResult.Error
	}

	return &httpResult.ResponseBody, nil

}

func AlertDing() {

	dingTalkHookToken := "asdfadf;ljql243hrj3brk2j34rklj234kl2h34l5kjh234l"
	DINGTALK_HOST_PREFIX := "https://oapi.dingtalk.com/robot/send"

	msgTemplate := `\n
&#x1F695; &#x1F695; &#x1F695; <font color=\"#dd0000\">title_here[ding]</font>\n
> 名字1: **%s** \n
`
	_msg := fmt.Sprintf(
		msgTemplate,
		"name1_here",
	)

	postMessage := fmt.Sprintf(`{"msgtype": "markdown", "markdown": { "title": "title_here[ding]", "text": "%s" }}`, _msg)

	httpResult := ghttp.SendRequest(
		ghttp.POST(DINGTALK_HOST_PREFIX, postMessage).AddHeader("Content-Type", "application/json").Query(utils.MapOf("access_token", dingTalkHookToken)),
	)

	logger.Infof(`DingTalk-Send-Success: Message=%d`, httpResult.Status)

}
