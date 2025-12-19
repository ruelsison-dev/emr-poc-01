resource "aws_kms_key" "dynamodb" {
  description             = "KMS key for DynamoDB table encryption for FHIR resources"
  deletion_window_in_days = 30
}

resource "aws_kms_alias" "dynamodb_alias" {
  name          = "alias/dynamodb-fhir"
  target_key_id = aws_kms_key.dynamodb.key_id
}

resource "aws_dynamodb_table" "fhir" {
  name           = var.table_name
  billing_mode   = var.billing_mode
  hash_key       = var.hash_key
  attribute {
    name = var.hash_key
    type = "S"
  }

  # optional GSI for resource type/time queries
  dynamic "global_secondary_index" {
    for_each = var.gsis
    content {
      name               = global_secondary_index.value.name
      hash_key           = global_secondary_index.value.hash_key
      range_key          = global_secondary_index.value.range_key
      projection_type    = global_secondary_index.value.projection_type
      read_capacity      = global_secondary_index.value.read_capacity
      write_capacity     = global_secondary_index.value.write_capacity
    }
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.dynamodb.arn
  }

  tags = var.tags
}
