package awsclient

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Clients struct {
	Dynamo *dynamodb.Client
	S3     *s3.Client
}

func NewClients(ctx context.Context) *Clients {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	d := dynamodb.NewFromConfig(cfg)
	s := s3.NewFromConfig(cfg)
	return &Clients{Dynamo: d, S3: s}
}
