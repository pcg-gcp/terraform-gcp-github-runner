resource "google_service_account" "control_plane" {
  count = var.disable_service_account_management ? 0 : 1

  project      = var.project_id
  account_id   = "ghr-control-plane"
  display_name = "Github Runner CP SA"
}

data "google_service_account" "control_plane" {
  count = var.disable_service_account_management ? 1 : 0

  project    = var.project_id
  account_id = var.control_plane_account_id
}

resource "google_service_account" "invoker" {
  count = var.disable_service_account_management ? 0 : 1

  project      = var.project_id
  account_id   = "ghr-cp-invoker"
  display_name = "Invoker Service Account"
}

data "google_service_account" "invoker" {
  count = var.disable_service_account_management ? 1 : 0

  project    = var.project_id
  account_id = var.invoker_account_id
}

locals {
  control_plane_email = var.disable_service_account_management ? data.google_service_account.control_plane[0].email : google_service_account.control_plane[0].email
  invoker_email       = var.disable_service_account_management ? data.google_service_account.invoker[0].email : google_service_account.invoker[0].email
}

resource "google_project_iam_member" "control_plane" {
  for_each = toset(var.disable_service_account_management ? [] : ["compute.admin", "secretmanager.admin"])

  project = var.project_id
  role    = "roles/${each.value}"
  member  = "serviceAccount:${local.control_plane_email}"
}

resource "google_secret_manager_secret_iam_member" "control_plane" {
  secret_id = var.private_key_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${local.control_plane_email}"
}

resource "google_service_account_iam_member" "runner_user" {
  count = var.disable_service_account_management ? 0 : 1

  service_account_id = var.runner_service_account_id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.control_plane_email}"
}

resource "google_cloud_run_v2_service" "control_plane" {
  project = var.project_id

  name     = "ghr-control-plane"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = local.control_plane_email
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
  members  = ["serviceAccount:${local.invoker_email}"]
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
      service_account_email = google_service_account.invoker.email
    }
  }
}

