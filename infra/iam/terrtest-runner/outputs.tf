output "role_arn" {
  value = aws_iam_role.terrtest_runner.arn
}

output "role_name" {
  value = aws_iam_role.terrtest_runner.name
}

output "policy_arn" {
  value = aws_iam_policy.terrtest_policy.arn
}