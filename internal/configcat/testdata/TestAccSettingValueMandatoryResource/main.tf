variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "value" {
  type = string
}

variable "mandatory_notes" {
  type    = string
  default = null
}

resource "configcat_setting" "test" {
  config_id = var.config_id
  key       = "MandatorySettingKey"
  name      = "MandatorySettingName"
  order     = 0
}

resource "configcat_setting_value" "test" {
  environment_id  = var.environment_id
  setting_id      = configcat_setting.test.id
  init_only       = false
  value           = var.value
  mandatory_notes = var.mandatory_notes
}
