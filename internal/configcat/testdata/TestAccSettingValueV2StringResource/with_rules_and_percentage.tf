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
  value          = { string_value = "vc" }

  targeting_rules = [
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "email"
            comparator           = "sensitiveTextEquals"
            comparison_value     = { string_value = "@configcat.com" }
          }
        }
      ],
      value = { string_value = "te" }
    },
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "color"
            comparator           = "isOneOf"
            comparison_value = {
              list_values = [
                { value = "#000000", hint = "black" },
                { value = "red" },
              ]
            }
          }
        }
      ],
      value = { string_value = "tf" }
    },
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
