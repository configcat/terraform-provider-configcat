variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

resource "configcat_webhook" "my_webhook" {
  config_id      = var.config_id
  environment_id = var.environment_id

  url         = "https://example.com"
  http_method = "post"
  content     = "The HTTP body content."
  webhook_headers = [
    { "webhookheaderkey" = "webhookheadervalue" }
  ]
  secure_webhook_headers = [
    { "securewebhookheaderkey" = "securewebhookheadervalue" }
  ]
}


output "webhook_id" {
  value = configcat_webhook.my_webhook.id
}
