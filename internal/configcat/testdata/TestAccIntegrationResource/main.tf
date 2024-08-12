
variable "product_id" {
  type    = string
  default = null
}
variable "integration_type" {
  type    = string
  default = null
}

variable "name" {
  type    = string
  default = null
}

variable "parameters" {
  type    = map(string)
  default = null
}

variable "configs" {
  type    = list(string)
  default = null
}

variable "environments" {
  type    = list(string)
  default = null
}

resource "configcat_integration" "test" {
  product_id = var.product_id

  integration_type = var.integration_type
  name             = var.name
  parameters       = var.parameters
  configs          = var.configs
  environments     = var.environments
}
