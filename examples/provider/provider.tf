terraform {
  required_providers {
    configcat = {
      source  = "configcat/configcat"
      version = "~> 3.0"
    }
  }
}

variable "configcat_basic_auth_username" {
  type      = string
  sensitive = true
}

variable "configcat_basic_auth_password" {
  type      = string
  sensitive = true
}

provider "configcat" {
  // Get your ConfigCat Public API credentials at https://app.configcat.com/my-account/public-api-credentials
  basic_auth_username = var.configcat_basic_auth_username
  basic_auth_password = var.configcat_basic_auth_password
}
