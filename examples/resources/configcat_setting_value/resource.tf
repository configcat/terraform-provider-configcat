variable "environment_id" {
  type = string
}

variable "setting_id" {
  type = string
}

resource "configcat_setting_value" "my_setting_value" {
  environment_id = var.environment_id
  setting_id     = var.setting_id

  mandatory_notes = "mandatory notes"

  value = "true"

  rollout_rules {
    comparison_attribute = "Email"
    comparator           = "contains"
    comparison_value     = "@mycompany.com"
    value                = "true"
  }
  rollout_rules {
    comparison_attribute = "custom"
    comparator           = "isOneOf"
    comparison_value     = "red"
    value                = "false"
  }

  percentage_items {
    percentage = 20
    value      = "true"
  }
  percentage_items {
    percentage = 80
    value      = "false"
  }
}
