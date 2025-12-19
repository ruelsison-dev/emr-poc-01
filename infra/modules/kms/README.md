# kms module (placeholder)

Purpose: Create and manage KMS keys with rotation and minimal key policy footprint.

Design notes:
- Keys are multi-purpose but prefer separate keys for PHI vs non-PHI data when audit separation is required.
- Enforce key rotation and strict IAM key usage policies.
- Document key aliasing strategy and grant practices.

TODOs:
- Define key creation, rotation, and grant management.
- Add checks to ensure keys are not accidentally made asymmetric unless required.
