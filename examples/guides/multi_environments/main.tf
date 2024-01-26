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

//  Organization Resource is ReadOnly.
data "configcat_organizations" "my_organizations" {
  name_filter_regex = "ConfigCat"
}

resource "configcat_product" "my_product" {
  organization_id = data.configcat_organizations.my_organizations.organizations.0.organization_id
  name            = "My product"
  order           = 0
}

resource "configcat_config" "my_config" {
  product_id = configcat_product.my_product.id
  name       = "My config"
  order      = 0
}

resource "configcat_setting" "is_awesome" {
  config_id    = configcat_config.my_config.id
  key          = "isAwesomeFeatureEnabled"
  name         = "My awesome feature flag"
  hint         = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
  order        = 0
}

resource "configcat_setting" "welcome_text" {
  config_id    = configcat_config.my_config.id
  key          = "welcomeText"
  name         = "Welcome text"
  hint         = "Welcome text message shown on homepage"
  setting_type = "string"
  order        = 1
}

resource "configcat_tag" "created_by_terraform_tag" {
  product_id = configcat_product.my_product.id
  name       = "Created by Terraform"
}

resource "configcat_setting_tag" "is_awesome_setting_tag" {
  setting_id = configcat_setting.is_awesome.id
  tag_id     = configcat_tag.created_by_terraform_tag.id
}

resource "configcat_setting_tag" "welcome_text_setting_tag" {
  setting_id = configcat_setting.welcome_text.id
  tag_id     = configcat_tag.created_by_terraform_tag.id
}

// Test module
module "test" {
  source = "./test"

  product_id              = configcat_product.my_product.id
  is_awesome_setting_id   = configcat_setting.is_awesome.id
  welcome_text_setting_id = configcat_setting.welcome_text.id
}

// Production module
module "production" {
  source = "./production"

  product_id              = configcat_product.my_product.id
  is_awesome_setting_id   = configcat_setting.is_awesome.id
  welcome_text_setting_id = configcat_setting.welcome_text.id
}
