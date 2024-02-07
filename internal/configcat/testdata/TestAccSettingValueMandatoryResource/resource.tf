variable "config_id" {
  type = string
}

resource "configcat_setting" "test" {
  config_id = var.config_id
  key       = "MandatorySettingKey"
  name      = "MandatorySettingName"
  order     = 0
}
