package aws

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
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	var cfg aws.Config
	var err error
	if endpoint != "" {
		// Use a custom endpoint resolver for localstack/testing
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{PartitionID: "aws", URL: endpoint, SigningRegion: region}, nil
			})),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	d := dynamodb.NewFromConfig(cfg)
	s := s3.NewFromConfig(cfg)
	return &Clients{Dynamo: d, S3: s}
}
