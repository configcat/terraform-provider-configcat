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
variable "color" {
  type    = string
  default = null
}
variable "order" {
  type    = number
  default = 0
}
variable "cleanup_auditlogs_on_destroy" {
  type    = bool
  default = false
}

resource "configcat_environment" "test" {
  product_id  = var.product_id
  name        = var.name
  description = var.description
  color       = var.color
  order       = var.order
  cleanup_auditlogs_on_destroy = var.cleanup_auditlogs_on_destroy
}
