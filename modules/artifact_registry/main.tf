resource "google_artifact_registry_repository" "my-repo" {
  location      = var.region
  repository_id = "ghr-image-cache"
  description   = "A repository for caching images for the Github Action Runner Cloud Run instances."
  format        = "DOCKER"
  mode          = "REMOTE_REPOSITORY"
  remote_repository_config {
    docker_repository {
      custom_repository {
        uri = var.repository_uri
      }
    }
  }
}
