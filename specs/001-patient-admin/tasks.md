---

description: "Task list for Patient Administration Sub-System"
---

# Tasks: Patient Administration Sub-System

**Input**: Design documents from `/specs/001-patient-admin/`  
**Prerequisites**: `plan.md` (required), `spec.md`, Phase 0 research outputs

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and baseline infra

- [ ] T001 Initialize repository structure and developer quickstart (create `backend/express-api/`, `backend/fhir-service-go/`, `infra/`, `specs/001-patient-admin/docs/`, `tests/`) — file paths: `backend/`, `infra/`, `specs/001-patient-admin/quickstart.md`
- [ ] T002 [P] Configure CI pipelines for linting, formatting, unit tests, dependency scanning, and secrets scanning (add workflows in `.github/workflows/ci-express.yml`, `.github/workflows/ci-go.yml`, `.github/workflows/iac.yml`)
- [x] T003 [P] Create baseline Terraform modules and skeletons for networking, RDS, S3, KMS, monitoring, and remote state (files: `infra/modules/*`) *(scaffolded placeholders added 2025-12-20)*
- [ ] T004 [P] Establish secure remote state and access controls (S3 SSE, DynamoDB lock, IAM policies) — files: `infra/envs/dev/remote-state.tf`, `infra/modules/iam-role-terraform/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infra and security that MUST be completed before user story work

- [ ] T005 Implement Identity/SSO integration and RBAC mapping (IdP config, SCIM where available) — files: `docs/identity.md`, `backend/express-api/src/middleware/auth.ts`
- [ ] T006 [P] Implement audit logging and central logging pipeline (CloudTrail → central logging) and retention policy (6y) — files: `infra/modules/logging/`, `specs/001-patient-admin/docs/compliance_evidence/logging.md`
- [ ] T007 [P] Configure CI gates to run unit, contract, and integration tests; enforce coverage thresholds and fail on regression — files: `.github/workflows/ci-*.yml`, `.github/workflows/iac.yml`
- [ ] T008 [P] Implement policy-as-code checks for Terraform plans (OPA/Conftest) in CI and add blocking manual approvals for production applies — files: `.github/workflows/iac.yml`, `infra/policies/`
- [ ] T009 Set up monitoring, SLOs, dashboards, and runbooks for core flows (booking, queries, notifications) — files: `docs/observability.md`, `infra/modules/monitoring/`
- [ ] T010 [P] Implement Terraform module tests (Terratest) for critical modules and create a runbook for safe CI execution (`infra/terratest/`, `docs/terratest-runbook.md`)
- [ ] T050 [P] Implement MFA & justification logging for privileged actions — design and implement enforcement for privileged IAM actions and app-layer privileged operations (files: `infra/modules/iam-privileged/`, `backend/*/src/middleware/mfa.ts`); add integration tests `backend/*/tests/integration/test_mfa_justification.*` and collect compliance evidence in `specs/001-patient-admin/docs/compliance_evidence/mfa/`. **Security sign-off required before Phase 3 user story merges.**
- [ ] T051 [P] Provision a Terratest runner IAM role with OIDC trust and least-privilege policy (`infra/iam/terrtest-runner.tf`, `infra/modules/iam-role-kms/`) and document the approval workflow and budget guardrails (`docs/terratest-runbook.md`). Add Terratest validation to assert the role's permissions in `infra/terratest/iam_role_test.go`.
- [ ] T052 [P] Implement autoscaling & SLO verification: add autoscaling and compute modules and deployment tuning (`infra/modules/autoscaling/`, `infra/modules/compute/`), add load test baseline and verification (`k6/scripts/booking_baseline.js`) and a task to iterate until SLOs are met (acceptance: 10 req/s sustained with <1% errors and p95 ≤ 500ms). Tie results to `T042` load tests and require performance gating for production deploys.

**Checkpoint**: Foundational infra and security gates complete — user story work may proceed.

---

## Phase 3: User Story 1 - Registry & Appointment Booking (Priority: P1) 🎯 (MVP)

**Goal**: Implement registry CRUD and appointment booking flow with reservation tokens, conflict resolution, and audit logging

**Independent Test**: Create, modify, cancel an appointment via API or UI; verify resource allocation, audit logs, and notification suppression on consent revocation.

### Tests (write first)

- [ ] T011 [P] [US1] Unit tests for models and validation (`backend/express-api/src/tests/unit/**`)
- [ ] T012 [P] [US1] Contract tests (Pact) for Registry and Appointment APIs (`tests/contract/pact_registry/`, `tests/contract/pact_appointments/`)
- [ ] T013 [P] [US1] Integration tests for booking lifecycle and conflict resolution (`backend/express-api/src/tests/integration/test_booking_lifecycle.ts`)

### Implementation

- [ ] T014 [P] [US1] Initialize `backend/express-api` skeleton (files: `backend/express-api/package.json`, `backend/express-api/tsconfig.json`, `backend/express-api/README.md`)
- [ ] T015 [P] [US1] Implement `Patient` model and CRUD controller (`backend/express-api/src/models/patient.ts`, `backend/express-api/src/controllers/patients.ts`)
- [ ] T016 [P] [US1] Implement `Person`, `Location`, `Organization` models and controllers (`backend/express-api/src/models/*.ts`, `backend/express-api/src/controllers/*.ts`)
- [ ] T017 [US1] Implement `Appointment` endpoints and `BookingService` (`backend/express-api/src/controllers/appointments.ts`, `backend/express-api/src/services/booking_service.ts`)
- [ ] T018 [US1] Implement `Schedule` and `Slot` generation (`backend/express-api/src/controllers/schedules.ts`, `backend/express-api/src/models/slot.ts`)
- [ ] T019 [US1] Implement reservation token holds with TTL and Redis locking (`backend/express-api/src/lib/holds.ts`, `backend/express-api/src/services/hold_manager.ts`)
- [ ] T020 [US1] Implement idempotency (idempotency keys) and transactional safety for booking operations (`backend/express-api/src/lib/idempotency.ts`)
- [ ] T021 [US1] Add audit logging hooks and ensure events are emitted to central logging (`backend/express-api/src/lib/audit.ts`)
- [ ] T022 [US1] Implement resource allocation and locking integration (Redis + DB) and add Terratest checks that verify infra configuration (`backend/express-api/src/services/resource_manager.ts`, `infra/terratest/booking_tests.go`)
- [ ] T023 [US1] Add Pact consumer/provider tests connecting Express APIs to FHIR backend contracts (`tests/contract/pact_appointments/`)
- [ ] T024 [US1] Add CI checks for Node: Snyk/npm audit, govulncheck (or equivalent), and DAST pre-production step (`.github/workflows/ci-express.yml`)

**Checkpoint**: US1 should be fully functional and independently testable (unit, contract, integration).

## Phase 3b: User Story 2 - FHIR Backend (Go) - Integration (Priority: P1)

**Goal**: Provide FHIR R4 compatible backend for storage and exports (DynamoDB + S3) with secure access

**Independent Test**: Provider contract tests + LocalStack integration verifying DynamoDB and S3 read/write and SSE-KMS enforcement.

### Tests (write first)

- [ ] T025 [P] [US2] Unit tests for Go handlers and validation (`backend/fhir-service-go/internal/**`)
- [ ] T026 [P] [US2] LocalStack integration tests for DynamoDB and S3 (`backend/fhir-service-go/tests/integration/**`)
- [ ] T027 [P] [US2] Pact contract tests where Express is consumer and Go FHIR service is provider (`tests/contract/pact_fhir/`)

### Implementation for FHIR Backend

- [ ] T028 [P] [US2] Initialize `backend/fhir-service-go` skeleton and README (`backend/fhir-service-go/cmd/main.go`, `backend/fhir-service-go/go.mod`)
- [ ] T029 [US2] Implement DynamoDB access patterns and table schema in `infra/modules/dynamodb-fhir/` and `backend/fhir-service-go/internal/store/dynamodb.go` (SSE-KMS, TTL)
- [ ] T030 [US2] Implement S3 exports and signed URL downloads in `backend/fhir-service-go/internal/store/s3.go` and `infra/modules/s3-fhir-docs/` (SSE-KMS)
- [ ] T031 [US2] Add LocalStack-based integration CI job (`.github/workflows/integration-localstack.yml`) and tests
- [ ] T032 [US2] Implement IAM least-privilege role and KMS policies (`infra/modules/iam-role-kms/`) and add Terratest IAM assertions (`infra/terratest/iam_policy_test.go`)
- [ ] T033 [US2] Add Go static analysis and vulnerability checks (`backend/fhir-service-go/.github/workflows/ci-go.yml`)

**Checkpoint**: US2 is independently testable and provides provider contracts for US1.

**Checkpoint**: US1 should be functional and testable independently

---

## Phase 4: Notifications, Consent & Security (Priority: P1)

**Goal**: Implement notification workflows with consent model and secure delivery

**Independent Test**: Consent grant/revoke flows enforce suppression and generate audit events; notifications sent according to consent and only include PHI when allowed.

- [ ] T034 [P] [US3] Implement consent record model and APIs (`backend/express-api/src/models/consent.ts`, `backend/express-api/src/controllers/consents.ts`)
- [ ] T035 [US3] Implement `notification-service` with adapters for SES/SMS and in-app delivery (`backend/notification-service/`)
- [ ] T036 [US3] Add tests for consent enforcement and notification content rules (`backend/notification-service/tests/**`)
- [ ] T037 [US3] Implement retry/exponential backoff and monitoring for delivery failures (`backend/notification-service/src/lib/retries.ts`)
- [ ] T038 [US3] Document consent and notification policies (`specs/001-patient-admin/docs/consent_policy.md`) and include suppression runbook (`docs/runbooks/suppression.md`)

---

## Phase 5: Exports & Queries (Priority: P2)

- [ ] T039 [US4] Implement export API for FHIR R4 JSON in `backend/fhir-service-go/internal/export/` and `backend/express-api/src/controllers/exports.ts` — ensure access controls and audit logging
- [ ] T040 [US4] Implement query endpoints and clinician views in `backend/express-api/src/controllers/queries.ts` with pagination and filters
- [ ] T041 [US4] Add performance tests for clinician queries using `k6/scripts/queries.js` and integrate into CI (`.github/workflows/load-tests.yml`)

---

## Phase 6: Performance, DR & Compliance (Priority: P2)

- [ ] T042 Run load tests for scheduling endpoints and validate p95 targets (add `k6` job: `.github/workflows/load-tests.yml`)
- [ ] T043 Implement backup & restore drills for registries and document RTO/RPO (`docs/drills/restore_playbook.md`)
- [ ] T044 Run compliance evidence collection and map controls to SOC 2 TSC and HIPAA (`specs/001-patient-admin/docs/compliance_evidence/`)
- [ ] T045 Conduct security testing (SAST, DAST) and threat model sign-off (`.github/workflows/security.yml`)

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T046 Accessibility & UX QA (WCAG AA checks) — `frontend/` or server-rendered pages
- [ ] T047 Documentation updates (`specs/001-patient-admin/quickstart.md`, `docs/runbooks/`)
- [ ] T048 Cost estimation and tagging enforcement for infra resources (`infra/tags.tf`)
- [ ] T049 Final operational readiness review and production cutover plan (`docs/release/cutover.md`)

---

## Dependencies & Execution Order
- **Phase 1 & 2** must complete before all user story work.  
- Tests for a user story **MUST** be added and failing before implementation (TST first).  
- IaC checks **MUST** run before applying to staging/prod.  
- Compliance tasks can be run in parallel, but final audit evidence must be available before production release.

---

## Notes
- Include SOC 2 TSC mapping and HIPAA control checklists in `/specs/001-patient-admin/docs/` and attach evidence artifacts to the compliance folder for audit readiness.
- Add follow-up task to record exact IdP and SCIM mappings once the organization selects an IdP.

---

## Parallel Execution Examples

- Run model + unit tests for Express and Go in parallel:
  - `pnpm -w -r test --filter ./backend/express-api ./backend/fhir-service-go`
- Run Terratest infra tests while developers implement models for US1:
  - CI job: `.github/workflows/infra-terratest.yml` runs `go test ./infra/terratest/...`

---

## Implementation Strategy

**MVP First (User Story 1 only)**
1. Complete Phase 1 (Setup).
2. Complete Phase 2 (Foundational) and obtain Architecture & Security sign-off.
3. Implement US1 with tests-first approach (unit → contract → integration).
4. Validate US1 independently and demo as MVP.

**Incremental & Parallel Delivery**
- After US1 is validated, implement US2 (FHIR backend) and US3 (Notifications) in parallel as capacity allows.
- Keep iterations small and ensure contract tests pass for each PR.

---

## Validation Summary

- Total tasks: 49 (T001–T049)
- Tasks per story:
  - US1 (Registry & Booking): 14 tasks (T011–T024)
  - US2 (FHIR Backend): 6 tasks (T025–T033)
  - US3 (Notifications/Consent): 5 tasks (T034–T038)
  - US4 (Exports & Queries): 3 tasks (T039–T041)
  - Foundational & Setup: 10 tasks (T001–T010)
  - Performance/DR/Polish: 11 tasks (T042–T049)
- Parallel opportunities identified: many [P] tasks (CI, infra, module tests, unit tests) and user-story-level parallelism once foundational phase is complete
- Independent test criteria for each story included (see each story's section above)
- Suggested MVP scope: **US1 (Registry & Booking)**

**Format validation**: All tasks conform to the speckit checklist format (`- [ ] T### [P?] [USx?] Description with file path`).

