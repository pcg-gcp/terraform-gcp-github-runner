variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
  default     = "europe-west1"
}

variable "image" {
  type        = string
  description = "The image to deploy"
}

variable "task_queue_path" {
  type        = string
  description = "The path to the task queue"
}

variable "control_plane_url" {
  type        = string
  description = "The control plane URL"
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

variable "webhook_secret_id" {
  type        = string
  description = "The secret ID for the webhook"
}

variable "webhook_secret_version" {
  type        = string
  description = "The secret version for the webhook"
}
