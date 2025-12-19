package integration

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	awsclient "github.com/ruelsison-dev/emr-poc-01/backend/fhir-service-go/internal/aws"
)

func randomSuffix() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Intn(100000))
}

func TestDynamoAndS3Integration(t *testing.T) {
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint == "" {
		t.Skip("AWS_ENDPOINT_URL not set; skipping LocalStack integration tests")
	}

	ctx := context.Background()
	clients := awsclient.NewClients(ctx)

	// DynamoDB table
	table := "test-fhir-" + randomSuffix()
	_, err := clients.Dynamo.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(table),
		AttributeDefinitions: []dtypes.AttributeDefinition{{AttributeName: aws.String("id"), AttributeType: dtypes.ScalarAttributeTypeS}},
		KeySchema:            []dtypes.KeySchemaElement{{AttributeName: aws.String("id"), KeyType: dtypes.KeyTypeHash}},
		BillingMode:          dtypes.BillingModePayPerRequest,
	})
	require.NoError(t, err)

	// Wait until table active
	var desc *dynamodb.DescribeTableOutput
	for i := 0; i < 10; i++ {
		desc, err = clients.Dynamo.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: &table})
		if err == nil && desc.Table.TableStatus == dtypes.TableStatusActive {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	require.NoError(t, err)
	require.Equal(t, dtypes.TableStatusActive, desc.Table.TableStatus)

	// Put an item
	item := map[string]interface{}{"id": "pat-test-1", "mrn": "MRN-123"}
	av, err := attributevalue.MarshalMap(item)
	require.NoError(t, err)
	_, err = clients.Dynamo.PutItem(ctx, &dynamodb.PutItemInput{TableName: &table, Item: av})
	require.NoError(t, err)

	// Get the item back
	res, err := clients.Dynamo.GetItem(ctx, &dynamodb.GetItemInput{TableName: &table, Key: map[string]dtypes.AttributeValue{"id": &dtypes.AttributeValueMemberS{Value: "pat-test-1"}}})
	require.NoError(t, err)
	require.NotNil(t, res.Item)

	// S3 bucket
	bucket := "test-fhir-bucket-" + randomSuffix()
	_, err = clients.S3.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &bucket})
	require.NoError(t, err)

	// Upload object
	key := "test/doc-1.txt"
	content := []byte("hello world")
	_, err = clients.S3.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: bytes.NewReader(content)})
	require.NoError(t, err)

	// Get object head
	_, err = clients.S3.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &bucket, Key: &key})
	require.NoError(t, err)

	// Cleanup (best-effort)
	_, _ = clients.S3.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &bucket, Key: &key})
	_, _ = clients.S3.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: &bucket})
	_, _ = clients.Dynamo.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: &table})
}
