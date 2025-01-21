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
  key          = "StringSettingV2Key"
  name         = "StringSettingV2Name"
  setting_type = "string"
  order        = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = { string_value = "va" }

  targeting_rules = [
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "email"
            comparator           = "sensitiveTextEquals"
            comparison_value     = { string_value = "jane@configcat.com" }
          }
        }
      ],
      value = { string_value = "tc" }
    },
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "email"
            comparator           = "containsAnyOf"
            comparison_value = {
              list_values = [
                { value = "@configcat.com", hint = "the greatest company of the world" },
                { value = "@example.com" }
              ]
            }
          }
        }
      ],
      value = { string_value = "td" }
    }
  ]
}
