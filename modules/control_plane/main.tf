data "google_service_account" "control_planer" {
  account_id = "ghr-control-plane-sa"
  project      = var.project_id
}

data "google_service_account" "invoker" {
  account_id = "ghr-cp-invoker-sa"
  project      = var.project_id
}

resource "google_cloud_run_v2_service" "control_plane" {
  project = var.project_id

  name     = "ghr-control-plane"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = data.google_service_account.control_plane.email
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
        name  = "ALLOWED_ZONES"
        value = join(",", var.allowed_zones)
      }
      env {
        name  = "USE_STRICT_ZONE_ORDER"
        value = var.use_strict_zone_order
      }
      env {
        name  = "ENABLE_DEBUG"
        value = var.enable_debug
      }
      env {
        name  = "INSTANCE_TEMPLATE_NAME"
        value = var.instance_template_name
      }
      env {
        name  = "GITHUB_APP_ID"
        value = var.github_app_id
      }
      env {
        name  = "EPHEMERAL"
        value = var.ephemeral
      }
      env {
        name  = "USE_JIT_CONFIG"
        value = var.use_jit_config
      }
      env {
        name  = "USE_ORG_RUNNERS"
        value = var.use_org_runners
      }
      env {
        name  = "MAX_RUNNER_COUNT"
        value = var.max_runner_count
      }
      env {
        name  = "MIN_RUNNER_COUNT"
        value = var.min_runner_count
      }
      env {
        name  = "RUNNER_LABELS"
        value = join(",", var.runner_labels)
      }
      env {
        name = "GITHUB_APP_PRIVATE_KEY"
        value_source {
          secret_key_ref {
            secret  = var.private_key_secret_id
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

resource "google_cloud_scheduler_job" "shutdown" {
  name             = "shutdown-scheduler"
  description      = "Trigger the control plane to check for runners to shutdown"
  schedule         = var.shutdown_schedule
  time_zone        = var.shutdown_schedule_timezone
  attempt_deadline = var.shutdown_attempt_timeout

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "${google_cloud_run_v2_service.control_plane.uri}/shutdown"

    oidc_token {
      service_account_email = data.google_service_account.invoker.email
    }
  }
}

