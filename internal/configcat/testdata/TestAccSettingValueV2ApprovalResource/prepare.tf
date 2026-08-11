resource "configcat_product" "test" {
  organization_id = "08d86d63-26dc-4276-86d6-eae122660e51"
  name            = "Approval test product"
  order           = 99
}

resource "configcat_product_preferences" "test" {
  product_id       = configcat_product.test.id
  approve_required = true
}

resource "configcat_config" "test" {
  product_id         = configcat_product.test.id
  name               = "Approval test config"
  evaluation_version = "v2"
  order              = 0
}

resource "configcat_environment" "test" {
  product_id = configcat_product.test.id
  name       = "Approval test environment"
  order      = 0
}

resource "configcat_setting" "test" {
  config_id = configcat_config.test.id
  key       = "ApprovalTestBoolKey"
  name      = "Approval test bool flag"
  order     = 0
}
