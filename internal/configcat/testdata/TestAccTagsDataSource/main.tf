variable "product_id" {
  type = string
}

variable "name_filter_regex" {
  type    = string
  default = null
}

data "configcat_tags" "test" {
  product_id        = var.product_id
  name_filter_regex = var.name_filter_regex
}

output "tag_id" {
  value = data.configcat_tags.test.tags[0].tag_id
}
