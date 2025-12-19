# Module Testing Guidance

This folder contains guidance for testing Terraform modules.

Recommended tools:
- Terratest (Go) for integration testing modules
- `terraform validate` + `terraform plan` smoke tests
- Conftest/OPA checks for policy validation

Suggested tests:
- Validate that S3 buckets created by the module enforce SSE-KMS and block public access
- Validate RDS modules include backup policies and correct parameter groups
- Validate network module produces private subnets and expected route tables

Add CI job to run Terratest for PRs that modify modules and to run conftest for policy checks on `terraform plan` outputs.
