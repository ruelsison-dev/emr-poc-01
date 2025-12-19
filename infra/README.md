# Infra Scaffolding — Patient Administration Sub-System

This directory contains initial Terraform module placeholders and policy-as-code examples to bootstrap secure, compliant infrastructure for the Patient Administration Sub-System.

Goals:
- Provide secure defaults (encryption, tagging, no public access)
- Use reusable, versioned Terraform modules (module boundaries defined below)
- Integrate policy-as-code checks (OPA) and Terraform plan validation in CI
- Include a `dev` environment backend example; production backends require manual approvals

Modules (initial placeholders):
- `network` — VPC, subnets, routing, NAT (if required)
- `rds` — managed relational DB module (RDS/Aurora)
- `s3` — buckets with SSE-KMS and lifecycle rules
- `kms` — KMS keys, rotation policies, and grants
- `monitoring` — dashboards, alerts, log routing
- `service-module` — patterns for service-level infra (IAM roles, security groups)

Policy-as-code examples live in `infra/policies/opa/` and should be expanded and executed by CI as part of `terraform plan` validation.

See `infra/tests/README.md` for guidance on testing modules (Terratest/Terraform compliance tests).
