package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ruelsison-dev/emr-poc-01/backend/fhir-service-go/internal/dynamo"
	"github.com/ruelsison-dev/emr-poc-01/backend/fhir-service-go/internal/s3"
)

// If DYNAMO_TABLE_NAME is set, use DynamoDB; otherwise fall back to in-memory store
var useDynamo = os.Getenv("DYNAMO_TABLE_NAME") != ""
var patients = map[string]map[string]interface{}{}
var ddbClient *dynamo.Client

func ensureDDB(ctx context.Context) {
	if ddbClient == nil {
		ddbClient = dynamo.NewClient(ctx)
	}
}

func createPatient(c *gin.Context) {
	ctx := c.Request.Context()
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_payload"})
		return
	}
	id := "pat_" + fmt.Sprint(time.Now().UnixNano())
	payload["id"] = id
	payload["created_at"] = time.Now().UTC().Format(time.RFC3339)

	if useDynamo {
		ensureDDB(ctx)
		if err := ddbClient.Put(ctx, payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "dynamo_error", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, payload)
		return
	}

	patients[id] = payload
	c.JSON(http.StatusCreated, payload)
}

func getPatient(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if useDynamo {
		ensureDDB(ctx)
		key := map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}}
		var out map[string]interface{}
		if err := ddbClient.Get(ctx, key, &out); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusOK, out)
		return
	}

	p, ok := patients[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
