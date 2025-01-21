variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "value" {
  type = bool
}

variable "mandatory_notes" {
  type    = string
  default = null
}

resource "configcat_setting" "test" {
  config_id = var.config_id
  key       = "MandatorySettingV2Key"
  name      = "MandatorySettingV2Name"
  order     = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id  = var.environment_id
  setting_id      = configcat_setting.test.id
  init_only       = false
  value           = { bool_value = var.value }
  mandatory_notes = var.mandatory_notes
}
