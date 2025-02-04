output "instance_template_name" {
  value = google_compute_instance_template.runner.name
}

output "runner_service_account_id" {
  value = data.google_service_account.runner.id
}
