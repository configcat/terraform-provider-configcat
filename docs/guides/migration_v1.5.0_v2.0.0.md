---
page_title: "Migration guide from v1.5.0 to v2.0.0"
---

# Migration guide from v1.5.0 to v2.0.0

## Breaking change in v2.0.0

Permission Group handling was introduced in v1.5.0 and it had a problem with handling custom Environment accesses. We had to refactor it and introduce a breaking change in v2.0.0.

You could define the custom Environment accesses in v1.5.0 with a list property:

```hcl
resource "configcat_permission_group" "my_permission_group" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "Read only except Test environment"

  accesstype = "custom"

  environment_access {
    environment_id = data.configcat_environments.my_test_environments.environments.0.environment_id
    environment_accesstype = "full"
  }

  environment_access {
    environment_id = data.configcat_environments.my_production_environments.environments.0.environment_id
    environment_accesstype = "none"
  }
}
```

The new way of defining custom Environment accesses is using a map property:

```hcl
resource "configcat_permission_group" "my_permission_group" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "Read only except Test environment"

  accesstype = "custom"

  environment_accesses {
    "${data.configcat_environments.my_test_environments.environments.0.environment_id}" = "full"
    "${data.configcat_environments.my_test_environments.environments.1.environment_id}" = "readOnly"
  }
}
```
