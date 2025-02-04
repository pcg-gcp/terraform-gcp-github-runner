data "google_service_account" "runner" {
  account_id = "ghr-webhook-sa"
  project      = var.project_id
}

resource "google_cloud_run_v2_service" "webhook" {
  project = var.project_id

  name     = "ghr-webhook"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = data.google_service_account.webhook.email
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
            secret  = var.webhook_secret_id
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
