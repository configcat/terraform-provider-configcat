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

data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_permission_groups" "my_permission_groups" {
  product_id        = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Administrators"
}

data "configcat_configs" "my_configs" {
  product_id        = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Main Config"
}

data "configcat_environments" "my_environments" {
  product_id        = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Test"
}

data "configcat_sdkkeys" "my_sdkkeys" {
  config_id      = data.configcat_configs.my_configs.configs.0.config_id
  environment_id = data.configcat_environments.my_environments.environments.0.environment_id
}

data "configcat_settings" "my_settings" {
  config_id        = data.configcat_configs.my_configs.configs.0.config_id
  key_filter_regex = "isAwesomeFeatureEnabled"
}

data "configcat_tags" "my_tags" {
  product_id        = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Test"
}
