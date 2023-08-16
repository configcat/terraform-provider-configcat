# configcat_sdkkeys Resource

Use this data source to access information about **SDK Keys**.
## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}
data "configcat_configs" "my_configs" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Main Config"
}
data "configcat_environments" "my_environments" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Test"
}

data "configcat_sdkkeys" "my_sdkkey" {
  config_id = data.configcat_configs.my_configs.configs.0.config_id
  environment_id = data.configcat_environments.my_environments.environments.0.environment_id
}

output "primary_sdkkey" {
  value = data.configcat_sdkkeys.my_sdkkey.primary
}

output "secondary_sdkkey" {
  value = data.configcat_sdkkeys.my_sdkkey.secondary
}
```

## Argument Reference

* `config_id` - (Required) The ID of the Config.
* `environment_id` - (Required) The ID of the Environment.

## Attribute Reference

* `primary` - The primary SDK Key associated with your **Config** and **Environment**.
* `secondary` - The secondary SDK Key associated with your **Config** and **Environment**.

## Endpoints used
- [Get SDK Key](https://api.configcat.com/docs/#tag/SDK-Keys/operation/get-sdk-keys)
