Autoscaling module

This module provides constructs for autoscaling compute based on service type (ECS, EC2 ASG, etc.).

Usage: populate this module with your service-specific resources and adust the inputs to target metrics and scaling behavior.

Security: ensure IAM roles used for scaling are least-privilege and that scaling rules do not expose sensitive endpoints.