variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
  default     = "europe-west1"
}

variable "zones" {
  type        = list(string)
  description = "The zones to deploy the runner instances into"
  default     = ["europe-west1-b", "europe-west1-c", "europe-west1-d"]
}


variable "zone" {
  type        = string
  description = "The zone to deploy runner instances into"
  default     = "europe-west1-b"
}

variable "vpc_name" {
  type        = string
  description = "The VPC to deploy runner instances into"
}

variable "subnet_name" {
  type        = string
  description = "The subnet to deploy runner instances into"
}

variable "runner_image_path" {
  type        = string
  description = "The image to deploy"
}

variable "runner_machine_type" {
  type        = string
  description = "The machine type to deploy"
}

variable "runner_user" {
  type        = string
  description = "The user to run the runner as"
  default     = "ghrunner"
}

variable "runner_dir" {
  type        = string
  description = "The directory to run the runner in"
  default     = "/opt/github-runner"
}

variable "control_plane_oci_image" {
  type        = string
  description = "The OCI image to deploy"
}

variable "control_plane_version" {
  type        = string
  description = "The version of the control plane to deploy"
  default     = "latest"
}

variable "webhook_secret" {
  type        = string
  description = "The secret to use"
  sensitive   = true
}

variable "github_app_id" {
  type        = string
  description = "The GitHub App ID"
}

variable "github_app_private_key_base64" {
  type        = string
  description = "The base64 encoded private key of the GitHub App"
  sensitive   = true
}

variable "webhook_oci_image" {
  type        = string
  description = "The OCI image to deploy"
  default     = "latest"
}

variable "webhook_version" {
  type        = string
  description = "The version of the webhook to deploy"
  default     = "latest"
}
