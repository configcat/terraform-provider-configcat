# configcat_settings Resource

Use this data source to access information about an existing Feature Flag/Setting.  

## Example Usage

```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_configs" "configs" {
  product_id = data.configcat_product.products.products.0.product_id
  name_filter_regex = "Main Config"
}

data "configcat_settings" "settings" {
  config_id = data.configcat_config.configs.configs.0.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}


output "setting_id" {
  value = data.configcat_settings.settings.settings.0.setting_id
}
```

## Argument Reference

* `config_id` - (Required) The ID of the Config.
* `key_filter_regex` - (Optional) Filter the Settings by key.

## Attribute Reference

* `settings` - A setting [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `settings` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `setting_id` - The unique Setting ID.
* `key` - The key of the Feature Flag/Setting.
* `name` - The name of the Setting.
* `hint` - The hint of the Setting.
* `setting_type` - The Setting's type. Available values: `boolean`|`string`|`int`|`double`.

## Used APIs:
- [Read](https://api.configcat.com/docs/index.html#operation/get-settings)