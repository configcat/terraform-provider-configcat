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
  config_id = var.config_id
  key       = "BoolSettingV2Key"
  name      = "BoolSettingV2Name"
  order     = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = { bool_value = true }

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
      value = { bool_value = true }
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
      value = { bool_value = false }
    },
    {
      percentage_options = [
        {
          percentage = var.percentage
          value      = { bool_value = true }
        },
        {
          percentage = 100 - var.percentage
          value      = { bool_value = false }
        }
      ]
    }
  ]
}
