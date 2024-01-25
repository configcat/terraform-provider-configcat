variable "product_id" {
  type = string
}

resource "configcat_environment" "my_environment" {
  product_id  = var.product_id
  name        = "Staging"
  description = "Staging description"
  color       = "blue"
  order       = 0
}


output "environment_id" {
  value = configcat_environment.my_environment.id
}
