# Feature Specification: Patient Administration Sub-System

**Feature Branch**: `001-patient-admin`  
**Created**: 2025-12-20  
**Status**: Draft  
**Input**: Develop an AI-Native EMR: Patient Administration Sub-System — add/modify registries (patients, persons, service delivery locations, providers, places, organizations); implement scheduling for appointments (requesting, booking, notifications, modifications), resource allocation, and queries. Must comply with HIPAA and SOC 2 (TSC mapping); PHI handling policies: encrypt in transit & at rest, audit logging, least privilege; use only HIPAA-Eligible AWS services; IaC via Terraform modules with policy-as-code checks; include SLOs, monitoring, and runbooks.

## Clarifications

### Session 2025-12-20

- Q: What should be the canonical patient identifier used across APIs and data models? → A: Use opaque system-generated UUIDs as the primary canonical identifier; store MRN and other external identifiers as mapped secondary identifiers.
- Q: How should the system resolve concurrent booking conflicts (same resource/time slot)? → A: Use reservation tokens (short-term holds with configurable TTL) issued at booking initiation, with background reconciliation and idempotent retry semantics to deterministically resolve conflicts.
- Q: How should patient consent for notifications be represented and stored? → A: Use per-purpose + per-channel consent records with timestamp and revocation history (audit-friendly); enforce consent checks at send time and retain consent audit trails per retention policy.
- Q: Should consent revocation trigger immediate notifications to stakeholders (patient and/or relevant staff)? → A: Yes — notify both the patient and relevant clinicians/staff immediately and immediately suppress pending notifications for the revoked purpose/channel; log suppression and notification events as audit records.
- Q: Who should be notified when consent is revoked? → A: Notify the patient, the patient's active care team (assigned clinicians and schedulers), and a configurable organizational escalation list (e.g., Privacy Officer); make the escalation list configurable per-organization.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Clinician schedules an appointment (Priority: P1)
A clinician schedules an appointment for a patient with a provider at a service delivery location and assigns required resources.

**Why this priority**: Core patient administration workflow required for care delivery.

**Independent Test**: Using a clinician account, create/modify/cancel an appointment via UI and API; verify resource assignment, audit entry, and notifications.

**Acceptance Scenarios**:
1. **Given** a clinician with appropriate role and a valid patient record, **When** they create an appointment with date/time and resource requirements, **Then** the system records the appointment, reserves assigned resources, persists an audit event, and returns a confirmation within the performance goal (SC-APPT-001).
2. **Given** an existing appointment, **When** the clinician modifies or cancels it, **Then** changes are persisted, conflicts are checked, audit events recorded, notifications issued, and resource allocations updated.

---

### User Story 2 - Scheduler or System processes booking requests (Priority: P1)
Schedulers (or automated workflows) receive appointment requests, accept/book as available, and handle conflicts and rescheduling.

**Independent Test**: Submit a booking request and verify automated or manual acceptance flows and conflict resolution.

**Acceptance Scenarios**:
1. **Given** a request for a slot where the provider is unavailable, **When** the system tries to auto-schedule, **Then** it returns recommended alternate slots and logs the decision.

---

### User Story 3 - Notification & Patient/Staff Query (Priority: P2)
Send appointment confirmations/reminders to patients and staff using minimal PHI and allow secure in-app notifications; enable clinicians to query appointment and registry data.

**Independent Test**: Trigger a notification and verify minimal content is used (no PHI leaked in clear channels), and perform common queries returning results within performance targets.

**Acceptance Scenarios**:
1. **Given** a booked appointment, **When** the system sends a reminder, **Then** the message may include the patient's name and appointment time (subject to patient consent and preference), or a secure link; all outbound messages containing PHI **MUST** use delivery controls (e.g., verified recipient addresses or encryption), be logged with consent metadata, and trigger an audit event.
2. **Given** a clinician query for a patient's appointments, **When** the query is executed, **Then** results return within response target (see Success Criteria) and are filtered by clinician's authorization.

---

