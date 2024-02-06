---
page_title: "Advanced usage of Resources with Segments"
---

# Advanced usage of Resources with Segments

Read more about [segmentation and segments](https://configcat.com/docs/advanced/segments) in ConfigCat.

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

# Organization Resource is ReadOnly.
data "configcat_organizations" "my_organizations" {
  name_filter_regex = "ConfigCat"
}

resource "configcat_product" "my_product" {
  organization_id = data.configcat_organizations.my_organizations.organizations[0].organization_id
  name            = "My product"
  description     = "My product description"
  order           = 0
}

resource "configcat_config" "my_config" {
  product_id  = configcat_product.my_product.id
  name        = "My config"
  description = "My config description"
  order       = 0
}

resource "configcat_environment" "my_environment" {
  product_id  = configcat_product.my_product.id
  name        = "Production"
  description = "Production description"
  color       = "blue"
  order       = 0
}

resource "configcat_segment" "sensitive_users" {
  product_id           = configcat_product.my_product.id
  name                 = "Sensitive users"
  description          = "Exclude these users from beta testings."
  comparison_attribute = "email"
  comparator           = "sensitiveIsOneOf"
  comparison_value     = "user@sensitivecompany.com,user2@sensitivecompany.com"
}

resource "configcat_segment" "dogfooding" {
  product_id           = configcat_product.my_product.id
  name                 = "Dogfooding"
  description          = "Eat your own dog food."
  comparison_attribute = "email"
  comparator           = "contains"
  comparison_value     = "@mycompany.com"
}

resource "configcat_setting" "is_awesome" {
  config_id    = configcat_config.my_config.id
  key          = "isAwesomeFeatureEnabled"
  name         = "My awesome feature flag"
  hint         = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
  order        = 0
}

resource "configcat_setting_value" "is_awesome_value" {
  environment_id = configcat_environment.my_environment.id
  setting_id     = configcat_setting.is_awesome.id

  value = "false"

  rollout_rules {
    segment_comparator = "isIn"
    segment_id         = configcat_segment.dogfooding.id
    value              = "true"
  }
  rollout_rules {
    segment_comparator = "isNotIn"
    segment_id         = configcat_segment.sensitive_users.id
    value              = "true"
  }
}
```
