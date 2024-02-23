resource "random_id" "id" {
  byte_length = 16
}


module "github_runners" {
  source = "git@github.com:pcg-gcp/terraform-gcp-github-runner.git"

  project_id = "cw-td-sandbox"
  region     = "europe-west1"
  zone       = "europe-west1-b"
  zones      = ["europe-west1-b", "europe-west1-c", "europe-west1-d"]

  runner_image_path       = "projects/cw-td-sandbox/global/images/packer-1708178791"
  runner_machine_type     = "e2-micro"
  control_plane_oci_image = "europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane:latest"
  webhook_secret          = random_id.id.hex
  webhook_oci_image       = "europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook:latest"
}