### Edge Cases
- Concurrent booking attempts for the same resource/time slot — system MUST issue reservation tokens (short-term holds with a configurable default TTL, e.g., 120 seconds) when a user begins a booking; holds expire automatically, and background reconciliation processes must re-evaluate pending requests and resolve conflicts deterministically (e.g., earliest hold wins or prioritized queue). Booking flows MUST use idempotency keys and provide clear UI/UX guidance when a hold cannot be secured or if the booking failed due to a conflict.
- Out-of-window modification requests (e.g., recurring appointments) — system must enforce business rules and notify affected parties.
- Failed notification delivery — system must retry with exponential backoff and surface alerts if delivery fails repeatedly.

---

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST allow creation, modification, cancellation, and querying of appointments for registered patients.
- **FR-002**: System MUST maintain registries for Patient, Person, Provider, Location, Place, Organization with unique identifiers (system-generated UUIDs as primary identifiers) and change history.
- **FR-003**: System MUST support resource allocation for appointments (rooms, equipment, staff), detect conflicts, and implement reservation tokens (short-term holds with configurable TTL) to temporarily reserve resources during booking; background reconciliation and idempotent booking retries MUST resolve remaining conflicts deterministically and surface clear guidance to the actor.
- **FR-004**: System MUST record audit events for all operations that create/modify/view PHI and make them available for compliance review.
- **FR-005**: System MUST provide search/query endpoints for authorized roles to retrieve registry and scheduling information.
- **FR-006**: System MUST provide notifications workflows (email/in-app/SMS proxies). Outbound messages **MAY** include limited PHI (patient name and appointment time) only when the patient has given explicit consent and preferences; otherwise only minimal identifiers or a secure link must be used. Consent **MUST** be represented as per-purpose and per-channel consent records (purpose, channel, granted_at, revoked_at, source, and metadata) with revocation history and enforceable at send time. Consent revocation **MUST** immediately suppress any pending outbound notifications for the affected purpose/channel and generate immediate notifications to the patient, the patient's active care team (assigned clinicians and schedulers), and a configurable organizational escalation list (e.g., Privacy Officer); the escalation list **MUST** be configurable per-organization. Suppression and notification events **MUST** be logged as audit records. Messages that include PHI **MUST** use delivery controls (encrypted channels or verified recipient addresses), be logged with consent metadata, and include retry and failure handling; consent audit trails **MUST** be retained and queryable for compliance.
- **FR-007**: System MUST enforce role-based access control so only authorized roles (Clinician, Scheduler, Admin) can perform sensitive actions.
- **FR-008**: System MUST validate inputs and enforce data integrity and idempotency for booking operations.

### Security & Compliance Requirements (MANDATORY)
- **SC-SEC-001**: For any PHI processing, designs **MUST** use only AWS services that are HIPAA-Eligible and in-scope for SOC 2 per the project's approved registry.
- **SC-SEC-002**: PHI **MUST** be encrypted in transit (TLS 1.2+) and at rest (KMS-managed keys); keys **MUST** be rotated and access audited.
- **SC-SEC-003**: Access to PHI **MUST** be least-privilege, enforced by IAM roles and RBAC in the application layer; privileged actions **MUST** require MFA and justification logging.
- **SC-SEC-004**: Audit logs capturing read/write access to PHI **MUST** be immutable, centrally collected, and retained for **6 years** for audit logs and **6 years** for PHI records (imaging **6 years**); consent records and consent audit trails **MUST** also be retained and queryable for **6 years**. Retention schedules must be documented, enforced, and periodically reviewed.
- **SC-SEC-005**: Map implemented controls to SOC 2 Trust Services Criteria (Security, Availability, Processing Integrity, Confidentiality, Privacy) and include mapping in the spec.

### Infrastructure & IaC Requirements
- **IAAC-001**: All infrastructure **MUST** be defined as Terraform modules with clear inputs/outputs, semantic versioning, and no environment-specific logic inside modules.
- **IAAC-002**: Remote state **MUST** be stored securely (e.g., S3 with SSE + DynamoDB locking) and access **MUST** be controlled and audited.
- **IAAC-003**: CI **MUST** run `terraform plan` with policy-as-code checks (e.g., OPA) and automated IaC tests (e.g., Terratest) before any apply.
- **IAAC-004**: Modules **MUST** enforce safe defaults (encryption, tagging, minimal public network exposure) and be covered by unit/regression tests.

