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
