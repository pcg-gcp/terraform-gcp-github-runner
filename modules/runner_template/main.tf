data "http" "github_runner_release_json" {
  url = "https://api.github.com/repos/actions/runner/releases/latest"

  request_headers = {
    Accept = "application/vnd.github+json"
    X-GitHub-Api-Version : "2022-11-28"
  }
}

locals {
  effective_node_version = coalesce(var.node_version, "lts")
  runner_version         = coalesce(var.runner_version, trimprefix(jsondecode(data.http.github_runner_release_json.response_body).tag_name, "v"))
}

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
  name = "setup_runner.sh"
  content = templatefile("../../templates/runners/setup_runner.tftpl", {
    include_install     = var.include_install_step,
    include_run         = true,
    node_version        = local.effective_node_version,
    runner_user         = var.runner_user,
    runner_dir          = var.runner_dir,
    runner_download_url = "https://github.com/actions/runner/releases/download/v${local.runner_version}/actions-runner-linux-x64-${local.runner_version}.tar.gz"
  })
  bucket = google_storage_bucket.runner_bucket.name
}

resource "google_storage_bucket_iam_member" "runner" {
  bucket = google_storage_bucket.runner_bucket.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.runner.email}"
}

locals {
  automatic_restart  = var.use_spot_vms ? false : true
  preemptible        = var.use_spot_vms
  provisioning_model = var.use_spot_vms ? "SPOT" : "STANDARD"
}

resource "google_compute_instance_template" "runner" {
  project     = var.project_id
  name        = "runner-template"
  description = "Runner Instance Template"
  region      = var.region

  machine_type   = var.machine_type
  can_ip_forward = false

  scheduling {
    automatic_restart   = local.automatic_restart
    on_host_maintenance = var.on_host_maintenance
    preemptible         = local.preemptible
    provisioning_model  = local.provisioning_model
  }

  disk {
    source_image = var.image_path
    disk_type    = var.disk_type
    disk_size_gb = var.disk_size_gb
    auto_delete  = true
    boot         = true
  }

  dynamic "disk" {
    for_each = var.additional_disks
    content {
      source_image     = disk.value.source_image
      disk_type        = disk.value.disk_type
      disk_size_gb     = disk.value.disk_size_gb
      auto_delete      = disk.value.auto_delete
      provisioned_iops = disk.value.provisioned_iops
      type             = disk.value.type
      source_snapshot  = disk.value.source_snapshot
      boot             = false
    }
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
