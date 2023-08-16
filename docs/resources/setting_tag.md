# configcat_setting_tag Resource

Adds/Removes **Tags** to/from **Feature Flags and Settings**.

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_configs" "my_configs" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Main Config"
}

data "configcat_tags" "my_tags" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Tag"
}

data "configcat_settings" "my_settings" {
  config_id = data.configcat_configs.my_configs.configs.0.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}

resource "configcat_setting_tag" "my_setting_tag" {
    setting_id = data.configcat_settings.my_settings.settings.0.setting_id
    tag_id = data.configcat_tags.my_tags.tags.0.tag_id
}
```

## Argument Reference

Parameters
* `setting_id` - (Required) The ID of the Feature Flag/Setting.
* `tag_id` - (Required) The ID of the Tag.

## Import

Tags can be imported using a combined SettingId:TagId ID.  
Get the SettingId using e.g. the [List Flags API](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/get-settings).  
Get the TagId using e.g. the [List Tags API](https://api.configcat.com/docs/#tag/Tags/operation/get-tags).  

```
$ terraform import configcat_setting_tag.example 1234:5678
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Flag](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/get-setting)
* [Update Flag](https://api.configcat.com/docs/#tag/Feature-Flags-and-Settings/operation/update-setting)
