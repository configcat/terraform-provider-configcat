variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "comparator_1" {
  type = string
}

variable "segment_id_1" {
  type = string
}
variable "comparator_2" {
  type = string
}

variable "segment_id_2" {
  type = string
}

resource "configcat_setting" "test" {
  config_id = var.config_id
  key       = "SegmentSettingV2Key"
  name      = "SegmentSettingV2Name"
  order     = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = false
  value          = { bool_value = true }


  targeting_rules = [
    {
      conditions = [
        {
          segment_condition = {
            comparator = var.comparator_1
            segment_id         = var.segment_id_1
          }
        }
      ]

      value = { bool_value = false }
    },
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "email"
            comparator           = "textEquals"
            comparison_value     = { string_value = "@configcat.com" }
          }
        },
        {
          segment_condition = {
            comparator = var.comparator_2
            segment_id         = var.segment_id_2
          }
        }
      ]

      value = { bool_value = true }
    }
  ]
}
