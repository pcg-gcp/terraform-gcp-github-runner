packer {
  required_plugins {
    googlecompute = {
      version = ">= 1.1.4"
      source  = "github.com/hashicorp/googlecompute"
    }
  }
}

variable "project_id" {
  type    = string
  default = ""
}

variable "zone" {
  type    = string
  default = ""
}

variable "builder_sa" {
  type    = string
  default = ""
}

variable "custom_shell_commands" {
  description = "Additional commands to run on the instance, to customize the instance, like installing packages"
  type        = list(string)
  default     = []
}

variable "runner_user" {
  description = "User to run the GitHub Runner"
  type        = string
  default     = "ghrunner"
}

variable "runner_dir" {
  description = "Directory to install the GitHub Runner"
  type        = string
  default     = "/opt/github-runner"
}

variable "runner_version" {
  description = "GitHub Runner version to install"
  type        = string
  default     = ""
}

variable "node_version" {
  description = "NodeJS version to install"
  type        = string
  default     = ""
}

data "http" github_runner_release_json {
  url = "https://api.github.com/repos/actions/runner/releases/latest"
  request_headers = {
    Accept = "application/vnd.github+json"
    X-GitHub-Api-Version : "2022-11-28"
  }
}

locals {
  runner_version         = coalesce(var.runner_version, trimprefix(jsondecode(data.http.github_runner_release_json.body).tag_name, "v"))
  effective_node_version = coalesce(var.node_version, "lts")
}

source "googlecompute" "github-runner" {
  project_id            = var.project_id
  source_image_family   = "ubuntu-2404-lts-amd64"
  zone                  = var.zone
  service_account_email = var.builder_sa
  ssh_username          = "root"
  image_name = "ubuntu-2404-ghr-${formatdate("YYYYMMDD-hhmmss",timestamp())}"
}

build {
  name    = "githubactions-runner"
  sources = ["sources.googlecompute.github-runner"]

  provisioner "file" {
    source      = "../../templates/runners/install_runner.sh"
    destination = "/tmp/install_runner.sh"
  }

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "NODE_VERSION=${local.effective_node_version}",
      "RUNNER_USER=${var.runner_user}",
      "RUNNER_DIR=${var.runner_dir}",
      "RUNNER_DOWNLOAD_URL=https://github.com/actions/runner/releases/download/v${local.runner_version}/actions-runner-linux-x64-${local.runner_version}.tar.gz",
    ]

    inline = concat([
      "/bin/bash /tmp/install_runner.sh",
      "rm -f /tmp/install_runner.sh",
    ], var.custom_shell_commands)
  }

}
