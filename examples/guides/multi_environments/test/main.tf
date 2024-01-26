terraform {
  required_providers {
    configcat = {
      source  = "configcat/configcat"
      version = "~> 3.0"
    }
  }
}

variable "product_id" { default = "" }
variable "is_awesome_setting_id" { default = "" }
variable "welcome_text_setting_id" { default = "" }

resource "configcat_environment" "test_environment" {
  product_id = var.product_id
  name       = "Test"
  order      = 0
}

resource "configcat_setting_value" "is_awesome_value" {
  environment_id = configcat_environment.test_environment.id
  setting_id     = var.is_awesome_setting_id
  value          = "true"
}

resource "configcat_setting_value" "welcome_text_value" {
  environment_id = configcat_environment.test_environment.id
  setting_id     = var.welcome_text_setting_id
  value          = "Welcome to ConfigCat"
}