### Observability, Performance & UX
- **OBS-001**: Define monitoring metrics, dashboards, alerts, and runbooks for operational incidents; instrument booking throughput, queue depth, error rates, and notification delivery success.
- **PERF-001**: Performance goals: support **10 req/s** for scheduling endpoints and return query results to clinicians within **500ms p95** in normal load; appointment booking flows (end-to-end) **MUST** complete within **2 minutes** for large batch operations.
- **UX-001**: UI components must meet **WCAG AA** accessibility and follow shared design system patterns.

### Testing & QA
- **TST-001**: Unit tests, integration tests, and contract tests **MUST** be provided for all new behavior; tests **MUST** run in CI and fail builds on regressions.
- **TST-002**: End-to-end tests in a staging environment that mirrors production **MUST** validate SLOs and recoverability scenarios (e.g., partial outage of a downstream service).
- **TST-003**: Security testing (SAST, dependency scanning, secrets scanning) **MUST** be part of CI and DAST prior to production deployment.

## Key Entities
- **Patient**: Primary identifier: opaque system-generated UUID (used as canonical primary id in APIs and DB PKs); demographic attributes; links to Person entities and consent metadata. MRN and other external identifiers are stored as secondary mapped identifiers (type/value) for interoperability and lookups.
- **Person**: Human entity (may be patient, provider, or contact); used for contact and identity.
- **Provider**: Healthcare provider with role(s), schedules, credentials, and availability.
- **Consent**: Per-purpose and per-channel consent records with fields: purpose (e.g., appointment_reminder), channel (email/SMS/in-app), granted_at (timestamp), revoked_at (timestamp, nullable), source (UI/API/SCIM), granted_by (user or system), and audit metadata. Consent records and revocation history **MUST** be immutable, auditable, and retained per retention policy (refer to SC-SEC-004 for retention schedule). Consent revocation **MUST** trigger immediate suppression of pending notifications for that purpose/channel and generate immediate notification events to the patient, the patient's active care team, and any configured escalation list; the escalation list **MUST** be configurable per organization, and these actions **MUST** be logged for compliance.
- **Location / Place / ServiceDeliveryLocation**: Physical place or virtual service location where services are delivered.
- **Organization**: Health organization or practice that owns providers and locations.
- **Appointment / Booking**: Scheduled service with participants, times, resources, status, and audit history.
- **Resource**: Rooms, equipment, or staff allocations required by Appointment.

## Success Criteria *(mandatory)*
- **SC-001**: Clinicians can create, modify or cancel an appointment with confirmation returned and audit logged in under **2 minutes** end-to-end for normal workloads.
- **SC-002**: Scheduling endpoints handle **10 req/s** sustained with less than **1%** error rate under load test conditions.
- **SC-003**: **95%** of queries for appointment listings return within **500ms** p95.
- **SC-004**: Security: All accesses to PHI are logged and **100%** of modifications/audit events are retained for **6 years** per the defined retention policy.
- **SC-005**: Compliance: Provide a mapping document showing how implemented controls satisfy SOC 2 TSC and HIPAA safeguards for PHI handling.

## Assumptions
- The system will operate primarily within jurisdictions covered by organization's current HIPAA BAA; data residency and legal retention specifics will follow local regulation unless otherwise specified.
- Notifications will avoid embedding PHI in cleartext channels; where necessary, they will include secure links requiring authentication.
- Default retention for audit logs and PHI-related artifacts is set to **6 years** (audit logs: 6y; PHI records: 6y; imaging: 6y) as the organization default, unless compliance owners specify otherwise.

## Out of Scope
- Billing and claims processing workflows.
- Full FHIR mediation and transformation beyond exporting records; the initial scope focuses on scheduling, registries, and appointment-related data with standardized exports as needed.

## Operational Considerations
- Runbooks for operational incidents, service degradation, and data breaches **MUST** be documented and tested quarterly.
- Consent revocation handling **MUST** be part of runbook procedures: revocations that affect active workflows (e.g., pending notifications) **MUST** trigger immediate suppression and generate alerts to on-call or duty staff; runbooks **MUST** include steps to validate suppression, escalate if suppression fails, and document remediation.
- Backup and restore procedures for registries **MUST** be executed and validated with backup restore drills.

## Dependencies
- Approved list of HIPAA-Eligible AWS services (project registry) and BAAs in place.
- Identity provider for enterprise SSO and RBAC integration.

---

**Next Steps**: Create implementation plan (`/speckit.plan`) and run the spec quality checklist to validate completeness, then schedule clarifications for the highlighted items.

