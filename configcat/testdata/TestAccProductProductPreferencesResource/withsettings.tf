variable "key_generation_mode" {
  type    = string
  default = null
}

variable "mandatory_setting_hint" {
  type    = bool
  default = null
}

variable "show_variation_id" {
  type    = bool
  default = null
}

variable "reason_required" {
  type    = bool
  default = null
}

resource "configcat_product" "product" {
  organization_id = "08d86d63-26dc-4276-86d6-eae122660e51"
  name            = "Product preferences test"
  order           = 1
}

resource "configcat_product_preferences" "preferences" {
  product_id = configcat_product.product.id

  key_generation_mode    = var.key_generation_mode
  mandatory_setting_hint = var.mandatory_setting_hint
  show_variation_id      = var.show_variation_id
  reason_required        = var.reason_required
}
