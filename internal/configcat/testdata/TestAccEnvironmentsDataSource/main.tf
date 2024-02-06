variable "product_id" {
  type = string
}

variable "name_filter_regex" {
  type    = string
  default = null
}

data "configcat_environments" "test" {
  product_id        = var.product_id
  name_filter_regex = var.name_filter_regex
}

output "environment_id" {
  value = data.configcat_environments.test.environments[0].environment_id
}
