provider "sops" {}

terraform {
  required_version = ">= 1.8"
  required_providers {
    sops = {
      source  = "carlpett/sops"
      version = "~> 1.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}
