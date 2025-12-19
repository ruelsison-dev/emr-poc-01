# service-module (template module)

Purpose: Common patterns used by all service modules (e.g., registry-service, scheduling-service): IAM roles, service-specific security groups, CloudWatch log groups, and optional autoscaling policies.

Safe defaults:
- Enforce least-privilege IAM roles scoped to service resources
- Enforce encryption-in-transit and at-rest for service communications and storage
- Require tags: `owner`, `environment`, `cost-center`, `service`
- Do not create public IPs by default; require explicit opt-in for public-facing resources

Inputs/Outputs:
- Inputs: `service_name`, `environment`, `tags`, `vpc_id`, `subnet_ids`
- Outputs: `iam_role_arn`, `security_group_id`, `cw_log_group`

TODOs:
- Implement module and add tests
- Add policy-as-code coverage specific to service-level constraints
