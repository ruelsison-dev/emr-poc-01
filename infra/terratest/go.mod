module github.com/ruelsison-dev/emr-poc-01/infra/terratest

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.22.0
	github.com/aws/aws-sdk-go-v2/config v1.14.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.19.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.19.0
	github.com/aws/aws-sdk-go-v2/service/kms v1.8.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.8.0
	github.com/gruntwork-io/terratest/modules/terraform v0.44.0
	github.com/stretchr/testify v1.8.4
)
