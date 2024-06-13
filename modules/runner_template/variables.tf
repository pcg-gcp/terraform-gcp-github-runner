variable "project_id" {
  type        = string
  description = "The GCP project ID to deploy all resources into"
}

variable "region" {
  type        = string
  description = "The region to deploy all resources into"
}

variable "vpc_name" {
  type        = string
  description = "The VPC to deploy runner instances into"
}

variable "subnet_name" {
  type        = string
  description = "The subnet to deploy runner instances into"
}

variable "machine_type" {
  type        = string
  description = "The machine type to deploy"
}

variable "image_path" {
  type        = string
  description = "The image to deploy"
}

variable "runner_user" {
  type        = string
  description = "The user to run the runner as"
}

variable "runner_dir" {
  type        = string
  description = "The directory to run the runner in"
}
