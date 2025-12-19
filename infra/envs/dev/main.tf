// Dev environment sample wiring (placeholder)
// TODO: Replace with real module calls and references during Phase 1 infra implementation

resource "random_id" "suffix" {
  byte_length = 4
}

module "dynamo" {
  source     = "../../modules/dynamodb-fhir"
  table_name = "dev-fhir-table-${random_id.suffix.hex}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"
  tags = { env = "dev" }
}

module "s3" {
  source = "../../modules/s3-fhir-docs"
  bucket_name = "dev-fhir-docs-${random_id.suffix.hex}"
  kms_key_id  = module.dynamo.kms_key_id
  tags = { env = "dev" }
}

module "iam" {
  source = "../../modules/iam-role-kms"
  name = "dev-fhir-service-role-${random_id.suffix.hex}"
  dynamo_table_arn = module.dynamo.table_name
  s3_bucket_arn = module.s3.bucket
  kms_key_arn = module.dynamo.kms_key_arn
  tags = { env = "dev" }
}
