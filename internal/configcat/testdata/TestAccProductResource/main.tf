variable "organization_id" {
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

resource "configcat_product" "test" {
  organization_id = var.organization_id
  name            = var.name
  description     = var.description
  order           = var.order
}
