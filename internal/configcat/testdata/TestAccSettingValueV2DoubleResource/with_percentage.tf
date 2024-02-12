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

variable "percentage" {
  type    = number
  default = null
}

variable "value" {
  type = number
}

resource "configcat_setting" "test" {
  config_id    = var.config_id
  key          = "DoubleSettingV2Key"
  name         = "DoubleSettingV2Name"
  setting_type = "double"
  order        = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = { double_value = var.value }

  targeting_rules = [
    {
      percentage_options = [
        {
          percentage = var.percentage
          value      = { double_value = 10.1 }
        },
        {
          percentage = 100 - var.percentage
          value      = { double_value = 11.1 }
        }
      ]
    }
  ]
}
