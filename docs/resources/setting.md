# configcat_setting Resource

Creates and manages a Feature Flag/Setting.  

## Example Usage

```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_configs" "configs" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Main Config"
}

resource "configcat_setting" "setting" {
  config_id = data.configcat_configs.configs.configs.0.config_id
  key = "isAwesomeFeatureEnabled"
  name = "My awesome feature flag"
  hint = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
}


output "setting_id" {
  value = configcat_setting.setting.id
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

## Import

Feature Flags/Settings can be imported using the SettingId. Get the SettingId using e.g. the [GetSettings API](https://api.configcat.com/docs/#operation/get-settings).

```
$ terraform import configcat_setting.example 1234
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Used APIs
* [Get Setting](https://api.configcat.com/docs/index.html#operation/get-setting)
* [Create Setting](https://api.configcat.com/docs/index.html#operation/create-setting)
* [Update Setting](https://api.configcat.com/docs/index.html#operation/update-setting)
* [Delete Setting](https://api.configcat.com/docs/index.html#operation/delete-setting)
