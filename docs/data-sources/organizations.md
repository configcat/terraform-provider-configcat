# configcat_organizations Resource

Use this data source to access information about existing Organizations.

## Example Usage

```hcl
data "configcat_organizations" "organizations" {
  name_filter_regex = "ConfigCat"
}


output "organization_id" {
  value = data.configcat_organizations.organizations.organizations.0.organization_id
}
```

## Argument Reference

* `name_filter_regex` - (Optional) Filter the Organizations by name.

## Attribute Reference

* `organizations` - An organization [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `organizations` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `organization_id` - The unique Organization ID.
* `name` - The name of the Organization.

## Used APIs
- [Get Organizations](https://api.configcat.com/docs/index.html#operation/get-organizations)