output "webhook_secret" {
  value     = var.webhook_secret
  sensitive = true
}

output "webhook_url" {
  value = module.webhook.service_url
}
