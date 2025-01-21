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
  key          = "IntSettingKey"
  name         = "IntSettingName"
  setting_type = "int"
  order        = 0
}

resource "configcat_setting_value" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = "40"

  rollout_rules {
    comparison_attribute = "email"
    comparator           = "contains"
    comparison_value     = "@configcat.com"
    value                = "50"
  }

  rollout_rules {
    comparison_attribute = "color"
    comparator           = "isOneOf"
    comparison_value     = "red"
    value                = "60"
  }


  percentage_items {
    percentage = 10
    value      = "70"
  }
  percentage_items {
    percentage = 20
    value      = "80"
  }
  percentage_items {
    percentage = 70
    value      = "90"
  }
}
