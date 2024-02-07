variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "init_only" {
  type    = bool
  default = false
}

variable "value" {
  type = string
}

resource "configcat_setting" "test" {
  config_id    = var.config_id
  key          = "StringSettingKey"
  name         = "StringSettingName"
  setting_type = "string"
  order        = 0
}

resource "configcat_setting_value" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = var.value
}
