# Terratest Runbook & CI Plan

Purpose: Provide a safe, auditable, and repeatable process to run Terratest-based integration tests for Terraform modules and infra in CI or locally. This runbook covers prerequisites, CI configuration, IAM & security gating, run steps, artifact collection, cost controls, and acceptance criteria.

---

## 1) Summary & Scope

- Scope: `infra/modules/*` critical modules (DynamoDB, S3, KMS, IAM roles). Tests live in `infra/terratest/`.
- Goal: Validate that Terraform modules create resources with required encryption (SSE-KMS), public-block settings, and least-privilege IAM attachments.
- Safe mode: Prefer LocalStack for fast validation; use dedicated **dev AWS account** for end-to-end verification only after Security sign-off.

---

## 2) Prerequisites

Local (developer):
- Go >= 1.21 (for Terratest)
- Terraform CLI (matching project version)
- AWS CLI (for authenticated runs)
- Optional: LocalStack for service emulation

CI (GitHub Actions):
- OIDC trust configured between GitHub and target AWS account (recommended) or a short-lived credential (approved secret) scoped to least-privilege role
- A dedicated `terrtestrunner` role in the target dev account with scoped permission to create/destroy only test-scoped resources in a dedicated test prefix
- Terraform state remote-backend pre-configured for the test environment (S3 bucket + DynamoDB lock)

Security prerequisites:
- Security & Architecture sign-off to run Terratest in the target dev account (must include review of IAM role policy and billing guardrails)
- Budget guardrails (AWS budget alarms for test account)

---

## 3) IAM & Least-Privilege Guidance (high level)

- Create a test role with a trust policy for GitHub Actions OIDC and allow `sts:AssumeRole` only from the repository/organization and branch(es) used for testing. **Provisioning this role is tracked in task `T051` in `specs/001-patient-admin/tasks.md`.**
- Minimum permissions for Terratest runs (examples, restrict to test tag/prefix):
  - `ec2:*`/`s3:*`/`dynamodb:*` only for names with prefix `test-<run-id>-*`
  - `kms:CreateKey`, `kms:PutKeyPolicy`, `kms:DescribeKey` scoped to test key resources
  - `iam:PassRole` only for roles created by the tests; prefer pre-created roles when possible
  - `sts:AssumeRole` to assume the `terrtestrunner` role
- Require MFA for more-privileged manual approvals when running against shared infra
- Keep role policies minimal and time-limited

---

## 4) CI Configuration (GitHub Actions) — recommended pattern

Key ideas:
- Use OIDC to obtain short-lived role credentials in the target AWS dev account
- Require an explicit label or review from Security/Architecture before running `infra-terratest` against the dev account
- Provide a `dry-run` job (LocalStack) that runs on every PR; `dev-run` job runs on explicit approval

Suggested workflow snippet (conceptual):

```yaml
# .github/workflows/infra-terratest.yml
on: workflow_dispatch
jobs:
  terratest-localstack:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run terratest (LocalStack)
        run: make test-localstack # runs `go test ./infra/terratest -run Integration`

  terratest-dev:
    if: ${{ github.event.inputs.run_in_dev == 'true' }} # only runs by manual trigger + approval
    runs-on: ubuntu-latest
    permissions:
      id-token: write
    steps:
      - uses: actions/checkout@v4
      - name: Configure AWS creds (OIDC)
        uses: aws-actions/configure-aws-credentials@v2
        with:
          role-to-assume: arn:aws:iam::123456789012:role/terrtest-runner
          aws-region: us-east-1
      - name: Run terratest (dev)
        env:
          TF_VAR_test_run_id: ${{ github.run_id }}
        run: make test-dev # runs `go test ./infra/terratest -run Integration -count=1`

# Important: require a SECURITY/ARCHITECTURE review before setting `run_in_dev=true` on a workflow dispatch
```

