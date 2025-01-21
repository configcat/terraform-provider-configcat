variable "config_id" {
  type = string
}

resource "configcat_setting" "test" {
  config_id = var.config_id
  key       = "MandatorySettingV2Key"
  name      = "MandatorySettingV2Name"
  order     = 0
}
