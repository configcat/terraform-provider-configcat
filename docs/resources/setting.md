# configcat_setting Resource

Creates and manages a **Feature Flag/Setting**. [Read more about the anatomy of a Feature Flag or Setting.](https://configcat.com/docs/main-concepts) 

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_configs" "my_configs" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Main Config"
}

resource "configcat_setting" "my_setting" {
  config_id = data.configcat_configs.my_configs.configs.0.config_id
  key = "isAwesomeFeatureEnabled"
  name = "My awesome feature flag"
  hint = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
}


output "setting_id" {
  value = configcat_setting.my_setting.id
}
```

## Argument Reference

* `config_id` - (Required) The ID of the Config.
* `key` - (Required) The key of the Feature Flag/Setting.
* `name` - (Required) The name of the Setting.
* `hint` - (Optional) The hint of the Setting.
* `setting_type` - (Optional) Default: `boolean`. The Setting's type.  
Available values: `boolean`|`string`|`int`|`double`.

## Attribute Reference

* `id` - The unique Setting ID.

## Import

Feature Flags/Settings can be imported using the SettingId. Get the SettingId using e.g. the [List Flags API](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/get-settings).

```
$ terraform import configcat_setting.example 1234
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Flag](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/get-setting)
* [Create Flag](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/create-setting)
* [Update Flag](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/update-setting)
* [Delete Flag](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/delete-setting)
