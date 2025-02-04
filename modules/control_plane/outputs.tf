output "service_url" {
  value = google_cloud_run_v2_service.control_plane.uri
}

output "invoker_service_account" {
  value = data.google_service_account.invoker.email
}

output "invoker_service_account_id" {
  value = data.google_service_account.invoker.id
}
