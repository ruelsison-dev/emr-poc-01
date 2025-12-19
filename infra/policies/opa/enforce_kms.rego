package terraform.policies

# Require that resources storing PHI reference a KMS key

deny[msg] {
  resource := input.resource_changes[_]
  resource.type == "aws_s3_bucket"
  not resource.change.after.server_side_encryption_configuration
  msg = "S3 buckets must enable server-side encryption (SSE-KMS)"
}
