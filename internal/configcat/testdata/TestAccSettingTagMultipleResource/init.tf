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

resource "configcat_setting_tag" "settingTag1" {
  setting_id = configcat_setting.testSetting.id
  tag_id     = configcat_tag.testTag1.id
}

resource "configcat_setting_tag" "settingTag2" {
  setting_id = configcat_setting.testSetting.id
  tag_id     = configcat_tag.testTag2.id
}

resource "configcat_setting_tag" "settingTag3" {
  setting_id = configcat_setting.testSetting.id
  tag_id     = configcat_tag.testTag3.id
}

