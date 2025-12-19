package s3client

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ruelsison-dev/emr-poc-01/backend/fhir-service-go/internal/aws"
)

func UploadDocument(ctx context.Context, key string, body []byte) (string, error) {
	c := aws.NewClients(ctx)
	bucket := os.Getenv("FHIR_S3_BUCKET")
	if bucket == "" {
		return "", fmt.Errorf("FHIR_S3_BUCKET not configured")
	}
	uploader := manager.NewUploader(c.S3)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(bytes.NewReader(body)),
		ServerSideEncryption: s3types.ServerSideEncryptionAwsKms,
	})
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("s3://%s/%s", bucket, key)
	return url, nil
}
