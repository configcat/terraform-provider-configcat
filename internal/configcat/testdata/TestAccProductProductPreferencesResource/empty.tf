resource "configcat_product" "product" {
  organization_id = "08d86d63-26dc-4276-86d6-eae122660e51"
  name            = "Product preferences test"
  order           = 1
}

resource "configcat_product_preferences" "preferences" {
  product_id = configcat_product.product.id
}
