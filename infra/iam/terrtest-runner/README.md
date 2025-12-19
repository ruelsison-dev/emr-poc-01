Terratest Runner Role

This module creates an IAM role intended for short-lived test runs (e.g., GitHub Actions via OIDC). It is deliberately conservative in scope for the initial pass and MUST be narrowed for production use.

Usage example (dev/test):

```
module "terrtest-runner" {
  source = "../../iam/terrtest-runner"

  oidc_provider_arn = "arn:aws:iam::123456789012:oidc-provider/token.actions.githubusercontent.com"
  github_repo       = "owner/repo"
  role_name         = "terrtest-runner-${var.test_run_id}"
}
```

Security notes:
- Narrow policy resources before production use
- Use OIDC trust with conditions scoped to the repository and branch(es)
- Require Security/Architecture approval before enabling runs against shared accounts
