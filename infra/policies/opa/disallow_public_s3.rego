package terraform.policies

# Deny S3 buckets with public access block not set or disabled

deny[msg] {
  resource := input.resource_changes[_]
  resource.type == "aws_s3_bucket_public_access_block"
  not resource.change.after.block_public_acls
  msg = "S3 public ACLs must be blocked"
}
