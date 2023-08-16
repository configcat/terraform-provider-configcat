# configcat_permission_group Resource

Creates and manages a **Permission Group**. [What is a Permission Group in ConfigCat?](https://configcat.com/docs/advanced/team-management/team-management-basics/#permissions--permission-groups-product-level)

## Example Usages

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

resource "configcat_permission_group" "my_permission_group" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "Administrators"

  accesstype = "full"

  can_manage_members = true
  can_createorupdate_config = true
  can_delete_config = true
  can_createorupdate_environment = true
  can_delete_environment = true
  can_createorupdate_setting = true
  can_tag_setting = true
  can_delete_setting = true
  can_createorupdate_tag = true
  can_delete_tag = true
  can_manage_webhook = true
  can_use_exportimport = true
  can_manage_product_preferences = true
  can_manage_integrations = true
  can_view_sdkkey = true
  can_rotate_sdkkey = true
  can_createorupdate_segment = true
  can_delete_segment = true
  can_view_product_auditlog = true
  can_view_product_statistics = true
}

output "permission_group_id" {
  value = configcat_permission_group.my_permission_group.id
}
```

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_environments" "my_environments" {
  name_filter_regex = "Test"
}

resource "configcat_permission_group" "my_permission_group" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "Read only except Test environment"

  accesstype = "custom"

  environment_access {
    environment_id = data.configcat_environments.my_environments.environments.0.environment_id
    environment_access_type = "full"
  }
}

output "permission_group_id" {
  value = configcat_permission_group.my_permission_group.id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name` - (Required) The name of the Permission Group.
* `can_manage_members` - (Optional) Group members can manage team members. Default: false.
* `can_createorupdate_config` - (Optional) Group members can create/update Configs. Default: false.
* `can_delete_config` - (Optional) Group members can delete Configs. Default: false.
* `can_createorupdate_environment` - (Optional) Group members can create/update Environments. Default: false.
* `can_delete_environment` - (Optional) Group members can delete Environments. Default: false.
* `can_createorupdate_setting` - (Optional) Group members can create/update Feature Flags and Settings. Default: false.
* `can_tag_setting` - (Optional) Group members can attach/detach Tags to Feature Flags and Settings. Default: false.
* `can_delete_setting` - (Optional) Group members can delete Feature Flags and Settings. Default: false.
* `can_createorupdate_tag` - (Optional) Group members can create/update Tags. Default: false.
* `can_delete_tag` - (Optional) Group members can delete Tags. Default: false.
* `can_manage_webhook` - (Optional) Group members can create/update/delete Webhooks. Default: false.
* `can_use_exportimport` - (Optional) Group members can use the export/import feature. Default: false.
* `can_manage_product_preferences` - (Optional) Group members can update Product preferences. Default: false.
* `can_manage_integrations` - (Optional) Group members can add and configure integrations. Default: false.
* `can_view_sdkkey` - (Optional) Group members has access to SDK keys. Default: false.
* `can_rotate_sdkkey` - (Optional) Group members can rotate SDK keys. Default: false.
* `can_createorupdate_segments` - (Optional) Group members can create/update Segments. Default: false.
* `can_delete_segments` - (Optional) Group members can delete Segments. Default: false.
* `can_view_product_auditlog` - (Optional) Group members has access to audit logs. Default: false.
* `can_view_product_statistics` - (Optional) Group members has access to product statistics. Default: false.
* `accesstype` - Represent the Feature Management permission. Possible values: readOnly, full, custom. Default: custom
* `new_environment_accesstype` - Represent the environment specific Feature Management permission for new Environments and for those that are not specified in the environment_accesses list. Possible values: full, readOnly, none. Default: none.
* `environment_accesses` - The environment specific permissions [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `environment_accesses` [list](https://www.terraform.io/docs/configuration/types.html#list-) block
* `environment_id` - (Required) The unique [Environment](https://configcat.com/docs/main-concepts/#environment) ID.
* `environment_access_type` - (Optional) Represent the environment specific Feature Management permission. Possible values: full, readOnly, none. Default: none.

## Attribute Reference

* `id` - The unique Permission Group ID.

## Import

Permission Groups can be imported using the PermissionGroupId. Get the PermissionGroupId using the [List Permission Groups API](https://api.configcat.com/docs/#tag/Permission-Groups/operation/get-permission-groups) for example.

```
$ terraform import configcat_permission_group.example 123
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Permission Group](https://api.configcat.com/docs/#tag/Permission-Groups/operation/get-permission-group)
* [Create Permission Group](https://api.configcat.com/docs/#tag/Permission-Groups/operation/create-permission-group)
* [Update Permission Group](https://api.configcat.com/docs/#tag/Permission-Groups/operation/update-permission-group)
* [Delete Permission Group](https://api.configcat.com/docs/#tag/Permission-Groups/operation/delete-permission-group)
