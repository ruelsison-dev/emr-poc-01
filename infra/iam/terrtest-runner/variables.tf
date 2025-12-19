variable "test_run_id" {
  type        = string
  description = "Unique id for namespacing test resources (recommended to use CI run id)"
  default     = "local"
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}
