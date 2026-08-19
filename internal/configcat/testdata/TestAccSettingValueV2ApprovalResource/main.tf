variable "value" {
  type = bool
}

variable "bypass_approval" {
  type    = bool
  default = null
}

resource "configcat_product" "test" {
  organization_id = "08d86d63-26dc-4276-86d6-eae122660e51"
  name            = "Approval test product"
  order           = 99
}

resource "configcat_product_preferences" "test" {
  product_id       = configcat_product.test.id
  approve_required = true
}

resource "configcat_config" "test" {
  product_id         = configcat_product.test.id
  name               = "Approval test config"
  evaluation_version = "v2"
  order              = 0
}

resource "configcat_environment" "test" {
  product_id = configcat_product.test.id
  name       = "Approval test environment"
  order      = 0
}

resource "configcat_setting" "test" {
  config_id = configcat_config.test.id
  key       = "ApprovalTestBoolKey"
  name      = "Approval test bool flag"
  order     = 0
}

resource "configcat_setting_value_v2" "test" {
  environment_id  = configcat_environment.test.id
  setting_id      = configcat_setting.test.id
  init_only       = false
  value           = { bool_value = var.value }
  bypass_approval = var.bypass_approval
}
