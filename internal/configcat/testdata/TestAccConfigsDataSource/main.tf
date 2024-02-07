variable "product_id" {
  type = string
}

variable "name_filter_regex" {
  type    = string
  default = null
}

data "configcat_configs" "test" {
  product_id        = var.product_id
  name_filter_regex = var.name_filter_regex
}


output "config_id" {
  value = length(data.configcat_configs.test.configs) > 0 ? data.configcat_configs.test.configs[0].config_id : null
}
