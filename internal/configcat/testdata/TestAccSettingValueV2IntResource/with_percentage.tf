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
  key          = "IntSettingV2Key"
  name         = "IntSettingV2Name"
  setting_type = "int"
  order        = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = { int_value = var.value }

  targeting_rules = [
    {
      percentage_options = [
        {
          percentage = var.percentage
          value      = { int_value = 10 }
        },
        {
          percentage = 100 - var.percentage
          value      = { int_value = 11 }
        }
      ]
    }
  ]
}
