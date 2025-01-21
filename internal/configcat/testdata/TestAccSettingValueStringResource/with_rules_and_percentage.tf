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
  value          = "Vvalue"

  rollout_rules {
    comparison_attribute = "email"
    comparator           = "contains"
    comparison_value     = "@configcat.com"
    value                = "RRValue1"
  }

  rollout_rules {
    comparison_attribute = "color"
    comparator           = "isOneOf"
    comparison_value     = "red"
    value                = "RRValue2"
  }


  percentage_items {
    percentage = 10
    value      = "P1value"
  }
  percentage_items {
    percentage = 20
    value      = "P2value"
  }
  percentage_items {
    percentage = 70
    value      = "P3value"
  }
}
