variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
  default     = "europe-west1"
}

variable "disable_service_account_management" {
  type        = bool
  description = "Whether the used service accounts should be create by this module."
}

variable "webhook_account_id" {
  type        = string
  description = "Account id of the webhook service account only used if service account management is disabled."
}

variable "image" {
  type        = string
  description = "The image to deploy"
}

variable "image_tag" {
  type        = string
  description = "The tag of the image to deploy"
}

variable "enable_debug" {
  type        = bool
  description = "Whether to enable debug mode"
}

variable "task_queue_path" {
  type        = string
  description = "The path to the task queue"
}

variable "control_plane_url" {
  type        = string
  description = "The control plane URL"
}

variable "forward_delay_seconds" {
  type        = number
  description = "The number of seconds the webhook handler delays events before forwarding them to the control plane"
}

variable "invoker_service_account" {
  type        = string
  description = "The service account to trigger the control plane"
}

variable "invoker_service_account_id" {
  type        = string
  description = "The service account ID to trigger the control plane"
}

variable "max_instance_count" {
  type        = number
  description = "The maximum number of instances to run"
}

variable "runner_labels" {
  description = "Github runner labels"
  type        = list(string)
}

variable "webhook_secret_id" {
  type        = string
  description = "The secret ID for the webhook"
}

variable "webhook_secret_version" {
  type        = string
  description = "The secret version for the webhook"
}
