variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy the repository to."
}

variable "region" {
  type        = string
  description = "The region to deploy the repository to."
  default     = "europe-west1"
}

variable "remote_repository_url" {
  type        = string
  description = "The URL of the remote repository to clone."
}
