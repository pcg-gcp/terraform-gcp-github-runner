resource "google_artifact_registry_repository" "image_cache" {
  location      = var.region
  project       = var.project_id
  repository_id = "ghr-image-cache"
  description   = "A repository for caching images for the Github Action Runner Cloud Run instances."
  format        = "DOCKER"
  mode          = "REMOTE_REPOSITORY"
  remote_repository_config {
    docker_repository {
      custom_repository {
        uri = var.remote_repository_url
      }
    }
  }
}
