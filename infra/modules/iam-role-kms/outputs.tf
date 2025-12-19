output "role_arn" {
  value = aws_iam_role.service_role.arn
}

output "policy_arn" {
  value = aws_iam_policy.fhir_policy.arn
}