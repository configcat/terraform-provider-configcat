variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "list_value" {
  type = string
}

resource "configcat_setting" "test" {
  config_id    = var.config_id
  key          = "ListValue"
  name         = "ListValue"
  setting_type = "boolean"
  order        = 0
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
          user_condition = {
            comparison_attribute = "Identifier"
            comparator           = "sensitiveIsOneOf"
            comparison_value = {
              list_values = [
                { value = "A" },
              ]
            }
          }
        }
      ]
      value = { bool_value = true }
    },
    {
      conditions = [
        {
          user_condition = {
            comparison_attribute = "Identifier"
            comparator           = "sensitiveIsOneOf"
            comparison_value = {
              list_values = [
                { value = var.list_value },
              ]
            }
          }
        }
      ]
      value = { bool_value = true }
    }
  ]
}
