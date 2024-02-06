variable "config_id" {
  type = string
}

variable "key_filter_regex" {
  type    = string
  default = null
}

data "configcat_settings" "test" {
  config_id        = var.config_id
  key_filter_regex = var.key_filter_regex
}

output "setting_id" {
  value = length(data.configcat_settings.test.settings) > 0 ? data.configcat_settings.test.settings[0].setting_id : null
}
