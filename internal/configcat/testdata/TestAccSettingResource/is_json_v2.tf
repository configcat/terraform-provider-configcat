variable "product_id" {
  type = string
}
variable "key" {
  type = string
}
variable "name" {
  type = string
}
variable "setting_type" {
  type = string
}
variable "order" {
  type    = number
  default = 0
}
variable "is_json" {
  type    = bool
  default = null
}

resource "configcat_config" "config" {
  product_id         = var.product_id
  name               = "Test config for setting is_json"
  evaluation_version = "v2"
  order              = 0
}

resource "configcat_setting" "test" {
  config_id    = configcat_config.config.id
  key          = var.key
  name         = var.name
  setting_type = var.setting_type
  order        = var.order
  is_json      = var.is_json
}
