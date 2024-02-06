# Retrieve your Product
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

# Retrieve your Config
data "configcat_configs" "my_configs" {
  product_id        = data.configcat_products.my_products.products[0].product_id
  name_filter_regex = "Main Config"
}

# Retrieve your Environment
data "configcat_environments" "my_environments" {
  product_id        = data.configcat_products.my_products.products[0].product_id
  name_filter_regex = "Test"
}

# Create a Feature Flag/Setting
resource "configcat_setting" "setting" {
  config_id    = data.configcat_configs.my_configs.configs[0].config_id
  key          = "isAwesomeFeatureEnabled"
  name         = "My awesome feature flag"
  hint         = "This is the hint for my awesome feature flag"
  setting_type = "boolean"
  order        = 0
}

# Set a value to the Feature Flag/Setting created above
resource "configcat_setting_value" "setting_value" {
  environment_id = data.configcat_environments.my_environments.environments[0].environment_id
  setting_id     = configcat_setting.setting.id
  value          = "false"
}
