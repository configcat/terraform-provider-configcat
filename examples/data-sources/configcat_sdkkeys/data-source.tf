variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

data "configcat_sdkkeys" "my_sdkkey" {
  config_id      = var.config_id
  environment_id = var.environment_id
}


output "primary_sdkkey" {
  value = data.configcat_sdkkeys.my_sdkkey.primary
}

output "secondary_sdkkey" {
  value = data.configcat_sdkkeys.my_sdkkey.secondary
}
