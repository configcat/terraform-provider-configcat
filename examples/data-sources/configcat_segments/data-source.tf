variable "product_id" {
  type = string
}

data "configcat_segments" "my_segments" {
  product_id        = var.product_id
  name_filter_regex = "Beta users"
}


output "segment_id" {
  value = data.configcat_segments.my_segments.segments[0].segment_id
}
