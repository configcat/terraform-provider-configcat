---
page_title: "Simple usage of Data Sources"
---

# Simple usage of Data Sources

## Prerequisites

[Get your Public Management API credentials](https://app.configcat.com/my-account/public-api-credentials) and set the following environment variables:

- CONFIGCAT_BASIC_AUTH_USERNAME
- CONFIGCAT_BASIC_AUTH_PASSWORD

## main.tf

```terraform
terraform {
  required_providers {
    configcat = {
      source  = "configcat/configcat"
      version = "~> 4.0"
    }
  }
}

provider "configcat" {
}

data "configcat_organizations" "my_organizations" {
  name_filter_regex = "ConfigCat"
}

output "organization_id" {
  value = data.configcat_organizations.my_organizations.organizations[0].organization_id
}

data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_permission_groups" "my_permission_groups" {
  product_id        = data.configcat_products.my_products.products[0].product_id
  name_filter_regex = "Administrators"
}

output "permission_group_id" {
  value = data.configcat_permission_groups.my_permission_groups.permission_groups[0].permission_group_id
}

data "configcat_configs" "my_configs" {
  product_id        = data.configcat_products.my_products.products[0].product_id
  name_filter_regex = "Main Config"
}

data "configcat_environments" "my_environments" {
  product_id        = data.configcat_products.my_products.products[0].product_id
  name_filter_regex = "Test"
}

data "configcat_sdkkeys" "my_sdkkeys" {
  config_id      = data.configcat_configs.my_configs.configs[0].config_id
  environment_id = data.configcat_environments.my_environments.environments[0].environment_id
}

output "primary_sdkkey" {
  value = data.configcat_sdkkeys.my_sdkkeys.primary
}

data "configcat_settings" "my_settings" {
  config_id        = data.configcat_configs.my_configs.configs[0].config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}

output "setting_id" {
  value = data.configcat_settings.my_settings.settings[0].setting_id
}

data "configcat_tags" "my_tags" {
  product_id        = data.configcat_products.my_products.products[0].product_id
  name_filter_regex = "Test"
}

output "tag_id" {
  value = data.configcat_tags.my_tags.tags[0].tag_id
}
```
