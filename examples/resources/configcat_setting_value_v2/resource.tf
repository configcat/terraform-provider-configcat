
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
  key       = "BoolDependencySettingV2"
  name      = "BoolDependencySettingV2"
  order     = 1
}

resource "configcat_setting" "string" {
  config_id    = var.config_id
  key          = "StringDependencySettingV2"
  name         = "StringDependencySettingV2"
  order        = 2
  setting_type = "string"
}

variable "dependency_setting_id" {
  type = string
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
          prerequisite_flag_condition = {
            prerequisite_setting_id = configcat_setting.int.id
            comparator              = "equals"
            comparison_value        = { int_value = 1 }
          }
        },
        {
          prerequisite_flag_condition = {
            prerequisite_setting_id = configcat_setting.double.id
            comparator              = "doesNotEqual"
            comparison_value        = { double_value = 1.1 }
          }
        }
      ],

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
