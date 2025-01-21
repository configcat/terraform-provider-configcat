variable "config_id" {
  type = string
}
variable "key" {
  type = string
}
variable "name" {
  type = string
}
variable "hint" {
  type    = string
  default = null
}
variable "setting_type" {
  type = string
}
variable "order" {
  type    = number
  default = 0
}

resource "configcat_setting" "test" {
  config_id    = var.config_id
  key          = var.key
  name         = var.name
  hint         = var.hint
  setting_type = var.setting_type
  order        = var.order
}
