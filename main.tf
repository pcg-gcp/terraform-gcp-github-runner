locals {
  module_version = "v0.4.1"

  runner_labels = sort(distinct(concat(["self-hosted", "linux", "x64"], var.runner_extra_labels)))
  required_services = concat(
    ["compute.googleapis.com", "run.googleapis.com", "cloudtasks.googleapis.com", "secretmanager.googleapis.com", "cloudscheduler.googleapis.com"],
    var.use_remote_repository ? ["artifactregistry.googleapis.com"] : []
  )

  effective_control_plane_version = var.control_plane_version == "" ? local.module_version : var.control_plane_version
  effective_webhook_version       = var.webhook_version == "" ? local.module_version : var.webhook_version
}

resource "google_project_service" "required_services" {
  for_each = toset(local.required_services)
  project  = var.project_id
  service  = each.key
}

resource "random_string" "queue_suffix" {
  length  = 5
  special = false
  upper   = false
  numeric = true
  lower   = true
}

resource "google_cloud_tasks_queue" "github_events" {
  name     = "github-job-events-${random_string.queue_suffix.result}"
  location = var.region

  depends_on = [google_project_service.required_services["cloudtasks.googleapis.com"]]
}

resource "google_secret_manager_secret" "webhook_secret" {
  secret_id = "webhook-secret"

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }

  depends_on = [google_project_service.required_services["secretmanager.googleapis.com"]]
}

resource "google_secret_manager_secret_version" "webhook_secret" {
  secret      = google_secret_manager_secret.webhook_secret.id
  secret_data = var.webhook_secret
}


resource "google_secret_manager_secret" "github_auth_secret" {
  secret_id = "github-auth-secret"

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }

  depends_on = [google_project_service.required_services["secretmanager.googleapis.com"]]
}

resource "google_secret_manager_secret_version" "github_auth_secret" {
  secret      = google_secret_manager_secret.github_auth_secret.id
  secret_data = var.github_app_private_key_base64
}

module "artifact_registry" {
  count  = var.use_remote_repository ? 1 : 0
  source = "./modules/artifact_registry"

  project_id            = var.project_id
  region                = var.region
  remote_repository_url = var.remote_repository_url
  depends_on            = [google_project_service.required_services]
}

module "runner_template" {
  source     = "./modules/runner_template"
  project_id = var.project_id
  region     = var.region

  vpc_name    = var.vpc_name
  subnet_name = var.subnet_name

  image_path          = var.runner_image_path
  machine_type        = var.runner_machine_type
  on_host_maintenance = var.runner_on_host_maintenance
  use_spot_vms        = var.runner_use_spot_vms


  disk_type        = var.runner_disk_type
  disk_size_gb     = var.runner_disk_size_gb
  additional_disks = var.runner_additional_disks

  runner_user          = var.runner_user
  runner_dir           = var.runner_dir
  runner_version       = var.runner_version
  node_version         = var.node_version
  include_install_step = var.include_install_step

  depends_on = [google_project_service.required_services]
}

module "control_plane" {
  source     = "./modules/control_plane"
  project_id = var.project_id
  region     = var.region

  allowed_zones         = var.allowed_zones
  use_strict_zone_order = var.use_strict_zone_order

  enable_debug = var.enable_debug

  max_instance_count = 2

  ephemeral       = var.ephemeral
  use_jit_config  = var.use_jit_config
  use_org_runners = var.use_org_runners

  min_runner_count = var.min_runner_count
  max_runner_count = var.max_runner_count

  instance_template_name    = module.runner_template.instance_template_name
  runner_service_account_id = module.runner_template.runner_service_account_id
  runner_labels             = local.runner_labels

  github_app_id              = var.github_app_id
  private_key_secret_id      = google_secret_manager_secret.github_auth_secret.id
  private_key_secret_version = google_secret_manager_secret_version.github_auth_secret.version

  image     = var.use_remote_repository ? "${module.artifact_registry[0].image_cache_url}/${var.remote_control_plane_image_name}" : var.control_plane_oci_image
  image_tag = local.effective_control_plane_version

  shutdown_schedule          = var.shutdown_schedule
  shutdown_schedule_timezone = var.shutdown_schedule_timezone
  shutdown_attempt_timeout   = var.shutdown_attempt_timeout

  depends_on = [google_project_service.required_services]
}

module "webhook" {
  source     = "./modules/webhook"
  project_id = var.project_id
  region     = var.region

  enable_debug = var.enable_debug

  max_instance_count = 2

  invoker_service_account    = module.control_plane.invoker_service_account
  invoker_service_account_id = module.control_plane.invoker_service_account_id
  control_plane_url          = module.control_plane.service_url

  image           = var.use_remote_repository ? "${module.artifact_registry[0].image_cache_url}/${var.remote_webhook_image_name}" : var.webhook_oci_image
  image_tag       = local.effective_webhook_version
  task_queue_path = google_cloud_tasks_queue.github_events.id

  webhook_secret_id      = google_secret_manager_secret.webhook_secret.id
  webhook_secret_version = google_secret_manager_secret_version.webhook_secret.version

  forward_delay_seconds = var.forward_delay_seconds
  runner_labels         = local.runner_labels

  depends_on = [google_project_service.required_services]
}
