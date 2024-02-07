variable "product_id" {
  type = string
}

variable "name_filter_regex" {
  type    = string
  default = null
}

data "configcat_segments" "test" {
  product_id        = var.product_id
  name_filter_regex = var.name_filter_regex
}

output "segment_id" {
  value = length(data.configcat_segments.test.segments) > 0 ? data.configcat_segments.test.segments[0].segment_id : null
}
