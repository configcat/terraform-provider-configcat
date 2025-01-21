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
            comparison_attribute = "utcNow"
            comparator           = "dateTimeBefore"
            comparison_value     = { double_value = 1707752290.1 }
          }
        }
      ],
    },
  ]
}
