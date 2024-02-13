
variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "beta_segment_id" {
  type = string
}

resource "configcat_setting" "bool" {
  config_id = var.config_id
  key       = "BoolSettingKey"
  name      = "Bool feature flag"
  order     = 1
}

resource "configcat_setting" "string" {
  config_id    = var.config_id
  key          = "StringSettingKey"
  name         = "String setting"
  order        = 2
  setting_type = "string"
}

resource "configcat_setting_value_v2" "bool_setting_value" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.bool.id

  mandatory_notes = "mandatory notes"

  value = { bool_value = true }

  targeting_rules = [
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "email"
            comparator           = "sensitiveTextEquals"
            comparison_value     = { string_value = "@configcat.com" }
          }
        },
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
        },
        {
          segment_condition = {
            segment_id = var.beta_segment_id
            comparator = "isIn"
          }
        },
        {
          prerequisite_flag_condition = {
            prerequisite_setting_id = configcat_setting.string.id
            comparator              = "doesNotEqual"
            comparison_value        = { string_value = "test" }
          }
        }
      ]
      value = { bool_value = true }
    },

    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "county"
            comparator           = "sensitiveTextEquals"
            comparison_value     = { string_value = "Hungary" }
          }
        },
      ],

      percentage_options = [
        {
          percentage = 50
          value      = { bool_value = true }
        },
        {
          percentage = 50
          value      = { bool_value = false }
        }
      ]
    },

    {
      percentage_options = [
        {
          percentage = 30
          value      = { bool_value = true }
        },
        {
          percentage = 70
          value      = { bool_value = false }
        }
      ]
    },
  ]
}


resource "configcat_setting_value_v2" "string_setting_value" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.string.id

  value = { string_value = "test" }

  targeting_rules = [
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "email"
            comparator           = "sensitiveTextEquals"
            comparison_value     = { string_value = "@configcat.com" }
          }
        },
      ]
      value = { string_value = "custom" }
    },
  ]
}
