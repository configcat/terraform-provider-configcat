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
  key       = "SegmentSettingKey"
  name      = "SegmentSettingName"
  order     = 0
}

resource "configcat_setting_value" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = false
  value          = "true"

  rollout_rules {
    segment_comparator = var.comparator_1
    segment_id         = var.segment_id_1
    value              = "false"
  }

  rollout_rules {
    comparison_attribute = "email"
    comparator           = "contains"
    comparison_value     = "@configcat.com"
    value                = "true"
  }

  rollout_rules {
    segment_comparator = var.comparator_2
    segment_id         = var.segment_id_2
    value              = "false"
  }
}
