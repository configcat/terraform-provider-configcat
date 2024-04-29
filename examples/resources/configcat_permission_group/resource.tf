variable "product_id" {
  type = string
}

variable "test_environment_id" {
  type = string
}

variable "productuction_environment_id" {
  type = string
}

resource "configcat_permission_group" "admin_permission_group" {
  product_id = var.product_id
  name       = "Administrators"

  accesstype = "full"

  can_manage_members             = true
  can_createorupdate_config      = true
  can_delete_config              = true
  can_createorupdate_environment = true
  can_delete_environment         = true
  can_createorupdate_setting     = true
  can_tag_setting                = true
  can_delete_setting             = true
  can_createorupdate_tag         = true
  can_delete_tag                 = true
  can_manage_webhook             = true
  can_use_exportimport           = true
  can_manage_product_preferences = true
  can_manage_integrations        = true
  can_view_sdkkey                = true
  can_rotate_sdkkey              = true
  can_createorupdate_segment     = true
  can_delete_segment             = true
  can_view_product_auditlog      = true
  can_view_product_statistics    = true
}


resource "configcat_permission_group" "custom_permission_group" {
  product_id = var.product_id
  name       = "Read only except Test environment"

  accesstype = "custom"

  environment_accesses = {
    (var.test_environment_id)          = "full"
    (var.productuction_environment_id) = "readOnly"
  }
}


output "admin_permission_group_id" {
  value = configcat_permission_group.admin_permission_group.id
}

output "custom_permission_group_id" {
  value = configcat_permission_group.custom_permission_group.id
}
