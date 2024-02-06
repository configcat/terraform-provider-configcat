variable "product_id" {
  type = string
}

data "configcat_permission_groups" "my_permission_groups" {
  product_id        = var.product_id
  name_filter_regex = "Administrators"
}


output "permission_group_id" {
  value = data.configcat_permission_groups.my_permission_groups.permission_groups[0].permission_group_id
}
