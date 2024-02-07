resource "configcat_permission_group" "my_permission_group" {
  product_id = data.configcat_products.my_products.products[0].product_id
  name       = "Read only except Test environment"

  accesstype = "custom"

  environment_access {
    environment_id         = data.configcat_environments.my_test_environments.environments[0].environment_id
    environment_accesstype = "full"
  }

  environment_access {
    environment_id         = data.configcat_environments.my_production_environments.environments[0].environment_id
    environment_accesstype = "none"
  }
}
