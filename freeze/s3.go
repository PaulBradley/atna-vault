package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (v *vault) s3UploadFile() {
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

	if v.gzValid {
		_, v.err = s3.New(s3session).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(*v.aws_s3_bucket),
			Key:                  aws.String(v.outputFilenamePrefix + v.outputFilename),
			ACL:                  aws.String("private"),
			Body:                 bytes.NewReader([]byte(v.gzBuffer.Bytes())),
			ContentLength:        aws.Int64(int64(binary.Size(v.gzBuffer.Bytes()))),
			ContentType:          aws.String(http.DetectContentType([]byte(v.gzBuffer.Bytes()))),
			ServerSideEncryption: aws.String("AES256"),
			StorageClass:         aws.String("GLACIER_IR"),
		})

		if v.err != nil {
			log.Println("WARN:" + v.err.Error())
		} else {
			log.Println("INFO:Uploaded " + v.outputFilenamePrefix + v.outputFilename)
		}
	}
}
