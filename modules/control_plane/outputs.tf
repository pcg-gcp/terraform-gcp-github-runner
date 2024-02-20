output "service_url" {
  value = google_cloud_run_v2_service.control_plane.uri
}

output "invoker_service_account" {
  value = google_service_account.invoker.email
}
