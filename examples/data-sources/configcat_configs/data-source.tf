variable "product_id" {
  type = string
}

data "configcat_configs" "my_configs" {
  product_id        = var.product_id
  name_filter_regex = "Main Config"
}


output "config_id" {
  value = data.configcat_configs.my_configs.configs.0.config_id
}
