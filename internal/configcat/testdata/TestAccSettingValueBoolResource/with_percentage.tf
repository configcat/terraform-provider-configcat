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

variable "percentage" {
  type    = number
  default = null
}

resource "configcat_setting" "test" {
  config_id = var.config_id
  key       = "BoolSettingKey"
  name      = "BoolSettingName"
  order     = 0
}

resource "configcat_setting_value" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = var.value

  percentage_items {
    percentage = var.percentage
    value      = "true"
  }
  percentage_items {
    percentage = 100 - var.percentage
    value      = "false"
  }
}
