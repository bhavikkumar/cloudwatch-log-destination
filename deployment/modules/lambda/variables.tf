variable "kms_key_arn" {
  type        = "string"
  description = "The KMS Key to use for encrypting environment variables"
}

variable "lambda_version" {
  type        = "string"
  description = "The version of the lambda function to deploy"
}

variable "log_destination_arn" {
  description = "The destination which log groups will be subscribed to"
}

variable "log_retention_period" {
  description = "The number of days to retain the logs for in CloudWatch"
  default     = 14
}

variable "tags" {
  type        = "map"
  description = "A map of tags to add to all resources"
  default     = {}
}
