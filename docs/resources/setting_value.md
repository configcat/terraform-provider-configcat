# configcat_setting_value Resource

Initializes/updates a Feature Flag/Setting's value.  
Used APIs:
* [Read](https://api.configcat.com/docs/#operation/get-setting-value)
* [Update](https://api.configcat.com/docs/#operation/replace-setting-value)

## Example Usage

```hcl
data "configcat_products" "products" {
  name = "ConfigCat's product"
}

data "configcat_configs" "configs" {
  product_id = configcat_products.products.products.0.id
  name = "Main Config"
}

data "configcat_environments" "environments" {
  product_id = configcat_products.products.products.0.id
  name = "Test"
}

data "configcat_settings" "settings" {
  config_id = configcat_configs.configs.configs.0.id
  key_filter_regex = "isAwesomeFeatureEnabled"
}

resource "configcat_setting_value" "setting_value" {
    environment_id = configcat_environments.environments.environments.0.id
    setting_id = configcat_settings.settings.settings.0.id
    setting_type = configcat_settings.settings.settings.0.setting_type
    
    value = "true"

    rollout_rules {
        comparison_attribute = "email"
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
* `environment_id` - (Required) The ID of the Environment
* `setting_id` - (Required) The ID of the Feature Flag/Setting
* `setting_type` - (Required) The Setting's type.
* `init_only` - (Optional) Default: true. Read more below  

The Feature Flag/Setting's value
* `value` - (Required) The Setting's value. Type: `string`. It must be compatible with the `setting_type`.
* `rollout_rules` - (Optional) A [list](https://www.terraform.io/docs/configuration/types.html#list-) to define [Rollout rules](https://configcat.com/docs/advanced/targeting/#anatomy-of-a-targeting-rule). Read more below
* `percentage_items` - (Optional) A [list](https://www.terraform.io/docs/configuration/types.html#list-) to define [Percentage items](https://configcat.com/docs/advanced/targeting/#targeting-a-percentage-of-users). Read more below

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

The `init_only` argument's default value is `true` which means that this resource's state will be only applied to the Feature Flag/Setting only once and only during resource creation.  
This prevents overriding the Feature Flag/Setting's modified values on the [ConfigCat Dashboard](https://app.configcat.com).  

If you want to fully manage the Feature Flag/Setting's value from Terraform, set `init_only` argument to `false`. After setting the`init_only` argument to `false` each terraform run will update the Feature Flag/Setting's value to the state provided in Terraform.
