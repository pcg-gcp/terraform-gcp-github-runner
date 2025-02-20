output "instance_template_name" {
  value = google_compute_instance_template.runner.name
}

output "runner_service_account_id" {
  value = var.disable_service_account_management ? data.google_service_account.runner[0].id : google_service_account.runner[0].id
}
