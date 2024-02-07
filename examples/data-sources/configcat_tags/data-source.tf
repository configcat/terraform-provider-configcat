variable "product_id" {
  type = string
}

data "configcat_tags" "my_tags" {
  product_id        = var.product_id
  name_filter_regex = "Test"
}


output "tag_id" {
  value = data.configcat_tags.my_tags.tags[0].tag_id
}
