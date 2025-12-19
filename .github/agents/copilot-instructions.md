# GitHub Copilot Project Context — Patient Administration Sub-System

**Project**: EMR POC — Patient Administration Sub-System
**Last Updated**: 2025-12-20

## Active Technologies
- Backend: Python 3.11, FastAPI
- Database: Amazon Aurora (Postgres-compatible)
- Caching/Locking: Redis (ElastiCache)
- Queueing: Amazon SQS
- Object Storage: Amazon S3 with SSE-KMS
- IAM & KMS: AWS IAM, KMS
- IaC: Terraform (modules), policy-as-code (OPA/Conftest), Terratest
- Testing: pytest, Pact (contracts), k6 (load)
- Monitoring: CloudWatch, CloudTrail, GuardDuty
- Node.js **20.x (LTS)** with **TypeScript** (primary runtime for the web application / REST APIs). This follows the user request to generate the front-end web application and REST APIs using **ExpressJS**. + ExpressJS (with TypeScript), Zod or Joi for request validation, Helmet, express-rate-limit, Winston/pino for structured logging, AWS SDK v3, TypeORM or Prisma (Postgres ORM), Jest + Supertest for tests, Pact for contract testing, ESLint/Prettier, Snyk/npm audit for dependency security. (001-patient-admin)
- Amazon Aurora (PostgreSQL) for registries and transactional data; Amazon S3 (SSE-KMS) for exports and blobs; ElastiCache (Redis) for caching and locks; SQS for async jobs and DLQs. (001-patient-admin)

## Project Structure (from plans)
```
services/
├── registry-service/
├── scheduling-service/
├── notification-service/
infra/
├── modules/
└── envs/
specs/001-patient-admin/
├── spec.md
├── plan.md
├── tasks.md
├── data-model.md
├── contracts/openapi.yaml
├── research.md
└── quickstart.md
```

## Important Commands
- Run unit tests: `pytest tests/unit`
- Run integration tests: `pytest tests/integration`
- IaC checks: `terraform init && terraform validate` (inside `infra/`)
- Policy-as-code: `conftest test plan.out`
- Create feature branch (script): `.specify/scripts/powershell/create-new-feature.ps1 'Description' -ShortName 'patient-admin'`

## Notes for Copilot assistance
- Comply with Constitution v1.1.0: Security, HIPAA, SOC 2 TSC mapping, IaC policy checks.
- Prefer code quality, test-first patterns, and accessibility for UI work (WCAG AA).
- When suggesting AWS services for PHI, prefer only those listed in the project's approved HIPAA-Eligible registry; if unsure, recommend consulting Security.

<!-- MANUAL ADDITIONS START -->

*Add helpful snippets and examples between these markers. Keep additions concise and reviewable.*

<!-- MANUAL ADDITIONS END -->
