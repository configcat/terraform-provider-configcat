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
  type = string
}

resource "configcat_setting" "test" {
  config_id    = var.config_id
  key          = "StringSettingV2Key"
  name         = "StringSettingV2Name"
  setting_type = "string"
  order        = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = { string_value = var.value }

  targeting_rules = [
    {
      percentage_options = [
        {
          percentage = var.percentage
          value      = { string_value = "tpa" }
        },
        {
          percentage = 100 - var.percentage
          value      = { string_value = "tpb" }
        }
      ]
    }
  ]
}
