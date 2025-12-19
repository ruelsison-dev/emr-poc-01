# s3 module (placeholder)

Purpose: Create S3 buckets with secure defaults for PHI data and exports.

Design notes:
- Default: Server-Side Encryption with KMS (SSE-KMS)
- Block public access by default
- Versioning and lifecycle rules for archival to meet retention policies
- Enforce bucket policies that limit access to approved principals and VPC endpoints

TODOs:
- Add automated policy checks (disallow public access)
- Add tests to validate encryption and lifecycle rules
