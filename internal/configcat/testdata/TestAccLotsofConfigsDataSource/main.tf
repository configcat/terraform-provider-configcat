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

output "config_id" {
  value = length(data.configcat_configs.test[0].configs) > 0 ? data.configcat_configs.test[0].configs[0].config_id : null
}
