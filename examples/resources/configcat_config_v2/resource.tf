variable "product_id" {
  type = string
}

resource "configcat_config_v2" "my_config" {
  product_id  = var.product_id
  name        = "My config V2"
  description = "My config V2 description"
  order       = 0
}


output "config_id" {
  value = configcat_config_v2.my_config.id
}
