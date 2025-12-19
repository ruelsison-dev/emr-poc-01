<!--
Sync Impact Report
- Version change: unspecified → 1.1.0
- Modified principles: Added/clarified Code Quality, Test-First, UX Consistency, Performance & Efficiency, Well-Architected mappings
- Added sections: Security & Compliance Requirements, Infrastructure & IaC Requirements, Observability & SLOs
- Removed sections: none
- Templates updated: ✅ .specify/templates/plan-template.md, ✅ .specify/templates/spec-template.md, ✅ .specify/templates/tasks-template.md
- Templates pending manual review: ⚠ .specify/templates/commands/*.md (verify no agent-specific names / update guidance as needed)
- Follow-up TODOs: RATIFICATION_DATE TODO(RATIFICATION_DATE): confirm original adoption date; manual review of commands/README to reflect HIPAA/SOC2 constraints
-->

# EMR POC Constitution

## Core Principles

### I. Code Quality & Maintainability (NON-NEGOTIABLE)
- Code **MUST** be readable, well-documented, and modular. Each module or service **MUST** have a clear purpose, documented API/contract, and a defined owner.
- Automated static analysis and linters **MUST** run in CI and block merges when quality gates fail. Security-sensitive code **MUST** undergo additional review.
- Dependency upgrades **MUST** be tracked and applied in a timely manner; breaking changes **MUST** follow the breaking change policy and deprecation schedule.

**Rationale**: Maintainable code reduces risk, lowers operational burden, and speeds iteration.

### II. Test-First & Continuous Validation (NON-NEGOTIABLE)
- Tests **MUST** be written prior to implementation for all new behavior: unit, integration, and contract tests as appropriate. Red-Green-Refactor cycle is the preferred workflow.
- CI **MUST** run test suites and enforce coverage thresholds for new and changed code; critical systems (including PHI handling) **MUST** have higher coverage targets and integration testing in an environment that mirrors production.
- Contract testing and consumer-driven contract verification **MUST** be used for service boundaries.

**Rationale**: Tests protect correctness, prevent regressions, and provide confidence for safe deployments.

### III. User Experience Consistency & Accessibility
- UX patterns must be consistent across products and documented in a shared design system. Changes that affect UX **MUST** be reviewed with product and design and have acceptance criteria.
- Accessibility (WCAG AA) **MUST** be considered for UI features. Error states, performance expectations, and localized content **MUST** be specified and tested.

**Rationale**: A consistent, accessible UX improves usability, reduces support burden, and broadens adoption.

### IV. Performance & Efficiency (SLO-Driven)
- Every user-facing or API surface **MUST** define performance objectives (e.g., p95, p99) and resource budgets. Performance regressions **MUST** be detected by CI and monitored in production.
- Right-sizing, caching, and efficient algorithms **SHOULD** be preferred to larger instance types. Autoscaling policies **MUST** be defined and tested.

**Rationale**: Efficient systems lower cost, meet user expectations, and scale predictably.

### V. Well-Architected Principles (Mapped to Design & Best Practices)
For all services and infrastructure designs, follow the six pillars of the AWS Well-Architected Framework. Each pillar below contains mandatory and recommended practices.

- **Operational Excellence** (MUST)
  - Runbooks, run/playbooks, deployment and recovery procedures **MUST** be documented and tested. Use IaC for reproducible environments and automation for frequent, small, reversible changes.
- **Security** (MUST)
  - Apply least privilege, MFA for privileged accounts, and role-based access control. Encrypt data in transit (TLS 1.2+) and at rest (KMS-managed keys). Enable audit logging (CloudTrail) and continuous monitoring (GuardDuty, Config) where available.
  - Use only AWS services that are designated **HIPAA-Eligible** and in-scope for SOC 2 for systems handling PHI; maintain an approved services registry and consult AWS compliance documentation when adding new services.
- **Reliability** (MUST)
  - Design for failure: health checks, multi-AZ deployments where required, backups, snapshot policies, and tested restore procedures. SLOs and incident escalation policies **MUST** be defined.
- **Performance Efficiency** (MUST)
  - Benchmark designs, implement performance budgets, and monitor utilization. Emphasize right-sizing, caching layers, and efficient data access patterns.
- **Cost Optimization** (SHOULD)
  - Tag all resources, estimate costs during design, and incorporate cost reviews into architecture decisions.
- **Sustainability** (SHOULD)
  - Favor efficient compute and storage options, avoid wasteful overprovisioning, and include lifecycle policies for data.

**Rationale**: Following Well-Architected principles ensures systems are secure, performant, resilient, and cost-effective.

### VI. Infrastructure as Code & Terraform Best Practices
- Cloud resources **MUST** be specified as code (Terraform preferred). Modules **MUST** be reusable, documented, and versioned; avoid environment-specific logic inside modules.
- Remote state **MUST** be encrypted and locked (S3 SSE + DynamoDB lock recommended). Access to state **MUST** be restricted and auditable.
- Policy-as-code (OPA/Sentinel/AWS Config rules) **MUST** run in CI to validate Terraform plans; all modules **MUST** include automated tests (e.g., Terratest) and be reviewed before release.
- IAM, KMS key usage, encryption defaults, tagging, and resource sizing **MUST** be enforced by module design and CI checks.

**Rationale**: IaC enforces consistency, enables repeatable deployments, and integrates governance earlier in the lifecycle.

## Security & Compliance Requirements
- Apply SOC 2 TSC-aligned controls for Security, Availability, Processing Integrity, Confidentiality, and Privacy. Document how each implemented control maps to TSC criteria in specs.
- For HIPAA: all PHI **MUST** be protected with strict access controls, encryption at-rest and in-transit, minimum necessary access, audit logging of accesses and modifications, and documented incident response and breach notification procedures.
- Maintain a compliance evidence collection plan (who collects what, where retained, retention periods) and schedule regular compliance reviews and audits.
- All AWS services used for regulated workloads **MUST** be in the project's approved HIPAA/SOC 2 service registry.

## Development Workflow & Quality Gates
- Pull requests **MUST** include necessary tests and satisfy CI gates (linting, formatting, unit/integration tests, IaC policy checks). Security-sensitive changes **MUST** include a security review and threat model if they touch sensitive data.
- Release process **MUST** include an automated plan/apply workflow for IaC with manual approvals for production, rollback strategies, and post-deploy verification steps.
- Monitoring, alerting, SLOs, and runbooks **MUST** be present before production release.

## Governance
- The Constitution supersedes local/adhoc practices for matters it governs. Any change to the Constitution **MUST** be made via an explicit amendment PR with a rationale, migration plan, and tests where applicable.
- **Amendment Procedure**:
  1. Create a draft amendment PR against `.specify/memory/constitution.md` with rationale and example impact.
  2. Obtain approvals: Architecture Owner and Security/Compliance Owner (both required) plus one additional maintainer.
  3. Run impact analysis and update affected templates and guidance files.
  4. Merge and record the **Last Amended** date and bump the constitution version as follows:
     - **MAJOR**: Removal or incompatible redefinition of a principle or governance change that breaks existing policies.
     - **MINOR**: Addition of a principle or material expansion of guidance (this change: +1 minor).
     - **PATCH**: Clarifications, typos, or non-semantic wording changes.
- **Enforcement & Validation**: CI checks, policy-as-code, and automated test suites **MUST** be used to enforce requirements; periodic compliance scans and quarterly reviews **MUST** be scheduled.

**Version**: 1.1.0 | **Ratified**: TODO(RATIFICATION_DATE) | **Last Amended**: 2025-12-19

