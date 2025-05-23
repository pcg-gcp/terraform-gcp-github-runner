variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
}

variable "allowed_zones" {
  type        = list(string)
  description = <<-EOT
  The zones to deploy the runner instances into.
  If not set the runner will be deployed into any zone in the region.
  Also see use_strict_zone_order
  EOT
}

variable "use_strict_zone_order" {
  type        = bool
  description = <<-EOT
  If this is set to true and allowed_zones is set the runner will always be deployed in the first available zone in the list unless it is unavailable.
  If allowed_zones is not set the first zone returned by the API will be used.
  Otherwise the runner will be deployed in a random zone either from the allowed_zones list or from the API.
  EOT
}

variable "disable_service_account_management" {
  type        = bool
  description = "Whether the used service accounts should be create by this module."
}

variable "control_plane_account_id" {
  type        = string
  description = "Account id of the control plane service account only used if service account management is disabled."
}

variable "invoker_account_id" {
  type        = string
  description = "Account id of the invoker service account only used if service account management is disabled."
}

variable "enable_debug" {
  type        = bool
  description = "Whether to enable debug mode"
}

variable "runner_service_account_id" {
  type        = string
  description = "The service account to run the runner as"
}

variable "runner_labels" {
  description = "Github runner labels"
  type        = list(string)
}

variable "instance_template_name" {
  type        = string
  description = "The name of the instance template to use for runner instances"
}

variable "image" {
  type        = string
  description = "The image to deploy"
}

variable "image_tag" {
  type        = string
  description = "The tag of the image to deploy"
}

variable "max_instance_count" {
  type        = number
  description = "The maximum number of instances to run"
}

variable "github_app_id" {
  type        = string
  description = "The GitHub App ID"
}

variable "private_key_secret_id" {
  type        = string
  description = "The secret ID of the private key"
}

variable "private_key_secret_version" {
  type        = string
  description = "The secret version of the private key"
}

variable "shutdown_schedule" {
  type        = string
  description = "The shutdown schedule in cron format"
}

variable "shutdown_schedule_timezone" {
  type        = string
  description = "The timezone of the shutdown schedule"
}

variable "shutdown_attempt_timeout" {
  type        = string
  description = "The timeout for the shutdown attempt"
}

variable "ephemeral" {
  type        = bool
  description = "Whether to use ephemeral runners"
}

variable "use_jit_config" {
  type        = bool
  description = "Whether to use JIT config"
  default     = false
}

variable "use_org_runners" {
  type        = bool
  description = "Whether to use github organization runners"
  default     = false
}

variable "max_runner_count" {
  type        = number
  description = "The maximum number of runners that should be deployed at the same time"
}

variable "min_runner_count" {
  type        = number
  description = "The minimum number of runners that should be deployed at all times"
}
