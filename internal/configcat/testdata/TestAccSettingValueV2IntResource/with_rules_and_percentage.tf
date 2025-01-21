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
  key          = "IntSettingV2Key"
  name         = "IntSettingV2Name"
  setting_type = "int"
  order        = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = { int_value = 30 }

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
      value = { int_value = 31 }
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
      value = { int_value = 32 }
    },
    {
      percentage_options = [
        {
          percentage = var.percentage
          value      = { int_value = 33 }
        },
        {
          percentage = 100 - var.percentage
          value      = { int_value = 34 }
        }
      ]
    }
  ]
}
