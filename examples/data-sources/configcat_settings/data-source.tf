variable "config_id" {
  type = string
}

data "configcat_settings" "settings" {
  config_id        = var.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}


output "setting_id" {
  value = data.configcat_settings.settings.settings[0].setting_id
}
