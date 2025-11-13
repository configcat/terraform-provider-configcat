
variable "product_id" {
  type = string
}

resource "configcat_config" "config" {
  product_id         = var.product_id
  name               = "Test config for predefined variations values"
  evaluation_version = "v2"
  order              = 0
}

resource "configcat_setting" "bool_setting" {
  config_id    = configcat_config.config.id
  key          = "boolKey"
  name         = "bool setting name"
  setting_type = "boolean"
  predefined_variations = [
    {
      value = {
        bool_value = false
      }
      name = "Off name"
    },
    {
      value = {
        bool_value = true
      }
      hint = "On hint"
    }
  ]
  order = 0
}

resource "configcat_setting" "string_setting" {
  config_id    = configcat_config.config.id
  key          = "stringKey"
  name         = "string setting name"
  setting_type = "string"
  predefined_variations = [
    {
      value = {
        string_value = "Variation A"
      }
      name = "Variation A name"
    },
    {
      value = {
        string_value = "Variation B"
      }
      hint = "Variation B hint"
    }
  ]
  order = 1
}

resource "configcat_setting" "int_setting" {
  config_id    = configcat_config.config.id
  key          = "intKey"
  name         = "int setting name"
  setting_type = "int"
  predefined_variations = [
    {
      value = {
        int_value = 10
      }
      name = "10 name"
    },
    {
      value = {
        int_value = 20
      }
      hint = "20 hint"
    }
  ]
  order = 1
}

resource "configcat_setting" "double_setting" {
  config_id    = configcat_config.config.id
  key          = "doubleKey"
  name         = "double setting name"
  setting_type = "double"
  predefined_variations = [
    {
      value = {
        double_value = 10.0
      }
      name = "10.0 name"
    },
    {
      value = {
        double_value = 20.2
      }
      hint = "20.2 hint"
    }
  ]
  order = 1
}
