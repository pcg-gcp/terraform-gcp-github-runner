variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
}

variable "zone" {
  type        = string
  description = "The zone to deploy runner instances into"
}

variable "runner_service_account_id" {
  type        = string
  description = "The service account to run the runner as"
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
