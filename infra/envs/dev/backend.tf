# Backend example for dev environment (placeholder)
# IMPORTANT: For production, use a secure production backend with proper encryption and access controls.

# Example (commented):
# terraform {
#   backend "s3" {
#     bucket         = "<org>-terraform-state-dev"
#     key            = "001-patient-admin/terraform.tfstate"
#     region         = "us-east-1"
#     encrypt        = true
#   }
# }

# Note: Ensure remote state access is restricted and DynamoDB lock table is used for locking.
