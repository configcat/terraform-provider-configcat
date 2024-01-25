variable "name_filter_regex" {
  type    = string
  default = null
}

data "configcat_products" "test" {
  name_filter_regex = var.name_filter_regex
}
