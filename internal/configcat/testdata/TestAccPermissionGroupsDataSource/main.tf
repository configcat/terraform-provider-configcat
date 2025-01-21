variable "product_id" {
  type = string
}

variable "name_filter_regex" {
  type    = string
  default = null
}

data "configcat_permission_groups" "test" {
  product_id        = var.product_id
  name_filter_regex = var.name_filter_regex
}

output "permission_group_id" {
  value = length(data.configcat_permission_groups.test.permission_groups) > 0 ? data.configcat_permission_groups.test.permission_groups[0].permission_group_id : null
}
