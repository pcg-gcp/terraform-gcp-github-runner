resource "random_id" "id" {
  byte_length = 16
}

data "sops_file" "secrets" {
  source_file = "secrets.yaml"
}

module "github_runners" {
  # source = "git@github.com:pcg-gcp/terraform-gcp-github-runner.git"
  source = "../../"

  project_id = "cw-td-sandbox"
  region     = "europe-west1"
  zone       = "europe-west1-b"
  zones      = ["europe-west1-b", "europe-west1-c", "europe-west1-d"]

  vpc_name    = "default"
  subnet_name = "default"

  github_app_private_key_base64 = data.sops_file.secrets.data["github.private_key"]
  github_app_id                 = data.sops_file.secrets.data["github.app_id"]

  runner_image_path       = "projects/cw-td-sandbox/global/images/ubuntu-2204-ghr-20240307-163710"
  runner_machine_type     = "n2-standard-2"
  control_plane_oci_image = "europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane"
  webhook_secret          = random_id.id.hex
  webhook_oci_image       = "europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook"
}
