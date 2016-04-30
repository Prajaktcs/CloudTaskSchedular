package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

var s *s3.S3 = s3.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds))

func putItem(key string, bucket string) {
	file, _ := os.Open(key)
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	}
	_, err := s.PutObject(params)
	if err != nil {
		fmt.Println(err.Error())
	}
}
