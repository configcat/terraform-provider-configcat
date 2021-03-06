# configcat_setting_value Resource

Initializes and updates **Feature Flag and Setting** values. [Read more about the anatomy of a Feature Flag or Setting.](https://configcat.com/docs/main-concepts) 

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

data "configcat_settings" "my_settings" {
  config_id = data.configcat_configs.my_configs.configs.0.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}

resource "configcat_setting_value" "my_setting_value" {
    environment_id = data.configcat_environments.my_environments.environments.0.environment_id
    setting_id = data.configcat_settings.my_settings.settings.0.setting_id

    mandatory_notes = "mandatory notes"
    
    value = "true"

    rollout_rules {
        comparison_attribute = "Email"
        comparator = "contains"
        comparison_value = "@mycompany.com"
        value = "true"
    }
    rollout_rules {
        comparison_attribute = "custom"
        comparator = "isOneOf"
        comparison_value = "red"
        value = "false"
    }

    percentage_items {
        percentage = 20
        value = "true"
    }
    percentage_items {
        percentage = 80
        value = "false"
    }
}
```

## Argument Reference

Parameters
* `environment_id` - (Required) The ID of the Environment.
* `setting_id` - (Required) The ID of the Feature Flag/Setting.
* `setting_type` - (Required) The Setting's type.
* `mandatory_notes` - (Optional) Default: "". If the Product's "Mandatory notes" preference is turned on for the Environment the Mandatory note must be passed.  
* `init_only` - (Optional) Default: true. Read more below.  

The Feature Flag/Setting's value
* `value` - (Required) The Setting's value. Type: `string`. It must be compatible with the `setting_type`.
* `rollout_rules` - (Optional) A [list](https://www.terraform.io/docs/configuration/types.html#list-) to define [Rollout rules](https://configcat.com/docs/advanced/targeting/#anatomy-of-a-targeting-rule). Read more below.
* `percentage_items` - (Optional) A [list](https://www.terraform.io/docs/configuration/types.html#list-) to define [Percentage items](https://configcat.com/docs/advanced/targeting/#targeting-a-percentage-of-users). Read more below.

### `rollout_rules` list

By adding a rule, you specify a group of your users and what feature flag - or other settings - value they should get.  

* `comparison_attribute` - (Required) The [comparison attribute](https://configcat.com/docs/advanced/targeting/#attribute).
* `comparator` - (Required) The [comparator](https://configcat.com/docs/advanced/targeting/#comparator).
* `comparison_value` - (Required) The [comparison value](https://configcat.com/docs/advanced/targeting/#comparison-value).
* `value` - (Required) The exact [value](https://configcat.com/docs/advanced/targeting/#served-value) that will be served to the users who match the targeting rule. Type: `string`. It must be compatible with the `setting_type`.

### `percentage_items` list

With percentage-based user targeting, you can specify a randomly selected fraction of your users whom a feature will be enabled or a different value will be served.

* `percentage` - (Required) Any [number](https://configcat.com/docs/advanced/targeting/#-value) between 0 and 100 that represents a randomly allocated fraction of your users.
* `value` - (Required) The exact [value](https://configcat.com/docs/advanced/targeting/#served-value-1) that will be served to the users that fall into that fraction. Type: `string`. It must be compatible with the `setting_type`.

### `init_only` argument

The main purpose of this resource to provide an initial value for the Feature Flag/Setting.  

The `init_only` argument's default value is `true`. Meaning that the Feature Flag or Setting's **value will be only be applied once** during resource creation. If someone modifies the value on the [ConfigCat Dashboard](https://app.configcat.com) those modifications will **not be overwritten** by the Terraform script.

If you want to fully manage the Feature Flag/Setting's value from Terraform, set `init_only` argument to `false`. After setting the`init_only` argument to `false` each terraform run will update the Feature Flag/Setting's value to the state provided in Terraform.

## Import

Feature Flag/Setting values can be imported using a combined EnvironmentID:SettingId ID.  
Get the SettingId using e.g. the [GetSettings API](https://api.configcat.com/docs/#operation/get-settings).  
Get the EnvironmentId using e.g. the [GetEnvironments API](https://api.configcat.com/docs/#operation/get-environments).

```
$ terraform import configcat_setting_value.example 08d86d63-2726-47cd-8bfc-59608ecb91e2:1234
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Setting Value](https://api.configcat.com/docs/#operation/get-setting-value)
* [Replace Setting Value](https://api.configcat.com/docs/#operation/replace-setting-value)
