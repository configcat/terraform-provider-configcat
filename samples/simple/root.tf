terraform {
  required_providers {
    configcat = {
      source = "configcat/configcat"
      version = "0.1.2-alpha"
    }
  }
}

provider "configcat" {
}

data "configcat_products" "products" {
  // Filter for your Product
  name_filter_regex = "Configcat's product"
}

data "configcat_configs" "configs" {
  product_id = data.configcat_products.products.products.0.product_id
  // Filter for your Config
  name_filter_regex = "Main Config"
}

data "configcat_environments" "environments" {
  product_id = data.configcat_products.products.products.0.product_id
  // Filter for your Environment
  name_filter_regex = "Test"
}

resource "configcat_setting" "setting" {
  config_id = data.configcat_configs.configs.configs.0.config_id
  key = "facebookSharingEnabled"
  name = "Facebook sharing enabled"
  hint = "Created by Terraform"
  setting_type = "boolean"
}

resource "configcat_setting_value" "setting_value" {
    environment_id = data.configcat_environments.environments.environments.0.environment_id
    setting_id = configcat_setting.setting.id
    init_only = false
    value = "false"

    rollout_rules {
        comparison_attribute = "Email"
        comparator = "contains"
        comparison_value = "@configcat.com"
        value = "true"
    }
}