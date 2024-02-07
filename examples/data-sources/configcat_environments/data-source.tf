variable "product_id" {
  type = string
}

data "configcat_environments" "my_environments" {
  product_id        = var.product_id
  name_filter_regex = "Test"
}


output "environment_id" {
  value = data.configcat_environments.my_environments.environments[0].environment_id
}
