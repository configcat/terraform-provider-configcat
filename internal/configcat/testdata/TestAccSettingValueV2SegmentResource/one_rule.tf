variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "comparator" {
  type = string
}

variable "segment_id" {
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
  value          = { bool_value = false }

  targeting_rules = [
    {
      conditions = [
        {
          segment_condition = {
            comparator = var.comparator
            segment_id         = var.segment_id
          }
        }
      ]

      value = { bool_value = false }
    }
  ]
}

