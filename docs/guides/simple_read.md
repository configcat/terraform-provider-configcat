# Simple usage of data sources


```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_configs" "configs" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Main Config"
}

data "configcat_settings" "settings" {
  config_id = data.configcat_configs.configs.configs.0.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}


output "setting_id" {
  value = data.configcat_settings.settings.settings.0.setting_id
}
```