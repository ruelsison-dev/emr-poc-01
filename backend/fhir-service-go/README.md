# FHIR Service (Go + Gin) — scaffold

This service implements a FHIR-compatible REST API (scaffold) using Go and `gin-gonic/gin`.

Quickstart:

- Install Go 1.21
- cd backend/fhir-service-go
- go mod download
- go run ./main.go

Notes:
- Current scaffold uses an in-memory store for Patient endpoints. Phase 1 will integrate DynamoDB persistence and S3 for document exports (with SSE-KMS).
- AWS SDK v2 is used for DynamoDB/S3 clients. Ensure `AWS_REGION` and AWS credentials are configured in environment or role.
