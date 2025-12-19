package dynamo

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ruelsison-dev/emr-poc-01/backend/fhir-service-go/internal/aws"
)

type Client struct {
	cli *dynamodb.Client
	table string
}

func NewClient(ctx context.Context) *Client {
	c := aws.NewClients(ctx)
	t := os.Getenv("DYNAMO_TABLE_NAME")
	if t == "" {
		panic("DYNAMO_TABLE_NAME not set")
	}
	return &Client{cli: c.Dynamo, table: t}
}

func (c *Client) Put(ctx context.Context, item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	_, err = c.cli.PutItem(ctx, &dynamodb.PutItemInput{TableName: &c.table, Item: av})
	if err != nil {
		return fmt.Errorf("putitem: %w", err)
	}
	return nil
}

func (c *Client) Get(ctx context.Context, key map[string]types.AttributeValue, out interface{}) error {
	res, err := c.cli.GetItem(ctx, &dynamodb.GetItemInput{TableName: &c.table, Key: key})
	if err != nil {
		return err
	}
	if res.Item == nil {
		return fmt.Errorf("not_found")
	}
	if err := attributevalue.UnmarshalMap(res.Item, out); err != nil {
		return err
	}
	return nil
}
