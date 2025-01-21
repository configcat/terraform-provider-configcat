variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

resource "configcat_setting" "main" {
  config_id = var.config_id
  key       = "MainSettingV2"
  name      = "MainSettingV2"
  order     = 0
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
resource "configcat_setting" "int" {
  config_id    = var.config_id
  key          = "IntDependencySettingV2"
  name         = "IntDependencySettingV2"
  order        = 3
  setting_type = "int"
}
resource "configcat_setting" "double" {
  config_id    = var.config_id
  key          = "DoubleDependencySettingV2"
  name         = "DoubleDependencySettingV2"
  order        = 4
  setting_type = "double"
}



resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.main.id
  init_only      = false
  value          = { bool_value = true }
  targeting_rules = [
    {
      conditions = [
        {
          prerequisite_flag_condition = {
            prerequisite_setting_id = configcat_setting.bool.id
            comparator              = "equals"
            
            comparison_value        = { 
              # Error 
              # bool_value = true 
            }
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
      ]
      value = { bool_value = false }
    }
  ]
}
