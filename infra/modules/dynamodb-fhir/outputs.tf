output "table_name" {
  value = aws_dynamodb_table.fhir.name
}

output "kms_key_id" {
  value = aws_kms_key.dynamodb.key_id
}

output "kms_key_arn" {
  value = aws_kms_key.dynamodb.arn
}