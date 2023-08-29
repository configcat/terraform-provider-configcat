---
page_title: "Simple usage of Resources"
---

# Simple usage of Resources

## Prerequisites

[Get your Public Management API credentials](https://app.configcat.com/my-account/public-api-credentials) and set the following environment variables:
- CONFIGCAT_BASIC_AUTH_USERNAME
- CONFIGCAT_BASIC_AUTH_PASSWORD

## root.tf

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
}

// Organization Resource is ReadOnly.
data "configcat_organizations" "my_organizations" {
  name_filter_regex = "ConfigCat"
}

resource "configcat_product" "my_product" {
  organization_id = data.configcat_organizations.my_organizations.organizations.0.organization_id
  name = "My product"
  description = "My product description"
}

resource "configcat_permission_group" "my_permission_group" {
  product_id = configcat_product.my_product.id
  name = "Administrators"
}

resource "configcat_config" "my_config" {
  product_id = configcat_product.my_product.id
  name = "My config"
  description = "My config description"
}

resource "configcat_environment" "my_environment" {
  product_id = configcat_product.my_product.id
  name = "Production"
  description = "Production description"
  color = "blue"
}

resource "configcat_setting" "is_awesome" {
  config_id = configcat_config.my_config.id
  key = "isAwesomeFeatureEnabled"
  name = "My awesome feature flag"
  hint = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
}

resource "configcat_setting_value" "is_awesome_value" {
    environment_id = configcat_environment.my_environment.id
    setting_id = configcat_setting.is_awesome.id
    
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

resource "configcat_tag" "my_tag" {
  product_id = configcat_product.my_product.id
  name = "Created by Terraform"
}

resource "configcat_setting_tag" "is_awesome_tag" {
    setting_id = configcat_setting.is_awesome.id
    tag_id = configcat_tag.my_tag.id
}
```
