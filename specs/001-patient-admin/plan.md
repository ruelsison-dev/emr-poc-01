# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

**Language/Version**: Node.js **20.x (LTS)** with **TypeScript** (primary runtime for the web application / REST APIs). This follows the user request to generate the front-end web application and REST APIs using **ExpressJS**.
**Primary Dependencies**: ExpressJS (with TypeScript), Zod or Joi for request validation, Helmet, express-rate-limit, Winston/pino for structured logging, AWS SDK v3, TypeORM or Prisma (Postgres ORM), Jest + Supertest for tests, Pact for contract testing, ESLint/Prettier, Snyk/npm audit for dependency security.
**Storage**: Amazon Aurora (PostgreSQL) for registries and transactional data; Amazon S3 (SSE-KMS) for exports and blobs; ElastiCache (Redis) for caching and locks; SQS for async jobs and DLQs.
**Testing**: Unit tests (Jest), integration tests (Supertest + local DB via Docker), contract tests (Pact), IaC tests (Terratest), policy-as-code (Conftest/OPA), security scanning (Snyk/npm audit, dependency-check), DAST in CI.
**Target Platform / Deployment**: AWS ECS (Fargate) or AWS App Runner as primary long-running runtime pattern (containerized). Lambdas may be used for lightweight async tasks where appropriate. Deploy with Terraform modules, ECS service definition, and ALB/NLB in private subnets.
**Project Type**: Web application + API (ExpressJS acts as API server and can serve server-side rendered pages or static front-end assets if needed). Consider a separate SPA (React/Next.js) for richer client UI in Phase 1 if desired.
**Performance Goals**: Maintain spec targets: support 10 req/s for scheduling endpoints and p95 <= 500ms for queries; booking flows (end-to-end) <= 2 minutes for batch operations. Adopt autoscaling policies and connection pooling for DB.
**Constraints**: Must comply with HIPAA & SOC 2 gates — only HIPAA-Eligible AWS services for PHI; KMS for key management, TLS 1.2+; logging and retention (6 years); least-privilege IAM. Ensure Node.js security best practices (no eval, limit headers, input validation, strict dependency vetting).
**Scale/Scope**: Initial pilot: support N clinicians in a single region (TBD). Capacity planning to be done in Phase 1 and load-tested in staging.

**Notes / NEEDS CLARIFICATION**:
- Does ExpressJS **replace** the previously chosen Python/FastAPI services for this feature, or is ExpressJS intended as the front-end/API gateway that forwards to existing FastAPI microservices? (Answer required — affects code ownership and migration plan)
- Confirm Node.js version and whether to standardize on TypeScript (strongly recommended) or permit JavaScript for prototyping.
- Hosting preference: ECS Fargate vs App Runner vs Lambda for the primary Express service (recommend ECS Fargate unless short-lived functions are required).
- Confirm whether SPA client (React/Next.js) will be used or the Express app should render server-side assets.
- Determine required throughput targets (beyond the 10 req/s for scheduling) for capacity planning in Phase 1.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- All features **MUST** satisfy the Constitution's mandatory gates before Phase 0 completes. Any deviation requires documented justification and an approved exception.
- Required gates include (non-exhaustive):
  - **Security & Compliance**: For systems handling regulated or sensitive data (e.g., PHI), designs **MUST** reference only AWS services that are *HIPAA-Eligible* and in-scope for SOC 2, demonstrate encryption at-rest and in-transit, least-privilege IAM, audit logging, and identify required BAAs. Our design follows these constraints (Aurora, S3 SSE-KMS, SQS, CloudTrail). **Note**: the runtime change to Node/Express requires Security/Compliance review to ensure toolchains and libraries meet scanning/security requirements for PHI workloads.
  - **Testing & Quality**: Unit, integration, and contract tests **MUST** be present for new behavior. CI must enforce test gates and target coverage thresholds; failing tests block progression. We will use Jest/Supertest/Pact for Node services and Terratest/Conftest for IaC.
  - **Operational Readiness**: Monitoring, alerting, runbooks, SLO/SLA targets, deployment and rollback plans **MUST** be documented. Our plan includes CloudWatch dashboards, runbooks for revocation handling and operational incidents.
  - **IaC & Policy Validation**: Infrastructure **MUST** be declared as code (Terraform preferred), use reusable modules with semantic versioning, and run policy-as-code checks (e.g., OPA/Sentinel/AWS Config) in CI before apply. The plan includes module templates and OPA examples.
  - **Performance & Efficiency**: Performance goals, budgets, and load-test plans **MUST** be defined; capacity/resizing strategy included. Scheduling endpoints already include p95 and throughput targets in the spec.
- **Outstanding gate**: Architecture and Security approval required for the runtime/stack change (from Python->Node) and for adding Node-specific CI/DAST tools to the pipeline. This must be resolved before Phase 1 implementation begins.
- Consult the Constitution (Version [CONSTITUTION_VERSION]) for the full set of mandatory checks that apply to regulated systems and PHI handling.

## Project Structure

### Documentation (this feature)

```text
specs/001-patient-admin/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

## Phase 0: Research & Clarifications
Objective: Resolve outstanding technical unknowns that materially affect the implementation approach.

Tasks:
1. Decide whether **ExpressJS (Node)** replaces or coexists with **FastAPI (Python)** for core services (owner: Product/Arch). [REQUIRED]
2. Confirm Node.js runtime version and TypeScript adoption (owner: Eng Lead). [REQUIRED]
3. Select hosting pattern for Express API (ECS Fargate vs App Runner vs Lambda) and document pros/cons (owner: Infra). [REQUIRED]
4. Identify security hardening checklist for ExpressJS handling PHI (Helmet, CSP, rate-limiting, input validation, dependency scanning) and integrate to CI (owner: Security). [REQUIRED]
5. Finalize API surface OpenAPI contract (`contracts/patient-admin-openapi.yaml`) and ensure contract tests (Pact) are planned (owner: API Lead). [REQUIRED]

Exit criteria: All items above marked resolved and recorded in `research.md`. Once complete, proceed to Phase 1 design work (data models, API contracts, module skeletons).
### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Selected structure: Web application + API (ExpressJS)
backend/express-api/
├── src/
│   ├── controllers/     # route handlers
│   ├── routes/          # route definitions
│   ├── services/        # business logic
│   ├── models/          # ORM models (Prisma/TypeORM)
│   ├── middlewares/     # auth, validation, error handling
│   ├── lib/             # shared utilities
│   └── tests/           # unit + integration tests (jest, supertest)
frontend/                # optional SPA served by separate build (React/Next.js recommended)
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# Structure Decision:
- Use `backend/express-api` as the primary implementation for the REST API endpoints requested (Patient, Person, Location, Organization, Encounter, Appointment, Schedule, Slot, Questionnaire).
- Optionally implement a separate `frontend/` SPA (React/Next.js) if richer client-side UX is required; for initial delivery the Express server may serve server-rendered pages or static assets.
- Keep the data access layer and domain models colocated under `backend/express-api` and expose a clear `contracts/` directory with OpenAPI specs for each public surface.

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Change of primary runtime from Python/FastAPI to Node.js/Express (user-requested) | The user requested an ExpressJS-based web application and REST API; implementing in Node.js aligns with front-end developer skillset and user preference. | Retaining FastAPI would require rewriting the front-end stacks or adding an API gateway layer; the chosen path reduces integration impedance for a JS-first team but requires Architecture/Compliance approval for runtime change and security sign-off. |

**Action Required**: Seek Architecture Owner approval for the runtime change; record decision and migration plan if legacy FastAPI services must be retained or phased out.
