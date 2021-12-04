package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (v *vault) getMessageObjectKeysFromS3() {
	var s3session *session.Session

	s3session, v.err = session.NewSession(&aws.Config{
		Region: aws.String(*v.aws_region),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("S3_KEY"),
			os.Getenv("S3_SECRET"),
			""),
	})

	if v.err != nil {
		log.Println("WARN:" + v.err.Error())
	}

	var continuationToken *string
	i := 0
	for {
		params := &s3.ListObjectsV2Input{
			Bucket: aws.String(*v.aws_s3_bucket),
			//Prefix:            aws.String("01FP0"),
			Prefix:            aws.String("AWSLogs/747808398410/elasticloadbalancing/eu-west-2/2021/"),
			ContinuationToken: continuationToken,
		}

		resp, err := s3.New(s3session).ListObjectsV2(params)

		if err != nil {
			log.Println(err.(awserr.Error).Code(), err.(awserr.Error).Error())
		}

		for _, obj := range resp.Contents {
			fmt.Println(*obj.Key)
			i = i + 1
		}

		if !aws.BoolValue(resp.IsTruncated) {
			break
		}
		continuationToken = resp.NextContinuationToken
	}

	fmt.Println(i)
}
