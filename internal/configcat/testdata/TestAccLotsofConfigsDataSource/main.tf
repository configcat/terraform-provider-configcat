variable "product_id" {
  type = string
}

variable "name_filter_regex" {
  type    = string
  default = null
}

variable "data_sources" {
  type = number
}

data "configcat_configs" "test" {
  count = var.data_sources

  product_id        = var.product_id
  name_filter_regex = var.name_filter_regex
}
