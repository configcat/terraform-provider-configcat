variable "organization_id" {
  type = string
}

resource "configcat_product" "my_config" {
  organization_id = var.organization_id
  name            = "My product"
  description     = "My product description"
  order           = 0
}


output "product_id" {
  value = configcat_product.my_product.id
}
