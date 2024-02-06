terraform {
  required_providers {
    configcat = {
      source  = "configcat/configcat"
      version = "~> 4.0"
    }
  }
}

variable "product_id" {
  type    = string
  default = ""
}
variable "is_awesome_setting_id" {
  type    = string
  default = ""
}
variable "welcome_text_setting_id" {
  type    = string
  default = ""
}

resource "configcat_environment" "production_environment" {
  product_id = var.product_id
  name       = "Production"
  order      = 1
}

resource "configcat_setting_value" "is_awesome_value" {
  environment_id = configcat_environment.production_environment.id
  setting_id     = var.is_awesome_setting_id

  value = "false"

  rollout_rules {
    comparison_attribute = "email"
    comparator           = "contains"
    comparison_value     = "@mycompany.com"
    value                = "true"
  }

  percentage_items {
    percentage = 20
    value      = "true"
  }
  percentage_items {
    percentage = 80
    value      = "false"
  }
}

resource "configcat_setting_value" "welcome_text_value" {
  environment_id = configcat_environment.production_environment.id
  setting_id     = var.welcome_text_setting_id

  value = "Welcome to ConfigCat"

  rollout_rules {
    comparison_attribute = "email"
    comparator           = "contains"
    comparison_value     = "@configcat.com"
    value                = "Hey buddies!"
  }

  percentage_items {
    percentage = 20
    value      = "Welcome to ConfigCat. 10 minutes trainable feature flag and configuration management service"
  }
  percentage_items {
    percentage = 80
    value      = "Welcome to ConfigCat. Unlimited team size, awesome support and no surprises."
  }
}
