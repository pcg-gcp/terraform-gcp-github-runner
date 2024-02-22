resource "google_service_account" "webhook" {
  project      = var.project_id
  account_id   = "ghr-webhook"
  display_name = "Github Runner Webhook SA"
}

resource "google_project_iam_member" "webhook" {
  project = var.project_id
  role    = "roles/cloudtasks.enqueuer"
  member  = "serviceAccount:${google_service_account.webhook.email}"
}

resource "google_service_account_iam_member" "admin-account-iam" {
  service_account_id = var.invoker_service_account_id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.webhook.email}"
}

resource "google_cloud_run_v2_service" "webhook" {
  project = var.project_id

  name     = "ghr-webhook"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.webhook.email
    scaling {
      max_instance_count = var.max_instance_count
    }

    containers {
      image = var.image

      env {
        name  = "TASK_QUEUE_PATH"
        value = var.task_queue_path
      }
      env {
        name  = "CONTROL_PLANE_URL"
        value = var.control_plane_url
      }
      env {
        name  = "INVOKER_SERVICE_ACCOUNT"
        value = var.invoker_service_account
      }
      env {
        name = "WEBHOOK_SECRET_KEY"
        value_source {
          secret_key_ref {
            secret  = var.webhook_secret_id
            version = var.webhook_secret_version
          }
        }
      }
    }
  }

}

resource "google_cloud_run_v2_service_iam_binding" "webhook" {
  project  = google_cloud_run_v2_service.webhook.project
  location = google_cloud_run_v2_service.webhook.location
  name     = google_cloud_run_v2_service.webhook.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
