# network module (placeholder)

Purpose: Provide VPC, subnets, routing and optional NAT gateways. This module MUST enforce private subnets for PHI-handling workloads and restrict inbound public access by default.

Inputs: (documented in variables.tf)
Outputs: (documented in outputs.tf)

TODOs:
- Define variables for CIDR ranges, subnet counts, and AZs
- Implement NAT and route table defaults
- Add automated tests (Terratest) and policy-as-code validations
