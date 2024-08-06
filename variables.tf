variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
  default     = "europe-west1"
}

variable "zone" {
  type        = string
  description = "The zone to deploy runner instances into"
  default     = "europe-west1-b"
}

variable "allowed_zones" {
  type        = list(string)
  description = "The zones to deploy runner instances into"
  default     = ["europe-west1-b", "europe-west1-c", "europe-west1-d"]
}

variable "enable_debug" {
  type        = bool
  description = "Whether to enable debug mode"
  default     = false
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
  default     = "ubuntu-os-cloud/ubuntu-2404-lts-amd64"
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

variable "include_install_step" {
  type        = bool
  description = "Whether to include the install step for the setup script"
  default     = true
}

variable "runner_version" {
  description = "GitHub Runner version to install"
  type        = string
  default     = ""
}

variable "runner_extra_labels" {
  description = "Github runner extra labels. These should not include github read-only labels like 'self-hosted' or 'linux'"
  type        = list(string)
  default     = []
}

variable "node_version" {
  description = "NodeJS version to install"
  type        = string
  default     = ""
}

variable "ephemeral" {
  type        = bool
  description = "Whether to use ephemeral runners"
  default     = false
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
  default     = 10
}

variable "min_runner_count" {
  type        = number
  description = "The minimum number of runners that should be deployed at all times"
  default     = 0
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

variable "forward_delay_seconds" {
  type        = number
  description = "The number of seconds the webhook handler delays events before forwarding them to the control plane"
  default     = 10
}

variable "shutdown_schedule" {
  type        = string
  description = "The shutdown schedule in cron format"
  default     = "*/5 * * * *"
}

variable "shutdown_schedule_timezone" {
  type        = string
  description = "The timezone of the shutdown schedule"
  default     = "Etc/UTC"
}

variable "shutdown_attempt_timeout" {
  type        = string
  description = "The timeout for the shutdown attempt"
  default     = "320s"
}
