---

description: "Task list for Patient Administration Sub-System"
---

# Tasks: Patient Administration Sub-System

**Input**: Design documents from `/specs/001-patient-admin/`  
**Prerequisites**: `plan.md` (required), `spec.md`, Phase 0 research outputs

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and baseline infra

- [ ] T001 Initialize repository structure and developer quickstart (`services/`, `infra/`, `tests/`, `docs/`) — include development environment, linting, and formatting tools.
- [ ] T002 [P] Configure CI pipelines for linting, formatting, unit tests, and secrets scanning.
- [x] T003 [P] Create baseline Terraform modules and skeletons for networking, RDS, S3, KMS, monitoring, and remote state. *(scaffolded placeholders added to `infra/` on 2025-12-20)*
- [ ] T004 [P] Establish access controls for Terraform state (S3 SSE, DynamoDB lock) and secure secrets handling.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infra and security that block feature development

- [ ] T005 Implement Identity/SSO integration (IdP config, RBAC mapping, SCIM provisioning where available).  
- [ ] T006 [P] Implement audit logging and central logging pipeline (CloudTrail → central logging); define retention (6y) and immutability.
- [ ] T007 [P] Configure CI gate: run unit/contract/integration tests and enforce coverage thresholds (blocking merge on failure).
- [ ] T008 [P] Implement policy-as-code checks in CI for Terraform plans (OPA/AWS Config) and unblock workflow with manual approvals for production applies.
- [ ] T009 Set up monitoring & SLO framework (Dashboards, Alerts, Targets for booking, queries, notification delivery).
- [ ] T010 [P] Implement Terraform module tests (Terratest or equivalent) and test harness for module behaviors.

**Checkpoint**: Foundational infra and security gates complete — user story work may proceed.

---

## Phase 3: User Story 1 - Registry & Appointment Booking (Priority: P1) 🎯 MVP

**Goal**: Implement registry CRUD and appointment booking flow with conflict detection and audit logging

**Independent Test**: Create, modify, and cancel an appointment; verify resource allocation and audit entries.

### Tests for User Story 1 (must be written first)

- [ ] T011 [P] [US1] Unit tests for registry models and validation (tests/unit/test_registry_*.py)
- [ ] T012 [P] [US1] Contract tests for Registry API endpoints (tests/contract/test_registry_api.py)
- [ ] T013 [P] [US1] Integration tests for booking lifecycle and conflict resolution (tests/integration/test_booking_lifecycle.py)

### Implementation for User Story 1

- [ ] T014 [P] [US1] Implement `registry-service` models and CRUD endpoints (services/registry-service/src)
- [ ] T015 [P] [US1] Implement `scheduling-service` booking APIs and conflict detection (services/scheduling-service/src)
- [ ] T016 [US1] Implement idempotency and transactional safety for booking operations
- [ ] T017 [US1] Add audit logging hooks for create/modify/delete/read actions (centralized logging)
- [ ] T018 [US1] Implement resource allocation and locking (Redis + DB coordination)
- [ ] T019 [US1] Add basic UI or admin endpoints for manual scheduling and conflict resolution (optional lightweight UI)

**Checkpoint**: US1 should be functional and testable independently

---

## Phase 4: Notifications, Consent & Security (Priority: P1)

**Goal**: Implement notification workflows with consent model and secure delivery

- [ ] T020 [P] Implement notification-service with consent store and verification of recipient reachability
- [ ] T021 Integrate notification delivery channels (email/SMS proxies) with safe defaults; ensure messages with PHI only sent when consented and via verified/encrypted channels
- [ ] T022 Add retry/exponential backoff for delivery failures and health checks for delivery providers
- [ ] T023 Add tests for consent flows, notification content rules, and audit logs (tests/integration/test_notifications.py)

---

## Phase 5: Exports & Queries (Priority: P2)

- [ ] T024 Implement export API to provide FHIR R4 JSON for authorized exports; enforce access controls and audit exports
- [ ] T025 [P] Add query endpoints for clinician views with pagination, filters, and performance tests
- [ ] T026 Add contract tests for export formats and query endpoints

---

## Phase 6: Performance, DR & Compliance (Priority: P2)

- [ ] T027 Performance testing: load tests for 10 req/s and p95 targets for queries (use k6 or similar)
- [ ] T028 Backup/restore drills for registries (document RTO, RPO, and validate restores)
- [ ] T029 Run compliance evidence collection tasks for SOC 2 and HIPAA (mapping doc, sample artifacts, retention proof)
- [ ] T030 Conduct security testing (SAST, DAST) and threat model sign-off

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T031 Accessibility & UX QA (WCAG AA checks)
- [ ] T032 Documentation updates (quickstart, runbooks, operations guides)
- [ ] T033 Cost estimation and tagging enforcement for infra resources
- [ ] T034 Final operational readiness review and production cutover plan

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
