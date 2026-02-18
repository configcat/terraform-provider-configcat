
variable "product_id" {
  type = string
}
variable "evaluation_version" {
  type    = string
  default = "v2"
}
variable "setting_type" {
  type = string
}

variable "predefined_variations" {
  type = list(object({
    value = object({
      bool_value   = optional(bool, null)
      string_value = optional(string, null)
      int_value    = optional(number, null)
      double_value = optional(number, null)
    })
    name = optional(string, null)
    hint = optional(string, null)
  }))
}

resource "configcat_config" "config" {
  product_id         = var.product_id
  name               = "Test config for predefined variations (setting type: ${var.setting_type})"
  evaluation_version = var.evaluation_version
  order = 0
}

resource "configcat_setting" "test" {
  config_id             = configcat_config.config.id
  key                   = var.setting_type
  name                  = var.setting_type
  setting_type          = var.setting_type
  predefined_variations = var.predefined_variations
  order = 0
}
