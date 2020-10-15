# Simple usage of data sources

## Organization
```hcl
data "configcat_organizations" "organizations" {
  name_filter_regex = "ConfigCat"
}
```

## Products
```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}
```

## Configs
```hcl
data "configcat_configs" "configs" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Main Config"
}
```

## Environments
```hcl
data "configcat_environments" "environments" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Test"
}
```

## Settings
```hcl
data "configcat_settings" "settings" {
  config_id = data.configcat_configs.configs.configs.0.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}
```

## Tags
```hcl
data "configcat_tags" "tags" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Test"
}
```