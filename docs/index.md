# ConfigCat Feature Flags Provider

Manage features and change your software configuration using [ConfigCat feature flags](https://configcat.com), without the need to re-deploy code.  
A 10 minute trainable dashboard allows even non-technical team members to manage application features.  
Supports A/B testing, soft launching or targeting a specific group of users first with new ideas. Deploy any time, release when confident.  
Open-source SDKs enable easy integration with any web, mobile or backend application.

ConfigCat Feature Flags Provider allows you to configure and access ConfigCat resources via [ConfigCat Public Management API](https://api.configcat.com/). 

## Authentication

ConfigCat Feature Flags Provider requires authentication with [ConfigCat Public API credentials](https://app.configcat.com/my-account/public-api-credentials).

## Provider registration

```hcl
terraform {
  required_providers {
    configcat = {
      source = "configcat/configcat"
      version = "~> 1.0"
    }
  }
}

provider "configcat" {
  version     = "~> 1.0"

  // Get your ConfigCat Public API credentials at https://app.configcat.com/my-account/public-api-credentials
  basic_auth_username = var.configcat_basic_auth_username
  basic_auth_password = var.configcat_basic_auth_password
}
```

## Example usage

```hcl
// Retrieve your Product
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

// Retrieve your Config
data "configcat_configs" "configs" {
  product_id = data.configcat_products.products.products.0.product_id
  name_filter_regex = "Main Config"
}

// Create a Feature Flag/Setting
resource "configcat_setting" "setting" {
  config_id = data.configcat_configs.configs.configs.0.config_id
  key = "isAwesomeFeatureEnabled"
  name = "My awesome feature flag"
  hint = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
}

// Initialize the Feature Flag/Setting's value
resource "configcat_setting_value" "setting_value" {
    environment_id = data.configcat_environments.environments.environments.0.environment_id
    setting_id = configcat_setting.setting.id
    value = "false"
}
```

## Argument Reference

The following arguments are supported:

* `basic_auth_username` - (Required) Get your `basic_auth_username` at [ConfigCat Public API credentials](https://app.configcat.com/my-account/public-api-credentials).  
This can also be sourced from the `CONFIGCAT_BASIC_AUTH_USERNAME` Environment Variable.

* `basic_auth_password` - (Required) Get your `basic_auth_password` at [ConfigCat Public API credentials](https://app.configcat.com/my-account/public-api-credentials).  
This can also be sourced from the `CONFIGCAT_BASIC_AUTH_PASSWORD` Environment Variable.

* `base_path` - (Optional) ConfigCat Public Management API's `base_path`. Defaults to https://api.configcat.com.  
This can also be sourced from the `CONFIGCAT_BASE_PATH` Environment Variable.
