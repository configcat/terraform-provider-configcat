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

variable "value" {
  type = string
}

variable "percentage1" {
  type    = number
  default = null
}
variable "percentage2" {
  type    = number
  default = null
}
variable "percentage3" {
  type    = number
  default = null
}

variable "percentage1value" {
  type    = string
  default = null
}
variable "percentage2value" {
  type    = string
  default = null
}
variable "percentage3value" {
  type    = string
  default = null
}

resource "configcat_setting" "test" {
  config_id    = var.config_id
  key          = "StringSettingKey"
  name         = "StringSettingName"
  setting_type = "string"
  order        = 0
}

resource "configcat_setting_value" "test" {
  environment_id = var.environment_id
  setting_id     = configcat_setting.test.id
  init_only      = var.init_only
  value          = var.value

  percentage_items {
    percentage = var.percentage1
    value      = var.percentage1value
  }
  percentage_items {
    percentage = var.percentage2
    value      = var.percentage2value
  }
  percentage_items {
    percentage = var.percentage3
    value      = var.percentage3value
  }
}
