package gaws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewSession() *session.Session {

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		panic(err)
	}

	return sess

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

func GetPresignedURL(sess *session.Session, bucket string, key string) string {

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
