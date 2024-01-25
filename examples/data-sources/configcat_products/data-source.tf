data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}


output "product_id" {
  value = data.configcat_products.my_products.products.0.product_id
}
