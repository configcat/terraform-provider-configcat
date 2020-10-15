# configcat_setting_tag Resource

Manages **Feature Flag/Setting**'s **Tags**.  

## Example Usage

```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_configs" "configs" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Main Config"
}

data "configcat_tags" "tags" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Tag"
}

data "configcat_settings" "settings" {
  config_id = data.configcat_configs.configs.configs.0.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}

resource "configcat_setting_tag" "setting_tag" {
    setting_id = data.configcat_settings.settings.settings.0.setting_id
    tag_id = data.configcat_tags.tags.tags.0.tag_id
}
```

## Argument Reference

Parameters
* `setting_id` - (Required) The ID of the Feature Flag/Setting.
* `tag_id` - (Required) The ID of the Tag.

## Import

Tags can be imported using a combined SettingId:TagId ID.  
Get the SettingId using e.g. the [GetSettings API](https://api.configcat.com/docs/#operation/get-settings).  
Get the TagId using e.g. the [GetTags API](https://api.configcat.com/docs/#operation/get-tags).  

```
$ terraform import configcat_setting_tag.example 1234:5678
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Used APIs
* [Get Setting](https://api.configcat.com/docs/#operation/get-setting)
* [Update Setting](https://api.configcat.com/docs/#operation/update-setting)
