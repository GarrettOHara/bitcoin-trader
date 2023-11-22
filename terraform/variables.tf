variable "name" {
  description = "The prefix for the name of the resources."
  type        = string
  default     = "bitcoin-trader"
}

variable "aws_region" {
  description = "The aws region the resources will be deployed to."
  type        = string
  default     = "us-west-1"
}

variable "flexible_time_window" {
  description = "Determines if the EventBridge Schedule will have a flexable time window."
  type        = bool
  default     = "false"
}
