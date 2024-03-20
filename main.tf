resource "google_cloud_tasks_queue" "github_events" {
  name     = "github-job-events"
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


resource "google_secret_manager_secret" "github_auth_secret" {
  secret_id = "github-auth-secret"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "github_auth_secret" {
  secret      = google_secret_manager_secret.github_auth_secret.id
  secret_data = var.github_app_private_key_base64
}

module "runner_template" {
  source     = "./modules/runner_template"
  project_id = var.project_id
  region     = var.region
  zone       = var.zone

  vpc_name    = var.vpc_name
  subnet_name = var.subnet_name

  image_path   = var.runner_image_path
  machine_type = var.runner_machine_type

  runner_user = var.runner_user
  runner_dir  = var.runner_dir
}

module "control_plane" {
  source     = "./modules/control_plane"
  project_id = var.project_id
  region     = var.region
  zone       = var.zone

  max_instance_count = 2

  ephemeral = var.ephemeral

  instance_template_name    = module.runner_template.instance_template_name
  runner_service_account_id = module.runner_template.runner_service_account_id

  github_app_id              = var.github_app_id
  private_key_secret_id      = google_secret_manager_secret.github_auth_secret.id
  private_key_secret_version = google_secret_manager_secret_version.github_auth_secret.version

  image     = var.control_plane_oci_image
  image_tag = var.control_plane_version

  shutdown_schedule          = var.shutdown_schedule
  shutdown_schedule_timezone = var.shutdown_schedule_timezone
  shutdown_attempt_timeout   = var.shutdown_attempt_timeout
}

module "webhook" {
  source     = "./modules/webhook"
  project_id = var.project_id
  region     = var.region

  max_instance_count = 2

  invoker_service_account    = module.control_plane.invoker_service_account
  invoker_service_account_id = module.control_plane.invoker_service_account_id
  control_plane_url          = module.control_plane.service_url

  image           = var.webhook_oci_image
  image_tag       = var.webhook_version
  task_queue_path = google_cloud_tasks_queue.github_events.id

  webhook_secret_id      = google_secret_manager_secret.webhook_secret.id
  webhook_secret_version = google_secret_manager_secret_version.webhook_secret.version
}
