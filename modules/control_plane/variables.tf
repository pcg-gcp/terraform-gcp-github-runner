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
}

variable "runner_dir" {
  type        = string
  description = "The directory to run the runner in"
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
