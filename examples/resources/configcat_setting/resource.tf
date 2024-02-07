variable "config_id" {
  type = string
}

resource "configcat_setting" "my_setting" {
  config_id    = var.config_id
  key          = "isAwesomeFeatureEnabled"
  name         = "My awesome feature flag"
  hint         = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
  order        = 0
}


output "setting_id" {
  value = configcat_setting.my_setting.id
}