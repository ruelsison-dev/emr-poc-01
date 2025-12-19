package test

import (
	"testing"
	"path/filepath"
	"fmt"
	"os"
	"strings"

	terraform "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestFhirInfra_EncryptionAndPolicies(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..","..","infra","envs","dev")
	opts := &terraform.Options{
		TerraformDir: root,
	}

	defer terraform.Destroy(t, opts)
	tf, err := terraform.InitAndApplyE(t, opts)
	require.NoError(t, err, fmt.Sprintf("terraform apply failed: %v", err))

	// Validate outputs exist (table_name, kms_key_arn, bucket, role_arn)
	table := terraform.Output(t, opts, "table_name")
	kms := terraform.Output(t, opts, "kms_key_arn")
	bucket := terraform.Output(t, opts, "bucket")
	roleArn := terraform.Output(t, opts, "role_arn")
	policyArn := terraform.Output(t, opts, "policy_arn")

	require.NotEmpty(t, table)
	require.NotEmpty(t, kms)
	require.NotEmpty(t, bucket)
	require.NotEmpty(t, roleArn)
	require.NotEmpty(t, policyArn)

	// Use AWS SDK to validate resource properties (requires AWS creds in CI)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	require.NoError(t, err)

	// DynamoDB SSE
	dClient := dynamodb.NewFromConfig(cfg)
	desc, err := dClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(table)})
	require.NoError(t, err)
	require.NotNil(t, desc.Table.SSEDescription)
	require.Equal(t, dtypes.SSETypeKms, desc.Table.SSEDescription.SSEType)

	// S3 encryption
	s3Client := s3.NewFromConfig(cfg)
	encrypt, err := s3Client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{Bucket: aws.String(bucket)})
	require.NoError(t, err)
	require.NotEmpty(t, encrypt.ServerSideEncryptionConfiguration.Rules)
	require.Equal(t, s3types.ServerSideEncryptionAwsKms, encrypt.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)

	// KMS key
	kmsClient := kms.NewFromConfig(cfg)
	resp, err := kmsClient.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: aws.String(kms)})
	require.NoError(t, err)
	require.Equal(t, kmsTypes.KeyStateEnabled, resp.KeyMetadata.KeyState)

	// KMS key policy - ensure there is a default policy with statements
	kp, err := kmsClient.GetKeyPolicy(ctx, &kms.GetKeyPolicyInput{KeyId: aws.String(kms), PolicyName: aws.String("default")})
	require.NoError(t, err)
	require.NotEmpty(t, aws.ToString(kp.Policy))
	require.Contains(t, aws.ToString(kp.Policy), "Statement")

	// S3 encryption
	s3Client := s3.NewFromConfig(cfg)
	encrypt, err := s3Client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{Bucket: aws.String(bucket)})
	require.NoError(t, err)
	require.NotEmpty(t, encrypt.ServerSideEncryptionConfiguration.Rules)
	require.Equal(t, s3types.ServerSideEncryptionAwsKms, encrypt.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)

	// S3 Public Access Block - ensure public access is blocked
	pab, err := s3Client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{Bucket: aws.String(bucket)})
	require.NoError(t, err)
	require.True(t, pab.PublicAccessBlockConfiguration.BlockPublicAcls)
	require.True(t, pab.PublicAccessBlockConfiguration.BlockPublicPolicy)
	require.True(t, pab.PublicAccessBlockConfiguration.IgnorePublicAcls)
	require.True(t, pab.PublicAccessBlockConfiguration.RestrictPublicBuckets)

	// IAM role existence and policy attachment
	iamClient := iam.NewFromConfig(cfg)
	// extract role name from ARN
	parts := strings.Split(roleArn, "/")
	require.True(t, len(parts) >= 2)
	roleName := parts[len(parts)-1]
	_, err = iamClient.GetRole(ctx, &iam.GetRoleInput{RoleName: aws.String(roleName)})
	require.NoError(t, err)

	// check policy attached
	listRes, err := iamClient.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{RoleName: aws.String(roleName)})
	require.NoError(t, err)
	found := false
	for _, p := range listRes.AttachedPolicies {
		if *p.PolicyArn == policyArn {
			found = true
			break
		}
	}
	require.True(t, found, "expected policy to be attached to role")

	fmt.Printf("Found resources: table=%s, kms=%s, bucket=%s, role=%s\n", table, kms, bucket, roleArn)
}
