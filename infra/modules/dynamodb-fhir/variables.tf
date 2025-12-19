variable "table_name" {
  type        = string
  description = "DynamoDB table name"
}

variable "billing_mode" {
  type        = string
  default     = "PROVISIONED"
  description = "Billing mode (PROVISIONED or PAY_PER_REQUEST)"
}

variable "hash_key" {
  type        = string
  default     = "id"
}

variable "gsis" {
  type        = any
  default     = []
}

variable "tags" {
  type    = map(string)
  default = {}
}
