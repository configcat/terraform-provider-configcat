variable "name_filter_regex" {
  type    = string
  default = null
}

data "configcat_organizations" "test" {
  name_filter_regex = var.name_filter_regex
}

output "organization_id" {
  value = length(data.configcat_organizations.test.organizations) > 0 ? data.configcat_organizations.test.organizations[0].organization_id : null
}
