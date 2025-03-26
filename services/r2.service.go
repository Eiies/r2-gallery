package services

import (
	"fmt"
	"io"
	"os"
	"r2-gallery/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToR2(file io.Reader, fileName string) (string, error) {
	_, err := config.R2.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(fileName),
		Body:   aws.ReadSeekCloser(file),
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", os.Getenv("YOU"), fileName)
	return url, nil
}

func DeleteFromR2(fileName string) error {
	_, err := config.R2.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(fileName),
	})

	return err
}

func GetR2ObjectURL(fileName string) string {
	return fmt.Sprintf("%s/%s/%s", os.Getenv("R2_PUBLIC_URL"), config.BucketName, fileName)
}
