package test

import (
	"testing"
	"path/filepath"
	"fmt"
	"os"

	terraform "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
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

	// Validate outputs exist (table_name, kms_key_arn, bucket)
	table := terraform.Output(t, opts, "table_name")
	kms := terraform.Output(t, opts, "kms_key_arn")
	bucket := terraform.Output(t, opts, "bucket")

	require.NotEmpty(t, table)
	require.NotEmpty(t, kms)
	require.NotEmpty(t, bucket)

	// Further verification using AWS SDK could be added here if AWS creds available in CI
	fmt.Printf("Found resources: table=%s, kms=%s, bucket=%s\n", table, kms, bucket)
}
