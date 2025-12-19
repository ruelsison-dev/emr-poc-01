terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
}

variable "oidc_provider_arn" {
  type        = string
  description = "ARN of the OIDC provider (e.g., arn:aws:iam::<acct>:oidc-provider/token.actions.githubusercontent.com)"
}

variable "github_repo" {
  type        = string
  description = "GitHub repo in the form owner/repo (used in role condition)"
}

variable "github_ref" {
  type        = string
  default     = "refs/heads/*"
  description = "GitHub ref pattern allowed to assume the role"
}

variable "role_name" {
  type    = string
  default = "terrtest-runner"
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Federated"
      identifiers = [var.oidc_provider_arn]
    }

    actions = ["sts:AssumeRoleWithWebIdentity"]

    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:${var.github_repo}:${var.github_ref}"]
    }
  }
}

resource "aws_iam_role" "terrtest_runner" {
  name               = var.role_name
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
  tags = {
    created_by = "terratest"
    purpose    = "terratest-runner"
  }
}

# Example least-privileged policy for test runs. Narrow this in production.
resource "aws_iam_policy" "terrtest_policy" {
  name        = "terrtest-runner-policy-${replace(var.role_name, "-", "-") }"
  description = "Least-privilege-ish policy for terratest runner. Narrow resource ARNs before production use."

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:GetBucketLocation",
          "s3:ListBucket",
          "s3:GetObject",
          "s3:PutObject",
          "dynamodb:DescribeTable",
          "dynamodb:Query",
          "dynamodb:PutItem",
          "kms:DescribeKey",
          "kms:Encrypt",
          "kms:Decrypt",
          "kms:GenerateDataKey"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "terrtest_attach" {
  role       = aws_iam_role.terrtest_runner.name
  policy_arn = aws_iam_policy.terrtest_policy.arn
}

output "role_arn" {
  value = aws_iam_role.terrtest_runner.arn
}

output "role_name" {
  value = aws_iam_role.terrtest_runner.name
}

output "policy_arn" {
  value = aws_iam_policy.terrtest_policy.arn
}
