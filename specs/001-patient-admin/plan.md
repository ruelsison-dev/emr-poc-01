# Implementation Plan: Patient Administration Sub-System

**Branch**: `001-patient-admin` | **Date**: 2025-12-20 | **Spec**: ../spec.md
**Input**: Feature specification for Patient Administration Sub-System (registries, scheduling, notifications, queries). 

## Summary

Deliver an AI-native Patient Administration Sub-System that provides robust, auditable registries (Patient, Person, Provider, Location, Organization) and a scalable scheduling system (requesting, booking, modification, notifications, resource allocation, queries). The implementation will prioritize HIPAA and SOC 2 compliance, run as microservices or bounded services, use Terraform modules for infrastructure, and include automated CI gates, policy-as-code checks, monitoring and runbooks. The initial MVP focuses on registry CRUD, appointment booking flows (P1 clinician workflows), notification opt-in/consent, and query APIs.

## Technical Context

**Language/Version (RECOMMENDATION — NEEDS CLARIFICATION)**: Backend: Python 3.11 (FastAPI) or an equivalent HTTP framework; Workers: same runtime or serverless workers (consistent language recommended).  
**Primary Dependencies (RECOMMENDATION)**: PostgreSQL-compatible DB (Amazon RDS/Aurora), Redis (cache, locks), Message queue (SQS or equivalent), Object storage (S3 with SSE-KMS), KMS for key management.  
**Storage**: Primary record store: RDS (Postgres); Long-term exports and blobs: S3 with KMS encryption.  
**Testing**: pytest (unit & integration), contract testing (Pact or similar), end-to-end tests in staging, Terratest for IaC modules.  
**Target Platform**: AWS (HIPAA-Eligible services only for PHI workloads)  
**Project Type**: Single/repo with clear service modules (recommended layout below)  
**Performance Goals**: Scheduling endpoints: 10 req/s sustained; appointment-list queries: p95 < 500ms; booking end-to-end (including notifications) < 2 minutes for heavy operations.  
**Constraints**: HIPAA & SOC 2 compliance, data retention (audit & PHI: 6y), explicit consent for PHI in notifications, use only approved HIPAA-Eligible AWS services.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- The implementation **MUST** satisfy mandatory constitution gates including Security & Compliance (HIPAA-Eligible services only), Testing & Quality (unit, contract, integration tests plus CI enforcement), Operational Readiness (monitoring, runbooks, SLOs), and IaC & Policy Validation (Terraform modules + OPA/Config checks).
- Security design **MUST** include encryption at-rest and in-transit, least privilege IAM, audit logging (CloudTrail), and a SOC 2 TSC mapping document.

## Project Structure

```text
specs/001-patient-admin/
├── plan.md              # This file
├── spec.md              # Feature spec
├── checklists/          # Validation checklists
│   └── requirements.md
├── tasks.md             # Task list (this file)
├── research.md          # Phase 0 research outputs
└── docs/                # Runbooks, operational docs, compliance mapping

repo/
├── services/
│   ├── registry-service/      # patient/person/provider/location models + APIs
│   ├── scheduling-service/    # booking engine, conflict detection
│   ├── notification-service/  # consented notifications + delivery controls
│   └── appointment-ui/        # optional lightweight web/admin UI
├── infra/                     # Terraform modules and environments
└── tests/                     # integration/contract/e2e tests
```

**Structure Decision**: Use a service-oriented layout to allow independent deployment and scaling of registry, scheduling, and notification components. Keep infra as separate Terraform modules to promote reuse and governance.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No required constitutional violations identified. Any proposed deviations from HIPAA-Eligible service usage or increased data retention must be documented and approved by Security/Compliance owner.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | N/A | N/A |

## Phase Plan & Milestones

### Phase 0: Research & Design (GATE)
- Architect data model for registries and appointment booking.
- Define SLOs, retention policy codification (6y), and incident runbooks.
- Security design: threat model for PHI flows and mapping to SOC 2 TSC/HIPAA safeguards.
- IaC design: Terraform module boundaries, state backend design, and policy rules (OPA/AWS Config). 
- Deliverables: design docs, data model, SOC 2/HIPAA control mapping, runbook drafts.

### Phase 1: Foundational (BLOCKING)
- Prepare infra baseline modules (networking, RDS, S3, KMS, monitoring) with safe defaults and tests.
- Implement CI pipelines with: linting, unit tests, contract tests, `terraform plan` + policy-as-code validation.
- Implement identity & SSO integration (IdP config, RBAC mapping, SCIM if available).
- Implement audit logging (CloudTrail/central logging) and encrypted remote state (S3 SSE + DynamoDB lock).
- Deliverables: infra modules v0.1, CI pipeline, IdP integration doc, runbooks verified in staging.

### Phase 2: Core Implementation (MVP)
- Implement Registry Service (Patient, Person, Provider, Location, Organization) including change history and audit events.
- Implement Scheduling Service (appointment lifecycle: request, book, modify, cancel) and conflict detection.
- Add resource allocation and locking (Redis/DB locks) and transaction design.
- Implement Notification Service (consent model, delivery controls, retry logic) and secure message links for PHI.
- Deliverables: functional services, automated tests, staging deploy, and integration tests passing.

### Phase 3: Integrations & Cross-Cutting (SECURE & OBSERVABLE)
- Implement monitoring, SLO dashboards, alerting, and runbooks for key incidents.
- Implement export API for FHIR R4 JSON (controlled access and logging).
- Performance testing and tuning to meet p95 and throughput goals.
- SOC 2 evidence collection tasks and security review sign-offs.

### Phase 4: Polish & Production Readiness
- Accessibility checks, localization, and UX polish.
- Backup/restore drills, disaster recovery test, and production runbook finalization.
- Final compliance audit readiness: address any gaps in SOC 2 TSC mapping and HIPAA safeguards.

## Acceptance & Release Criteria
- Unit, contract, and integration tests **pass** in CI; coverage thresholds met for critical modules.
- IaC policy checks (OPA/AWS Config) **pass** for all Terraform plans targeting staging or production.
- SLOs validated in staging: 10 req/s sustained with <1% error, p95 query latencies met.
- Security review completed (threat model, penetration testing if required) and control mapping documented.
- Runbooks and incident procedures validated by an operational readiness review.

## Risks & Mitigations
- **Risk**: Misconfiguration of cloud resources leading to PHI exposure.  **Mitigation**: Policy-as-code checks, enforced module defaults, and pre-apply manual approvals.
- **Risk**: Incomplete IdP/SSO mappings causing access errors.  **Mitigation**: Early IdP integration and end-to-end test coverage for RBAC.
- **Risk**: Performance regressions under heavy scheduling load. **Mitigation**: Capacity testing early and use of async workers/queueing.

## Implementation Strategy

1. Complete Phase 0 research outputs and security design (small, focused PRs).  
2. Implement foundational infra and CI gates (blocking).  
3. Deliver Core Implementation for P1 scenarios (registries + booking + notifications).  
4. Validate SLOs in staging, run DR/backup drills, finalize runbooks and compliance evidence.

---

**Next steps**: Convert Plan tasks into `tasks.md` (generated) and open `/speckit.tasks` to create actionable tasks, then commence Phase 0 activities.