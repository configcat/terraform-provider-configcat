# configcat_setting Resource

Use this data source to access information about an existing Feature Flag/Setting. [Used API](https://api.configcat.com/docs/index.html#operation/get-settings)

## Example Usage

```hcl
data "configcat_product" "product" {
  name = "ConfigCat's product"
}

data "configcat_config" "product" {
  product_id = configcat_product.product.id
  name = "Main Config"
}

data "configcat_setting" "setting" {
  config_id = configcat_config.config.id
  key = "isAwesomeFeatureEnabled"
}
```

## Argument Reference

* `config_id` - (Required) The ID of the Config
* `key` - (Required) The key of the Feature Flag/Setting

## Attribute Reference

* `id` - The unique Setting ID.
* `name` - The name of the Setting.
* `hint` - The hint of the Setting.
* `setting_type` - The Setting's type. Available values: `boolean`|`string`|`int`|`double`