# Simple usage of Data sources

## Prerequisites

Get your ConfigCat Public API credentials at https://app.configcat.com/my-account/public-api-credentials and set the following environment variables:
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
data "configcat_organizations" "organizations" {
  name_filter_regex = "ConfigCat"
}

resource "configcat_product" "product" {
  organization_id = data.configcat_organizations.organizations.organizations.0.organization_id
  name = "My product"
}

resource "configcat_config" "config" {
  product_id = configcat_product.product.id
  name = "My config"
}

resource "configcat_setting" "is_awesome" {
  config_id = configcat_config.config.id
  key = "isAwesomeFeatureEnabled"
  name = "My awesome feature flag"
  hint = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
}

resource "configcat_setting" "welcome_text" {
  config_id = configcat_config.config.id
  key = "welcomeText"
  name = "Welcome text"
  hint = "Welcome text message shown on homepage"
  setting_type = "text"
}

resource "configcat_tag" "tag" {
  product_id = configcat_product.product.id
  name = "Created by Terraform"
}

resource "configcat_setting_tag" "setting_tag" {
    setting_id = configcat_setting.is_awesome.id
    tag_id = configcat_tag.tag.id
}

resource "configcat_setting_tag" "setting_tag" {
    setting_id = configcat_setting.welcome_text.id
    tag_id = configcat_tag.tag.id
}

// Test module
module "test" {
  source = "./test"

  product_id = configcat_product.product.id
  setting_is_awesome_id  = configcat_setting.is_awesome.id
  setting_welcome_text_id  = configcat_setting.welcome_text.id
}

// Production module
module "production" {
  source = "./production"

  product_id = configcat_product.product.id
  setting_is_awesome_id  = configcat_setting.is_awesome.id
  setting_welcome_text_id  = configcat_setting.welcome_text.id
}
```

## test.tf

```hcl
variable "product_id" { default = "" }
variable "setting_is_awesome_id" { default = "" }
variable "setting_welcome_text_id" { default = "" }

resource "configcat_environment" "environment" {
  product_id = var.product_id
  name = "Test"
}

resource "configcat_setting_value" "setting_value_is_awesome" {
    environment_id = configcat_environment.environment.id
    setting_id = var.setting_is_awesome_id
    value = "true"
}

resource "configcat_setting_value" "setting_value_setting_welcome_text_id" {
    environment_id = configcat_environment.environment.id
    setting_id = var.setting_welcome_text_id
    value = "Welcome to ConfigCat"
}
```


## production.tf

```hcl
variable "product_id" { default = "" }
variable "setting_is_awesome_id" { default = "" }
variable "setting_welcome_text_id" { default = "" }

resource "configcat_environment" "environment" {
  product_id = var.product_id
  name = "Production"
}

resource "configcat_setting_value" "setting_value_is_awesome" {
    environment_id = configcat_environment.environment.id
    setting_id = var.setting_is_awesome_id
    
    value = "false"

    rollout_rules {
        comparison_attribute = "email"
        comparator = "contains"
        comparison_value = "@mycompany.com"
        value = "true"
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

resource "configcat_setting_value" "setting_value_setting_welcome_text_id" {
    environment_id = configcat_environment.environment.id
    setting_id = var.setting_welcome_text_id
    
    value = "Welcome to ConfigCat"

    rollout_rules {
        comparison_attribute = "email"
        comparator = "contains"
        comparison_value = "@configcat.com"
        value = "Hey buddies!"
    }

    percentage_items {
        percentage = 20
        value = "Welcome to ConfigCat. 10 minutes trainable feature flag and configuration management service"
    }
    percentage_items {
        percentage = 80
        value = "Welcome to ConfigCat. Unlimited team size, awesome support and no surprises."
    }
}
```