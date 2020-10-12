# configcat_setting Resource

Creates and manages a Feature Flag/Setting.  
Used APIs:
* [Read](https://api.configcat.com/docs/index.html#operation/get-setting)
* [Create](https://api.configcat.com/docs/index.html#operation/create-setting)
* [Update](https://api.configcat.com/docs/index.html#operation/update-setting)
* [Delete](https://api.configcat.com/docs/index.html#operation/delete-setting)

## Example Usage

```hcl
data "configcat_product" "product" {
  name = "ConfigCat's product"
}

data "configcat_config" "product" {
  product_id = configcat_product.product.id
  name = "Main Config"
}

resource "configcat_setting" "setting" {
  config_id = configcat_config.config.id
  key = "isAwesomeFeatureEnabled"
  name = "My awesome feature flag"
  hint = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
}
```

## Argument Reference

* `config_id` - (Required) The ID of the Config
* `key` - (Required) The key of the Feature Flag/Setting
* `name` - (Required) The name of the Setting.
* `hint` - (Optional) The hint of the Setting.
* `setting_type` - (Optional) Default: `boolean`. The Setting's type. Available values: `boolean`|`string`|`int`|`double`

## Attribute Reference

* `id` - The unique Setting ID.
