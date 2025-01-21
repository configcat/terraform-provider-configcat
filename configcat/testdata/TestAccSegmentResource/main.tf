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

variable "comparison_attribute" {
  type = string
}

variable "comparator" {
  type = string
}

variable "comparison_value" {
  type = string
}

resource "configcat_segment" "test" {
  product_id           = var.product_id
  name                 = var.name
  description          = var.description
  comparison_attribute = var.comparison_attribute
  comparator           = var.comparator
  comparison_value     = var.comparison_value
}
