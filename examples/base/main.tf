resource "random_id" "id" {
  byte_length = 16
}

data "sops_file" "secrets" {
  source_file = "secrets.yaml"
}

module "github_runners" {
  # source = "git@github.com:pcg-gcp/terraform-gcp-github-runner.git"
  source = "../../"

  enable_debug = true

  project_id = "cw-td-sandbox"
  region     = "europe-west1"
  zone       = "europe-west1-b"

  vpc_name    = "default"
  subnet_name = "default"

  github_app_private_key_base64 = data.sops_file.secrets.data["github.private_key"]
  github_app_id                 = data.sops_file.secrets.data["github.app_id"]

  ephemeral      = true
  use_jit_config = true

  runner_image_path       = "projects/cw-td-sandbox/global/images/ubuntu-2404-ghr-20240614-080646"
  runner_machine_type     = "n2-standard-2"
  control_plane_oci_image = "europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane"
  webhook_secret          = random_id.id.hex
  webhook_oci_image       = "europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook"

  include_install_step = false

  forward_delay_seconds = 0
}
