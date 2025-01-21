variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

resource "configcat_setting" "bool" {
  config_id = var.config_id
  key       = "BoolDependencySettingV2"
  name      = "BoolDependencySettingV2"
  order     = 1
}
resource "configcat_setting" "string" {
  config_id    = var.config_id
  key          = "StringDependencySettingV2"
  name         = "StringDependencySettingV2"
  order        = 2
  setting_type = "string"
}
resource "configcat_setting" "int" {
  config_id    = var.config_id
  key          = "IntDependencySettingV2"
  name         = "IntDependencySettingV2"
  order        = 3
  setting_type = "int"
}
resource "configcat_setting" "double" {
  config_id    = var.config_id
  key          = "DoubleDependencySettingV2"
  name         = "DoubleDependencySettingV2"
  order        = 4
  setting_type = "double"
}


