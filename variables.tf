variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
  default     = "europe-west1"
}

variable "allowed_zones" {
  type        = list(string)
  description = <<-EOT
  The zones to deploy the runner instances into.
  If not set the runner will be deployed into any zone in the region.
  Also see use_strict_zone_order
  EOT
  default     = []
}

variable "use_strict_zone_order" {
  type        = bool
  description = <<-EOT
  If this is set to true and allowed_zones is set the runner will always be deployed in the first available zone in the list unless it is unavailable.
  If allowed_zones is not set the first zone returned by the API will be used.
  Otherwise the runner will be deployed in a random zone either from the allowed_zones list or from the API.
  EOT
  default     = false
}

variable "disable_service_account_management" {
  type        = bool
  description = "Whether the used service accounts should be create by this module."
  default     = false
}

variable "enable_apis" {
  type        = bool
  description = "Whether this module should enable the required apis."
  default     = true
}

variable "runner_account_id" {
  type        = string
  description = "Account id of the runner service account only used if service account management is disabled."
  default     = "ghr-runner"
}

variable "control_plane_account_id" {
  type        = string
  description = "Account id of the control plane service account only used if service account management is disabled."
  default     = "ghr-control-plane"
}

variable "invoker_account_id" {
  type        = string
  description = "Account id of the invoker service account only used if service account management is disabled."
  default     = "ghr-cp-invoker"
}

variable "webhook_account_id" {
  type        = string
  description = "Account id of the webhook service account only used if service account management is disabled."
  default     = "ghr-webhook"
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

variable "use_remote_repository" {
  type        = bool
  description = "Whether to use a remote repository for caching images"
  default     = true
}

variable "remote_repository_url" {
  type        = string
  description = "The URL of the remote repository to clone. This is only used if use_remote_repository is set to true"
  default     = "https://ghcr.io"
}

variable "remote_webhook_image_name" {
  type        = string
  description = "The name of the webhook image in the remote repository. This should only be the image name as it is combined with the repository url to form the full path. This is only used if use_remote_repository is set to true"
  default     = "pcg-gcp/terraform-gcp-github-runner/webhook"
}

variable "remote_control_plane_image_name" {
  type        = string
  description = "The name of the control plane image in the remote repository. This should only be the image name as it is combined with the repository url to form the full path. This is only used if use_remote_repository is set to true"
  default     = "pcg-gcp/terraform-gcp-github-runner/control-plane"
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

variable "runner_disk_type" {
  type        = string
  description = "The disk type to deploy"
  default     = "pd-balanced"
}

variable "runner_disk_size_gb" {
  type        = number
  description = "The disk size in GB to deploy"
  default     = null
}

variable "runner_additional_disks" {
  type        = list(any)
  description = "Additional disks to attach to the runner"
  default     = []
}

variable "runner_on_host_maintenance" {
  type        = string
  description = "The maintenance policy for the runner"
  default     = "MIGRATE"
}

variable "runner_use_spot_vms" {
  type        = bool
  description = "Whether to use spot VMs for the runner"
  default     = false
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
  description = "The control plane OCI image to deploy. This needs to be the full path without the image tag. This is only used if use_remote_repository is set to false"
  default     = ""
}

variable "control_plane_version" {
  type        = string
  description = "The version of the control plane to deploy. If not set the module version will be used"
  default     = ""
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
  description = "The webhook OCI image to deploy. This needs to be the full path withouth the image tag. This is only used if use_remote_repository is set to false"
  default     = ""
}

variable "webhook_version" {
  type        = string
  description = "The version of the webhook to deploy. If not set the module version will be used"
  default     = ""
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
