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

resource "google_project_iam_member" "control_plane" {
  for_each = toset(["compute.admin", "iam.serviceAccountUser"])
  project  = var.project_id
  role     = "roles/${each.value}"
  member   = "serviceAccount:${google_service_account.control_plane.email}"
}

resource "google_secret_manager_secret_iam_member" "control_plane" {
  secret_id = var.private_key_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.control_plane.email}"
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
      image = var.image

      env {
        name  = "PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "ZONE"
        value = var.zone
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
