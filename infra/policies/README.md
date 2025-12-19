# Policy-as-Code (OPA) rules

This folder contains example OPA rules that should be executed as part of CI when running `terraform plan`.

Suggested integration:
- Run `conftest test` or `opa eval` against the Terraform plan JSON output.
- Fail the CI job if any rules produce `deny` results.

Example rules included:
- `disallow_public_s3.rego` — denies S3 buckets exposing public ACLs
- `enforce_kms.rego` — denies S3 buckets without SSE-KMS enabled

**Note**: Expand these rules to cover IAM, security group exposure, public IP creation, and resource tagging.
