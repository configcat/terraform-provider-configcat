terraform {
  required_providers {
    configcat = {
      source  = "configcat/configcat"
      version = "~> 5.0"
    }
  }
}

provider "configcat" {
}

variable "product_id" {
  type = string
}

resource "configcat_config" "my_config" {
  product_id         = var.product_id
  name               = "My config"
  description        = "My config description"
  order              = 0
  evaluation_version = "v2"
}

resource "configcat_environment" "my_environment" {
  product_id  = var.product_id
  name        = "Production"
  description = "Production description"
  color       = "blue"
  order       = 0
}

resource "configcat_setting" "my_string_setting" {
  config_id    = configcat_config.my_config.id
  key          = "myStringSetting"
  name         = "My awesome text setting"
  hint         = "This is the hint for my text setting"
  setting_type = "string"
  order        = 0

  predefined_variations = [
    {
      value = {
        string_value = "Variation A"
      }
    },
    {
      value = {
        string_value = "Variation B"
      }
      name = "Display name of Variation B"
    },
    {
      value = {
        string_value = "Variation C"
      }
      hint = "Hint for Variation C"
    }
  ]
}


resource "configcat_setting" "my_bool_setting" {
  config_id    = configcat_config.my_config.id
  key          = "myBoolSetting"
  name         = "My awesome bool feature flag"
  hint         = "This is the hint for my awesome bool feature flag"
  setting_type = "string"
  order        = 0

  predefined_variations = [
    {
      value = {
        bool_value = false
      }
      name = "Display name for Off value"
    },
    {
      value = {
        bool_value = true
      }
      hint = "Hint for the On value"
    },
  ]
}

resource "configcat_setting_value_v2" "my_bool_setting_value" {
  environment_id = configcat_environment.my_environment.id
  setting_id     = configcat_setting.my_bool_setting.id

  value = { predefined_variation_id = configcat_setting.my_bool_setting.predefined_variations[0].predefined_variation_id }
}

resource "configcat_setting_value_v2" "my_string_setting_value" {
  environment_id = configcat_environment.my_environment.id
  setting_id     = configcat_setting.my_string_setting.id

  value = { predefined_variation_id = configcat_setting.my_string_setting.predefined_variations[0].predefined_variation_id }

  targeting_rules = [
    {
      conditions = [
        {
          prerequisite_flag_condition = {
            prerequisite_setting_id = configcat_setting.my_bool_setting.id
            comparator              = "doesNotEqual"
            comparison_value        = { predefined_variation_id = configcat_setting.my_bool_setting.predefined_variations[0].predefined_variation_id }
          }
        }
      ]
      value = { predefined_variation_id = configcat_setting.my_string_setting.predefined_variations[0].predefined_variation_id }
    },

    {
      percentage_options = [
        {
          percentage = 30
          value      = { predefined_variation_id = configcat_setting.my_string_setting.predefined_variations[1].predefined_variation_id }
        },
        {
          percentage = 70
          value      = { predefined_variation_id = configcat_setting.my_string_setting.predefined_variations[2].predefined_variation_id }
        }
      ]
    },
  ]
}

