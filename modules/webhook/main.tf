resource "google_service_account" "webhook" {
  count = var.disable_service_account_management ? 0 : 1

  project      = var.project_id
  account_id   = "ghr-webhook"
  display_name = "Github Runner Webhook SA"
}

data "google_service_account" "webhook" {
  count = var.disable_service_account_management ? 1 : 0

  project    = var.project_id
  account_id = var.webhook_account_id
}

locals {
  webhook_email = var.disable_service_account_management ? data.google_service_account.webhook[0].email : google_service_account.webhook[0].email
}

resource "google_project_iam_member" "webhook" {
  count = var.disable_service_account_management ? 0 : 1

  project = var.project_id
  role    = "roles/cloudtasks.enqueuer"
  member  = "serviceAccount:${local.webhook_email}"
}

resource "google_service_account_iam_member" "admin-account-iam" {
  count = var.disable_service_account_management ? 0 : 1

  service_account_id = var.invoker_service_account_id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.webhook_email}"
}

resource "google_secret_manager_secret_iam_member" "webhook" {
  secret_id = var.webhook_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${local.webhook_email}"
}

resource "google_cloud_run_v2_service" "webhook" {
  project = var.project_id

  name     = "ghr-webhook"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = local.webhook_email
    scaling {
      max_instance_count = var.max_instance_count
    }

    containers {
      image = "${var.image}:${var.image_tag}"

      env {
        name  = "TASK_QUEUE_PATH"
        value = var.task_queue_path
      }
      env {
        name  = "CONTROL_PLANE_URL"
        value = "${var.control_plane_url}/startup"
      }
      env {
        name  = "INVOKER_SERVICE_ACCOUNT"
        value = var.invoker_service_account
      }
      env {
        name  = "ENABLE_DEBUG"
        value = var.enable_debug
      }
      env {
        name  = "DELAY_SECONDS"
        value = var.forward_delay_seconds
      }
      env {
        name  = "RUNNER_LABELS"
        value = join(",", var.runner_labels)
      }
      env {
        name = "WEBHOOK_SECRET_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret_iam_member.webhook.secret_id
            version = var.webhook_secret_version
          }
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [client, client_version]
  }

}

resource "google_cloud_run_v2_service_iam_binding" "webhook" {
  project  = google_cloud_run_v2_service.webhook.project
  location = google_cloud_run_v2_service.webhook.location
  name     = google_cloud_run_v2_service.webhook.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
