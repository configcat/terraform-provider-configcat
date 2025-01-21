variable "webhook_id" {
  type = number
}

data "configcat_webhook_signing_keys" "test" {
  webhook_id = var.webhook_id
}

output "key1" {
  value = data.configcat_webhook_signing_keys.test.key1
}

output "key2" {
  value = data.configcat_webhook_signing_keys.test.key2
}
