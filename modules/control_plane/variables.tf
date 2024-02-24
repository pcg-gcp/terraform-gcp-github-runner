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

variable "runner_image_path" {
  type        = string
  description = "The image to deploy"
}

variable "runner_machine_type" {
  type        = string
  description = "The machine type to deploy"
}

variable "image" {
  type        = string
  description = "The image to deploy"
}

variable "max_instance_count" {
  type        = number
  description = "The maximum number of instances to run"
}

variable "private_key_secret_id" {
  type        = string
  description = "The secret ID of the private key"
}

variable "private_key_secret_version" {
  type        = string
  description = "The secret version of the private key"
}
