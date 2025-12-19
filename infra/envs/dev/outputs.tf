output "table_name" {
  value = module.dynamo.table_name
}

output "kms_key_arn" {
  value = module.dynamo.kms_key_arn
}

output "bucket" {
  value = module.s3.bucket
}

output "role_arn" {
  value = module.iam.role_arn
}

output "policy_arn" {
  value = module.iam.policy_arn
}