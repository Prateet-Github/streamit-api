package s3

import (
	"context"
	"time"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func GenerateUploadURL(
	client *awss3.Client,
	bucket string,
	key string,
	contentType string,
) (string, error) {

	presignClient := awss3.NewPresignClient(client)

	req, err := presignClient.PresignPutObject(
		context.Background(),
		&awss3.PutObjectInput{
			Bucket:      &bucket,
			Key:         &key,
			ContentType: &contentType,
		},
		awss3.WithPresignExpires(1*time.Hour),
	)

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
