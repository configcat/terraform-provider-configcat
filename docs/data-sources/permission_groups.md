# configcat_permission_groups Resource

Use this data source to access information about existing **Permission Groups**. [What is a Permission Group in ConfigCat?](https://configcat.com/docs/advanced/team-management/team-management-basics/#permissions--permission-groups-product-level)

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_permission_groups" "my_permission_groups" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Administrators"
}


output "permission_group_id" {
  value = data.configcat_permission_groups.my_permission_groups.permission_groups.0.permission_group_id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name_filter_regex` - (Optional) Filter the Permission Groups by name.

## Attribute Reference

* `permission_groups` - A permission group [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `permission_groups` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `permission_group_id` - The unique Permission Groups ID.
* `name` - The name of the Permission Group.
* `can_manage_members` - Group members can manage team members.
* `can_createorupdate_config` - Group members can create/update Configs.
* `can_delete_config` - Group members can delete Configs.
* `can_createorupdate_environment` - Group members can create/update Environments.
* `can_delete_environment` - Group members can delete Environments.
* `can_createorupdate_setting` - Group members can create/update Feature Flags and Settings.
* `can_tag_setting` - Group members can attach/detach Tags to Feature Flags and Settings.
* `can_delete_setting` - Group members can delete Feature Flags and Settings.
* `can_createorupdate_tag` - Group members can create/update Tags.
* `can_delete_tag` - Group members can delete Tags.
* `can_manage_webhook` - Group members can create/update/delete Webhooks.
* `can_use_exportimport` - Group members can use the export/import feature.
* `can_manage_product_preferences` - Group members can update Product preferences.
* `can_manage_integrations` - Group members can add and configure integrations.
* `can_view_sdkkey` - Group members has access to SDK keys.
* `can_rotate_sdkkey` - Group members can rotate SDK keys.
* `can_createorupdate_segments` - Group members can create/update Segments.
* `can_delete_segments` - Group members can delete Segments.
* `can_view_product_auditlog` - Group members has access to audit logs.
* `can_view_product_statistics` - Group members has access to product statistics.
* `accesstype` - Represent the Feature Management permission. Possible values: readOnly, full, custom
* `new_environment_accesstype` - Represent the environment specific Feature Management permission for new Environments. Possible values: full, readOnly, none
* `environment_accesses` - The environment specific permissions [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `environment_accesses` [list](https://www.terraform.io/docs/configuration/types.html#list-) block
* `environment_id` - The unique [Environment](https://configcat.com/docs/main-concepts/#environment) ID.
* `environment_access_type` - Represent the environment specific Feature Management permission. Possible values: full, readOnly, none

## Endpoints used
- [Get Permission Groups](https://api.configcat.com/docs/index.html#tag/Permission-Groups/operation/get-permission-groups)
