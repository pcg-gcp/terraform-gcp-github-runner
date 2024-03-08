resource "google_service_account" "runner" {
  project      = var.project_id
  account_id   = "ghr-runner"
  display_name = "Runner Service Account"
}

resource "google_project_iam_member" "runner" {
  for_each = toset(["logging.logWriter", "monitoring.metricWriter"])
  project  = var.project_id
  role     = "roles/${each.value}"
  member   = "serviceAccount:${google_service_account.runner.email}"
}

resource "google_storage_bucket" "runner_bucket" {
  project       = var.project_id
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

resource "google_compute_instance_template" "runner" {
  project     = var.project_id
  name        = "runner-template"
  description = "Runner Instance Template"

  machine_type   = var.machine_type
  can_ip_forward = false

  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
  }

  disk {
    source_image = var.image_path
    disk_type    = "pd-balanced"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network    = var.vpc_name
    subnetwork = var.subnet_name
  }

  metadata = {
    startup-script-url = "gs://${google_storage_bucket.runner_bucket.name}/${google_storage_bucket_object.startup_script.name}"
    runner_user        = var.runner_user
    runner_dir         = var.runner_dir
  }

  service_account {
    email  = google_service_account.runner.email
    scopes = ["cloud-platform"]
  }
}
