# Feature Specification: [FEATURE NAME]

**Feature Branch**: `[###-feature-name]`  
**Created**: [DATE]  
**Status**: Draft  
**Input**: User description: "$ARGUMENTS"

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - [Brief Title] (Priority: P1)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently - e.g., "Can be fully tested by [specific action] and delivers [specific value]"]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]
2. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 2 - [Brief Title] (Priority: P2)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 3 - [Brief Title] (Priority: P3)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- What happens when [boundary condition]?
- How does system handle [error scenario]?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST [specific capability, e.g., "allow users to create accounts"]
- **FR-002**: System MUST [specific capability, e.g., "validate email addresses"]  
- **FR-003**: Users MUST be able to [key interaction, e.g., "reset their password"]
- **FR-004**: System MUST [data requirement, e.g., "persist user preferences"]
- **FR-005**: System MUST [behavior, e.g., "log all security events"]

### Security & Compliance Requirements (MANDATORY when applicable)

- **SC-SEC-001**: If the feature processes PHI or regulated data, all designs **MUST** use only AWS services that are designated HIPAA-Eligible and in-scope for SOC 2; any service outside the approved registry **MUST** be justified and approved by Security.
- **SC-SEC-002**: All data classified as PHI **MUST** be encrypted at-rest (KMS-managed keys preferred) and in-transit (TLS 1.2+).
- **SC-SEC-003**: Access to PHI **MUST** follow least-privilege principles, be gated by IAM roles and MFA for privileged accounts, and be auditable (CloudTrail or equivalent).
- **SC-SEC-004**: Audit logging, monitoring, and alerting **MUST** capture access to sensitive data and be retained per retention policy for compliance review.
- **SC-SEC-005**: A mapping of implemented controls to SOC 2 Trust Services Criteria (Security, Availability, Processing Integrity, Confidentiality, Privacy) **MUST** be included in the spec when applicable.

### Infrastructure & IaC Requirements (MANDATORY when infra change required)

- **IAAC-001**: Infrastructure **MUST** be defined as Terraform modules with clear inputs, outputs, and semantic versioning.
- **IAAC-002**: Remote state **MUST** be stored securely (S3 with SSE + DynamoDB locking) and access to state **MUST** be controlled and audited.
- **IAAC-003**: Terraform modules **MUST** include automated tests (e.g., Terratest), and CI **MUST** run `terraform plan` + policy-as-code checks (OPA/Sentinel/AWS Config rules) before apply.
- **IAAC-004**: Resource sizing, tagging, and cost estimation **MUST** be documented.

### Observability, Performance & UX

- **OBS-001**: Monitoring metrics, dashboards, and alert thresholds **MUST** be specified; SLOs and objectives **MUST** be stated.
- **PERF-001**: Performance goals (p95/p99 latencies, throughput) and test plans **MUST** be included in the spec.
- **UX-001**: Accessibility (WCAG AA) and UX consistency guidelines **MUST** be followed for UI work; error states and localized messages **MUST** be specified.

*Example of marking unclear requirements:*

- **FR-006**: System MUST authenticate users via [NEEDS CLARIFICATION: auth method not specified - email/password, SSO, OAuth?]
- **FR-007**: System MUST retain user data for [NEEDS CLARIFICATION: retention period not specified]

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: [Measurable metric, e.g., "Users can complete account creation in under 2 minutes"]
- **SC-002**: [Measurable metric, e.g., "System handles 1000 concurrent users without degradation"]
- **SC-003**: [User satisfaction metric, e.g., "90% of users successfully complete primary task on first attempt"]
- **SC-004**: [Business metric, e.g., "Reduce support tickets related to [X] by 50%"]
