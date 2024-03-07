resource "google_service_account" "control_plane" {
  project      = var.project_id
  account_id   = "ghr-control-plane"
  display_name = "Github Runner CP SA"
}

resource "google_service_account" "invoker" {
  project      = var.project_id
  account_id   = "ghr-cp-invoker"
  display_name = "Invoker Service Account"
}

resource "google_service_account" "runner" {
  project      = var.project_id
  account_id   = "ghr-runner"
  display_name = "Runner Service Account"
}

resource "google_project_iam_member" "control_plane" {
  for_each = toset(["compute.admin"])
  project  = var.project_id
  role     = "roles/${each.value}"
  member   = "serviceAccount:${google_service_account.control_plane.email}"
}

resource "google_secret_manager_secret_iam_member" "control_plane" {
  secret_id = var.private_key_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.control_plane.email}"
}

resource "google_service_account_iam_member" "runner_user" {
  service_account_id = google_service_account.runner.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.control_plane.email}"
}

resource "google_project_iam_member" "runner" {
  for_each = toset(["logging.logWriter", "monitoring.metricWriter"])
  project  = var.project_id
  role     = "roles/${each.value}"
  member   = "serviceAccount:${google_service_account.runner.email}"
}

resource "google_storage_bucket" "runner_bucket" {
  name          = "ghr-scripts-bucket"
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "startup_script" {
  name   = "startup.sh"
  source = "../../templates/runners/start_runner.sh"
  bucket = google_storage_bucket.runner_bucket.name
}

resource "google_storage_bucket_iam_member" "runner" {
  bucket = google_storage_bucket.runner_bucket.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.runner.email}"
}

resource "google_cloud_run_v2_service" "control_plane" {
  project = var.project_id

  name     = "ghr-control-plane"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.control_plane.email
    scaling {
      max_instance_count = var.max_instance_count
    }

    containers {
      image = "${var.image}:${var.image_tag}"

      env {
        name  = "PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "REGION"
        value = var.region
      }
      env {
        name  = "ZONE"
        value = var.zone
      }
      env {
        name  = "NETWORK"
        value = var.vpc_name
      }
      env {
        name  = "SUBNET"
        value = var.subnet_name
      }
      env {
        name  = "RUNNER_SERVICE_ACCOUNT"
        value = google_service_account.runner.email
      }
      env {
        name  = "RUNNER_USER"
        value = var.runner_user
      }
      env {
        name  = "RUNNER_DIR"
        value = var.runner_dir
      }
      env {
        name  = "IMAGE_PATH"
        value = var.runner_image_path
      }
      env {
        name  = "MACHINE_TYPE"
        value = var.runner_machine_type
      }
      env {
        name  = "GITHUB_APP_ID"
        value = var.github_app_id
      }
      env {
        name  = "STARTUP_SCRIPT_URL"
        value = google_storage_bucket_object.startup_script.self_link
      }
      env {
        name = "GITHUB_APP_PRIVATE_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret_iam_member.control_plane.secret_id
            version = var.private_key_secret_version
          }
        }
      }

    }
  }

  lifecycle {
    ignore_changes = [client, client_version]
  }

}

resource "google_cloud_run_v2_service_iam_binding" "control_plane" {
  project  = google_cloud_run_v2_service.control_plane.project
  location = google_cloud_run_v2_service.control_plane.location
  name     = google_cloud_run_v2_service.control_plane.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}
