# Phase 0 Research: Patient Administration Sub-System

**Feature**: Patient Administration Sub-System  
**Branch**: `001-patient-admin`  
**Date**: 2025-12-20

This document resolves outstanding technical clarifications and records decisions, rationale, and alternatives considered for Phase 0 design.

---

## Decision: Language & Framework
- **Previous Decision**: Use **Python 3.11** and **FastAPI** for backend services (registry-service, scheduling-service, notification-service). Rationale: FastAPI provides async support, built-in request validation (Pydantic), fast developer experience, strong test ecosystem (pytest), and good performance for IO-bound APIs; Python is widely used by data/AI teams.
- **Updated Decision (user request)**: Implement the **front-end web application and REST APIs using Node.js (20.x LTS) and ExpressJS with TypeScript**. This will be the initial implementation for the Patient Administration feature (Patient, Person, Location, Organization, Encounter, Appointment, Schedule, Slot, Questionnaire APIs).
- **Rationale for update**: The user requested ExpressJS; a TypeScript + Express approach aligns well with front-end teams, enables end-to-end JS development, and leverages mature Node.js tooling (Jest, ESLint, Pact). Express is simple and widely adopted for REST APIs; TypeScript improves maintainability and type safety across the stack.
- **Alternatives Considered**: Keep FastAPI as the primary backend and implement an Express-based frontend/API gateway that proxies to Python services (adds integration complexity), or use NestJS (structure and DI) instead of Express for stronger architectural patterns. The chosen approach favors direct development velocity and aligns with the user's request but introduces a runtime change that requires Architecture/Compliance sign-off.
- **Action Items**: Document runtime change in the plan (Complexity Tracking) and obtain Architecture Owner approval; update CI IaC pipelines to include Node.js build/test stages; update contracts (OpenAPI) to be API-first and share between consumer/producer teams using Pact.

---

## Decision: Primary Data Stores
- **Decision**: Use **Amazon Aurora (PostgreSQL-compatible)** for primary relational registries; use **Amazon S3 (SSE-KMS)** for exports and large blobs; use **ElastiCache (Redis)** for caching and distributed locking; **Amazon SQS** for durable asynchronous work queues.
- **Rationale**: Registries (Patients, Providers, Locations) require ACID guarantees and relational queries; Postgres is well-understood and Aurora offers managed HA, backups, and performance. Redis provides efficient locks and fast cache. S3 is secure for exports and complies with encryption and lifecycle needs. SQS provides a simple, reliable decoupling mechanism.
- **Alternatives Considered**: DynamoDB for extreme scale (high throughput, but modeling joins and transaction semantics is more complex). Amazon MQ (ActiveMQ) for more advanced broker semantics (but added operational complexity).
- **Action Item**: Confirm with Security whether **SQS** and preferred notification channels (SES, SNS) are acceptable / HIPAA-Eligible for PHI workloads; maintain an approved services registry.

---

## Decision: Notification Channels & PHI Delivery
- **Decision**: Prefer **in-app notifications** and secure links for any PHI; allow limited PHI (patient name + appointment time) in external messages only with **explicit patient consent** and using **SES with verified sender addresses** or a HIPAA-compliant third-party provider with a BAA.
- **Rationale**: Minimizes PHI exposure in transit while supporting UX needs; SES can be configured with secure sending and verification. External SMS/email introduces elevated risk; requiring consent and delivery controls reduces that risk.
- **Alternatives Considered**: Always use secure link only (best privacy, potentially worse UX); allow more PHI with strict encryption and per-recipient encryption policies (complex and potentially brittle).

---

## Decision: Identity & Access Management
- **Decision**: Integrate with the organization's enterprise IdP (Okta or Microsoft Entra) using **SAML/OIDC** for SSO and **SCIM** for provisioning where available; implement RBAC within services and map IdP groups/claims to roles.
- **Rationale**: Enterprise IdP enables centralized access control, MFA enforcement, and auditability. SCIM simplifies user provisioning and role membership lifecycle.
- **Alternatives Considered**: Managed Cognito (works but less enterprise-centred for SSO/SCIM) — Cognito remains an option for certain service integrations.
- **Action Item**: Document exact IdP, SCIM mappings and claim-to-role mapping as part of Phase 1 IdP integration task.

---

