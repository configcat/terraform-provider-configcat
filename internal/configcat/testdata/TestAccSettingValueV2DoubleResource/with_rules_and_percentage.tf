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
  key          = "DoubleSettingV2Key"
  name         = "DoubleSettingV2Name"
  setting_type = "double"
  order        = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = { double_value = 30.1 }

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
      value = { double_value = 31.1 }
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
      value = { double_value = 32.1 }
    },
    {
      percentage_options = [
        {
          percentage = var.percentage
          value      = { double_value = 33.1 }
        },
        {
          percentage = 100 - var.percentage
          value      = { double_value = 34.1 }
        }
      ]
    }
  ]
}
