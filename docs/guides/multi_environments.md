---
page_title: "Advanced usage of Resources in multiple environments"
---

# Advanced usage of Resources in multiple environments

## Prerequisites

[Get your Public Management API credentials](https://app.configcat.com/my-account/public-api-credentials) and set the following environment variables:
- CONFIGCAT_BASIC_AUTH_USERNAME
- CONFIGCAT_BASIC_AUTH_PASSWORD

## Folder/file structure

    .
    ├── root.tf
    ├── test
        ├── root.tf
    ├── production
        ├── root.tf

### root.tf

```hcl
terraform {
  required_providers {
    configcat = {
      source = "configcat/configcat"
      version = "~> 3.0"
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
  order = 0
}

resource "configcat_config" "my_config" {
  product_id = configcat_product.my_product.id
  name = "My config"
  order = 0
}

resource "configcat_setting" "is_awesome" {
  config_id = configcat_config.my_config.id
  key = "isAwesomeFeatureEnabled"
  name = "My awesome feature flag"
  hint = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
  order = 0
}

resource "configcat_setting" "welcome_text" {
  config_id = configcat_config.my_config.id
  key = "welcomeText"
  name = "Welcome text"
  hint = "Welcome text message shown on homepage"
  setting_type = "string"
  order = 1
}

resource "configcat_tag" "created_by_terraform_tag" {
  product_id = configcat_product.my_product.id
  name = "Created by Terraform"
}

resource "configcat_setting_tag" "is_awesome_setting_tag" {
    setting_id = configcat_setting.is_awesome.id
    tag_id = configcat_tag.created_by_terraform_tag.id
}

resource "configcat_setting_tag" "welcome_text_setting_tag" {
    setting_id = configcat_setting.welcome_text.id
    tag_id = configcat_tag.created_by_terraform_tag.id
}

// Test module
module "test" {
  source = "./test"

  product_id = configcat_product.my_product.id
  is_awesome_setting_id  = configcat_setting.is_awesome.id
  welcome_text_setting_id  = configcat_setting.welcome_text.id
}

// Production module
module "production" {
  source = "./production"

  product_id = configcat_product.my_product.id
  is_awesome_setting_id  = configcat_setting.is_awesome.id
  welcome_text_setting_id  = configcat_setting.welcome_text.id
}
```

### test/root.tf

```hcl
terraform {
  required_providers {
    configcat = {
      source = "configcat/configcat"
      version = "~> 3.0"
    }
  }
}

variable "product_id" { default = "" }
variable "setting_is_awesome_id" { default = "" }
variable "setting_welcome_text_id" { default = "" }

resource "configcat_environment" "test_environment" {
  product_id = var.product_id
  name = "Test"
  order = 0
}

resource "configcat_setting_value" "is_awesome_value" {
    environment_id = configcat_environment.test_environment.id
    setting_id = var.is_awesome_setting_id
    value = "true"
}

resource "configcat_setting_value" "welcome_text_value" {
    environment_id = configcat_environment.test_environment.id
    setting_id = var.welcome_text_setting_id
    value = "Welcome to ConfigCat"
}
```


### production/root.tf

```hcl
terraform {
  required_providers {
    configcat = {
      source = "configcat/configcat"
      version = "~> 3.0"
    }
  }
}

variable "product_id" { default = "" }
variable "setting_is_awesome_id" { default = "" }
variable "setting_welcome_text_id" { default = "" }

resource "configcat_environment" "production_environment" {
  product_id = var.product_id
  name = "Production"
  order = 1
}

resource "configcat_setting_value" "is_awesome_value" {
    environment_id = configcat_environment.production_environment.id
    setting_id = var.is_awesome_setting_id
    
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

resource "configcat_setting_value" "welcome_text_value" {
    environment_id = configcat_environment.production_environment.id
    setting_id = var.welcome_text_setting_id
    
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