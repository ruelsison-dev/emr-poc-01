# MFA & Justification for Privileged Actions

Purpose: Define the design and implementation plan for enforcing multi-factor authentication (MFA) and collecting a short, auditable justification when users perform privileged actions (e.g., modifying high-risk PHI, escalating privileges, performing destructive infra operations).

Scope
- Application-level privileged actions (e.g., role/permission changes, bulk deletions, PHI exports, overriding consent suppression) and privileged CI/IaC actions (e.g., running Terratest in a shared account, applying infra changes to production).

Design Goals
- Enforce MFA for privileged actions: require a recent MFA validation tied to the user session or a short-lived step-up flow.
- Capture a justification string for the privileged action (min length, limited character set) and log metadata (actor id, timestamp, action, resource id, auth method, client ip, user agent).
- Store justification and evidence in an immutable audit store with retention per policy (6 years).
- Provide integration patterns for IdP-based MFA (Okta/MS Entra) or step-up using OAuth/OIDC where necessary.
- Include tests (unit + integration) that validate enforcement, logging, and audit retention.

Implementation Plan

1. IdP Integration (Auth):
   - Recommend using enterprise IdP (Okta/MS Entra) to perform MFA step-up via OIDC claims (e.g., `amr` claim or custom claim indicating `mfa`), or use re-auth endpoints where IdP supports step-up.
   - For CI flows (e.g., Terratest runs), require temporary role assumption with MFA gating or require Security approval via protected environments (manual step) before enabling runs.

2. API Enforcement (App middleware):
   - Add middleware `mfaEnforce` which checks:
     - `req.user.mfa_verified` OR a recent `mfa_verified_at` timestamp within X minutes, OR
     - Presence of a `X-MFA-Verified: true` header for synthetic tests (only for tests and local dev — must be disabled in production)
   - Require `X-Justification` header or `justification` field in request body for privileged endpoints.
   - Validate justification (min length 15 chars, max 1024 chars) and sanitize inputs.

3. Justification Logging & Audit Storage:
   - Create helper `lib/justification` to create a justification record referencing: `id, actor_id, action, resource_id, justification, mfa_verified, timestamp, metadata`.
   - Persist to central audit logging (CloudWatch/aggregator) and also write a compliance-evidence artifact (e.g., a signed JSON file in S3 or in `specs/001-patient-admin/docs/compliance_evidence/mfa/` during tests).

4. Tests & Evidence:
   - Add integration tests that prove: unauthenticated or non-MFA requests to privileged endpoints are rejected; authenticated + MFA + justification requests succeed and produce a justification audit entry.
   - Capture logs and a sample evidence file that can be attached to PRs for compliance review.

5. Acceptance Criteria:
   - Privileged endpoints require MFA and justification to perform state-changing actions.
   - Justifications are recorded and are immutable for the retention period.
   - Security & Architecture have approved the IdP flow and evidence capture approach.

Security Notes
- Never accept justification via email or other unverified channels; only capture via secure, authenticated channels.
- Sanitize justification fields to prevent log injection and do not store PII in free-form justification unless absolutely necessary.
- Record evidence of MFA (timestamp, method) with each justification record.

Next Steps
- Implement middleware & logging helper with integration tests (T050b–T050d).
- Iterate on IdP integration with Security to select exact claims and step-up flow.

