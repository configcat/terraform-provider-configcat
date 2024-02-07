variable "product_id" {
  type = string
}

variable "name" {
  type = string
}
variable "color" {
  type    = string
  default = null
}


resource "configcat_tag" "test" {
  product_id = var.product_id
  name       = var.name
  color      = var.color
}
