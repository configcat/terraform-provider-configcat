variable "product_id" {
  type = string
}

variable "config_id" {
  type = string
}

resource "configcat_setting" "testSetting" {
  config_id = var.config_id
  key       = "testkeywithtag"
  name      = "test"
  order     = 0
}

resource "configcat_tag" "testTag1" {
  product_id = var.product_id
  name       = "tag1"
}

resource "configcat_tag" "testTag2" {
  product_id = var.product_id
  name       = "tag2"
}

resource "configcat_tag" "testTag3" {
  product_id = var.product_id
  name       = "tag3"
}