## Decision: Testing & Contracting
- **Decision**: Use **pytest** for unit/integration tests, **Pact** or similar consumer-driven contract testing for services, **Terratest** for IaC integration tests, **Conftest/OPA** for policy-as-code checks, and **k6** for load testing.
- **Rationale**: These tools integrate well with the chosen stack and support CI-driven validation across layers.
- **Alternatives Considered**: Other contract frameworks (e.g., Spring Contract); alternative load tools (Locust). Chosen tools are pragmatic and well-supported.

---

## Decision: IaC & Policy-as-Code
- **Decision**: Use **Terraform** with a module-first approach. Enforce remote state (S3 + DynamoDB locks), require PR-level `terraform plan` with OPA/Conftest checks in CI, and include Terratest-based tests for critical modules.
- **Rationale**: Terraform is widely used and supports modular patterns, versioning, and mature community tooling. Policy-as-code provides automated validation of security and compliance controls.
- **Alternatives Considered**: CloudFormation/Terraform CDK (higher learning curve); chosen approach balances governance and dev ergonomics.

---

## Decision: Monitoring, Logging & Security Services
- **Decision**: Use **CloudWatch** for metrics and dashboards, **CloudTrail** for audit events (centralized and immutable), and enable **GuardDuty** and **AWS Config** rules for continuous security posture checks. Store logs centrally with retention aligned to audit policy (6 years as per spec).
- **Rationale**: AWS native tooling integrates with other AWS services and supports compliance auditing requirements.
- **Alternatives Considered**: Splunk or Elastic Stack (more flexibility but added operational weight and potentially costlier for HIPAA workflows); can be added later if needed.

---

## Decision: HIPAA & SOC 2 Service Eligibility
- **Decision**: Enforce a policy that **only AWS services that are HIPAA-Eligible and in-scope for SOC 2** (per project registry) will be used to process PHI. If using a service that is not in the registry, it **MUST** be justified and approved by Security/Compliance.
- **Rationale**: Ensures the project maintains HIPAA and SOC 2 compliance boundaries and aligns with the Constitution.
- **Action Item**: Create `docs/approved-aws-services.md` listing approved services and BAAs; update during security reviews.

---

## Decision: Performance Architecture
- **Decision**: Design booking as a hybrid sync/async flow: synchronous API for request/acknowledgment and asynchronous background processing for heavy conflict resolution and resource allocation. Autoscale services according to metrics and use queue depth for backpressure.
- **Rationale**: Maintains responsive UX for clinicians while enabling durable processing for complex scheduling operations.
- **Alternatives Considered**: Fully synchronous booking (simpler) vs fully async (more robust but worse initial UX). Hybrid approach balances UX and reliability.

---

## Decision: Export Formats
- **Decision**: Export patient longitudinal records as **FHIR R4 JSON** (controlled via role-based export permissions; exports logged and available for compliance review). Provide chunked exports for large datasets.
- **Rationale**: FHIR R4 is a common, interoperable format for EMR exchange; chunked exports keep memory and transfer bounded.

---

## Outcome & Next Steps
- **Outcome**: All previously flagged NEEDS CLARIFICATION items in `plan.md` have recommended decisions recorded here. These choices resolve the immediate design unknowns so Phase 1 (data model, contracts, IaC module design) can proceed.

**Next Actions**:
- Add `research.md` to the feature spec directory (done).
- Proceed to Phase 1: create `data-model.md`, API contracts (`/contracts/`), `quickstart.md` and update agent context if desired.
- Create `docs/approved-aws-services.md` as part of Phase 1 gating with Security.

---

## Outstanding research tasks (Phase 0)
- Confirm whether ExpressJS will fully replace FastAPI for the Patient Administration APIs or coexist as an API gateway. (owner: Architecture)
- Select Node.js runtime version (recommend Node 20.x LTS) and confirm TypeScript adoption (recommend: TypeScript required). (owner: Eng Lead)
- Select hosting model for the Express service (recommend ECS Fargate with ALB); evaluate App Runner and Lambda for specific use-cases. (owner: Infra)
- Produce ExpressJS security hardening checklist for PHI workloads and add CI enforcement items (Snyk, DAST configuration, dependency pinning). (owner: Security)
- Confirm migration strategy or integration pattern if existing FastAPI services are retained (owner: Architecture/Engineering)

