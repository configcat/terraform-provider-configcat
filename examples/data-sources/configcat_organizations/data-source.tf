data "configcat_organizations" "my_organizations" {
  name_filter_regex = "ConfigCat"
}


output "organization_id" {
  value = data.configcat_organizations.my_organizations.organizations[0].organization_id
}