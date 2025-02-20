output "service_url" {
  value = google_cloud_run_v2_service.control_plane.uri
}

output "invoker_service_account" {
  value = local.invoker_email
}

output "invoker_service_account_id" {
  value = var.disable_service_account_management ? data.google_service_account.invoker[0].id : google_service_account.invoker[0].id
}
