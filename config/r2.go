package config

import (
	"os"
	"r2-gallery/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var R2 *s3.S3
var BucketName string

func InitR2() {
	BucketName = os.Getenv("R2_BUCKET_NAME")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("auto"),
		Endpoint:         aws.String(os.Getenv("R2_PUBLIC_URL")),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("R2_ACCESS_KEY_ID"),
			os.Getenv("R2_SECRET_ACCESS_KEY"),
			"",
		),
	})

	if err != nil {
		utils.LogError("无法初始化 R2", err)
	}

	R2 = s3.New(sess)
}
