variable "product_id" {
  type = string
}

variable "name" {
  type = string
}
variable "description" {
  type    = string
  default = null
}
variable "order" {
  type    = number
  default = 0
}

resource "configcat_config" "test" {
  product_id  = var.product_id
  name        = var.name
  description = var.description
  order       = var.order
}
