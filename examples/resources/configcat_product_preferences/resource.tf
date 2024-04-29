variable "organization_id" {
  type = string
}

resource "configcat_product" "product" {
  organization_id = var.organization_id
  name            = "My product"
  order           = 0
}

resource "configcat_environment" "test" {
  product_id = configcat_product.product.id
  name       = "Test"
  order      = 0
}

resource "configcat_environment" "production" {
  product_id = configcat_product.product.id
  name       = "Production"
  order      = 1
}

resource "configcat_product_preferences" "preferences" {
  product_id = configcat_product.product.id

  key_generation_mode    = "kebabCase"
  mandatory_setting_hint = true
  show_variation_id      = false
  reason_required        = false
  reason_required_environments = {
    "${configcat_environment.test.id}"       = false,
    "${configcat_environment.production.id}" = true,
  }
}
