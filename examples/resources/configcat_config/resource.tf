variable "product_id" {
  type = string
}

resource "configcat_config" "my_config" {
  product_id  = var.product_id
  name        = "My config"
  description = "My config description"
  order       = 0
}


output "config_id" {
  value = configcat_config.my_config.id
}