Notes:
- Use `github.run_id` to tag resources with deterministic names (e.g., `test-${GITHUB_RUN_ID}-dynamo`) so tear-down is easy
- Keep `-count=1` to avoid flaky reruns masking issues

---

## 5) Running Terratest locally

Developers should follow this sequence:
1. Ensure local AWS CLI credentials point to a dedicated test account or run LocalStack
2. From repo root: `cd infra/terratest && go test -v ./... -run TestIntegration -timeout 30m`
3. Use `TF_VAR_test_run_id` environment variable to namespace resources: `export TF_VAR_test_run_id=local123`
4. Inspect artifacts and logs (`go test` output), and run `terraform destroy` if tests leave resources behind (tests should clean up automatically)

Quick local LocalStack run (fast verification):
- `docker-compose -f dev/localstack/docker-compose.yml up -d`
- Run `make test-localstack` or `go test` against mocked endpoints

---

## 6) Artifacts & Evidence Collection

Ensure these artifacts are uploaded as workflow artifacts or to a secure bucket:
- `go test` output (full logs)
- `terraform plan` output (plain text and JSON)
- Test run metadata: run id, commit SHA, branch, workflow URL
- Terraform state snapshot (if ephemeral) — avoid leaking secrets
- Any debug artifacts (failed assertion stack traces, AWS CLI debug logs)

Attaching to PR:
- Post a comment summarizing test results and add artifact links and a short checklist for Security review

---

## 7) Timeouts, Retries & Cost Controls

- Job timeout: set CI job timeout (e.g., 60–90 minutes) to avoid runaway costs
- Terratest timeout: use `-timeout` argument (e.g., `-timeout 30m`) and default retries only for idempotent tests
- Tag all test resources with `test:true`, `run-id:<id>`, and `created-by:terratest`
- Create an automatic cleanup job (scheduled) to destroy resources older than a configured TTL (e.g., 24 hours)
- Set AWS budget alarms on the dev test account and require Security notification on threshold

---

## 8) Failure Handling & Rollback

- Tests MUST attempt to destroy created resources in `defer` blocks in Go tests. If teardown fails, collect resource IDs and escalate to manual cleanup.
- If a module test fails, the CI job should:
  - Upload logs and Terraform plan output
  - Create a PR comment with findings and attach artifacts
  - Optionally `terraform destroy` to remove partially created resources (if safe/approved)

---

## 9) Approval & Gating

- **Gating requirement**: Running `terratest-dev` against an actual AWS account requires explicit Architecture & Security approval (recorded as a PR review or an explicit comment from approvers)
- Add a checklist item to the PR for approvals and include the `run-id` for traceability

---

## 10) Acceptance Criteria (Test Run)

A successful Terratest run is one where:
- All Terratest assertions pass (no failed assertions)
- Terraform plan and apply were validated for the expected configuration (encryption, public-block settings, KMS policies)
- All created resources are destroyed at the end of the run (or flagged for cleanup and approved for retention)
- Artifacts (test logs, plan output) are uploaded as CI artifacts and linked in PR
- No excessive permissions were used (audit of the assumed role shows only expected actions)

---

## 11) Checklist (Quick)

- [ ] Confirm Security & Architecture approval for dev account
- [ ] Ensure OIDC role exists and trust policy is configured
- [ ] Ensure TF remote state backend is configured for test runs
- [ ] Run LocalStack dry-run on PR automation (every PR)
- [ ] On manual trigger (after approval), run `terratest-dev` workflow and collect artifacts
- [ ] Validate acceptance criteria and attach evidence to PR

---

## 12) Notes & Next Steps

- Consider adding a dedicated `test-runner` AWS account for destructive tests and cross-account read-only views into staging/production accounts for verification without write access.
- Keep the runbook in `docs/` and reference it from the PR template for infra changes.

---

If you'd like, I can commit this runbook (create `docs/terratest-runbook.md`), add a small PR with the change, and prepare the GitHub Actions snippet and an example IAM role policy for review by Security.