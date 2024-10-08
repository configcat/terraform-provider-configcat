variable "product_id" {
  type = string
}

variable "name" {
  type = string
}
variable "accesstype" {
  type    = string
  default = null
}
variable "new_environment_accesstype" {
  type    = string
  default = null
}


variable "can_manage_members" {
  type    = bool
  default = null
}
variable "can_createorupdate_config" {
  type    = bool
  default = null
}
variable "can_delete_config" {
  type    = bool
  default = null
}
variable "can_createorupdate_environment" {
  type    = bool
  default = null
}
variable "can_delete_environment" {
  type    = bool
  default = null
}
variable "can_createorupdate_setting" {
  type    = bool
  default = null
}
variable "can_tag_setting" {
  type    = bool
  default = null
}
variable "can_delete_setting" {
  type    = bool
  default = null
}
variable "can_createorupdate_tag" {
  type    = bool
  default = null
}
variable "can_delete_tag" {
  type    = bool
  default = null
}
variable "can_manage_webhook" {
  type    = bool
  default = null
}
variable "can_use_exportimport" {
  type    = bool
  default = null
}
variable "can_manage_product_preferences" {
  type    = bool
  default = null
}
variable "can_manage_integrations" {
  type    = bool
  default = null
}
variable "can_view_sdkkey" {
  type    = bool
  default = null
}
variable "can_rotate_sdkkey" {
  type    = bool
  default = null
}
variable "can_createorupdate_segment" {
  type    = bool
  default = null
}
variable "can_delete_segment" {
  type    = bool
  default = null
}
variable "can_view_product_auditlog" {
  type    = bool
  default = null
}
variable "can_view_product_statistics" {
  type    = bool
  default = null
}
variable "can_disable_2fa" {
  type    = bool
  default = null
}
variable "environment_accesses" {
  type    = map(string)
  default = null
}

resource "configcat_permission_group" "test" {
  product_id                     = var.product_id
  name                           = var.name
  accesstype                     = var.accesstype
  new_environment_accesstype     = var.new_environment_accesstype
  environment_accesses           = var.environment_accesses
  can_manage_members             = var.can_manage_members
  can_createorupdate_config      = var.can_createorupdate_config
  can_delete_config              = var.can_delete_config
  can_createorupdate_environment = var.can_createorupdate_environment
  can_delete_environment         = var.can_delete_environment
  can_createorupdate_setting     = var.can_createorupdate_setting
  can_tag_setting                = var.can_tag_setting
  can_delete_setting             = var.can_delete_setting
  can_createorupdate_tag         = var.can_createorupdate_tag
  can_delete_tag                 = var.can_delete_tag
  can_manage_webhook             = var.can_manage_webhook
  can_use_exportimport           = var.can_use_exportimport
  can_manage_product_preferences = var.can_manage_product_preferences
  can_manage_integrations        = var.can_manage_integrations
  can_view_sdkkey                = var.can_view_sdkkey
  can_rotate_sdkkey              = var.can_rotate_sdkkey
  can_createorupdate_segment     = var.can_createorupdate_segment
  can_delete_segment             = var.can_delete_segment
  can_view_product_auditlog      = var.can_view_product_auditlog
  can_view_product_statistics    = var.can_view_product_statistics
  can_disable_2fa                = var.can_disable_2fa
}
