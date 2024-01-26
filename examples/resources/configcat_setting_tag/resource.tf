variable "setting_id" {
  type = string
}

variable "tag_id" {
  type = string
}

resource "configcat_setting_tag" "my_setting_tag" {
  setting_id = var.setting_id
  tag_id     = var.tag_id
}
