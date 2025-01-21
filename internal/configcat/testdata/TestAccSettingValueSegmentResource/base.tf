variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

resource "configcat_setting" "test" {
  config_id = var.config_id
  key       = "SegmentSettingKey"
  name      = "SegmentSettingName"
  order     = 0
}

resource "configcat_setting_value" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = false
  value          = "true"
}
