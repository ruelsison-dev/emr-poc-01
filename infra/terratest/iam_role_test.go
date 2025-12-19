package test

import (
	"testing"
	"path/filepath"
	"fmt"
	"os"
	"time"
	"strconv"
	"context"
	"strings"

	terraform "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestTerrtestRunnerRole(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	root := filepath.Join("..","..","infra","iam","terrtest-runner")
	opts := &terraform.Options{
		TerraformDir: root,
		Vars: map[string]interface{}{
			"oidc_provider_arn": "arn:aws:iam::123456789012:oidc-provider/token.actions.githubusercontent.com",
			"github_repo":       "ruelsison-dev/emr-poc-01",
			"role_name":         fmt.Sprintf("terrtest-runner-%d", time.Now().Unix()),
		},
	}

	defer terraform.Destroy(t, opts)
	_, err := terraform.InitAndApplyE(t, opts)
	require.NoError(t, err, fmt.Sprintf("terraform apply failed: %v", err))

	roleArn := terraform.Output(t, opts, "role_arn")
	policyArn := terraform.Output(t, opts, "policy_arn")

	require.NotEmpty(t, roleArn)
	require.NotEmpty(t, policyArn)

	// Validate role exists and assume policy contains the OIDC provider string
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	require.NoError(t, err)

	iamClient := iam.NewFromConfig(cfg)

	// extract role name from ARN
	parts := strings.Split(roleArn, "/")
	require.True(t, len(parts) >= 2)
	roleName := parts[len(parts)-1]

	getRole, err := iamClient.GetRole(ctx, &iam.GetRoleInput{RoleName: aws.String(roleName)})
	require.NoError(t, err)

	require.Contains(t, aws.ToString(getRole.Role.AssumeRolePolicyDocument), "oidc")

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

	fmt.Printf("Created role: %s, policy: %s\n", roleArn, policyArn)
}
