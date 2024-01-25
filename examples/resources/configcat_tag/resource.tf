variable "product_id" {
  type = string
}

resource "configcat_tag" "my_tag" {
  product_id = var.product_id
  name       = "Created by Terraform"
  color      = "panther"
}


output "tag_id" {
  value = configcat_tag.my_tag.id
}
