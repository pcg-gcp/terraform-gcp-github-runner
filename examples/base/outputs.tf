output "webhook_secret" {
  value     = module.github_runners.webhook_secret
  sensitive = true
}

output "webhook_url" {
  value = module.github_runners.webhook_url
}
