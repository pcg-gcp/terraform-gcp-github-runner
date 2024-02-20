resource "google_cloud_tasks_queue" "github_events" {
  name     = "github-events"
  location = var.region
}

resource "google_secret_manager_secret" "webhook_secret" {
  secret_id = "webhook-secret"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "webhook_secret" {
  secret      = google_secret_manager_secret.webhook_secret.id
  secret_data = var.webhook_secret
}

module "control_plane" {
  source     = "./modules/control_plane"
  project_id = var.project_id
  region     = var.region
  zone       = var.zone

  max_instance_count = 2

  runner_image_path   = var.runner_image_path
  runner_machine_type = var.runner_machine_type

  image = var.control_plane_oci_image
}

module "webhook" {
  source     = "./modules/webhook"
  project_id = var.project_id
  region     = var.region

  max_instance_count = 2

  invoker_service_account = module.control_plane.invoker_service_account
  control_plane_url       = module.control_plane.service_url

  image           = var.webhook_oci_image
  task_queue_path = google_cloud_tasks_queue.github_events.id

  webhook_secret_id      = google_secret_manager_secret.webhook_secret.id
  webhook_secret_version = google_secret_manager_secret_version.webhook_secret.version
}
