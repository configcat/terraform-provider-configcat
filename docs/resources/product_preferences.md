---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "configcat_product_preferences Resource - terraform-provider-configcat"
subcategory: ""
description: |-
  Manages the Product Preferences.
---

# configcat_product_preferences (Resource)

Manages the **Product Preferences**.

## Example Usage

```terraform
variable "organization_id" {
  type = string
}

resource "configcat_product" "product" {
  organization_id = var.organization_id
  name            = "My product"
  order           = 0
}

resource "configcat_environment" "test" {
  product_id = configcat_product.product.id
  name       = "Test"
  order      = 0
}

resource "configcat_environment" "production" {
  product_id = configcat_product.product.id
  name       = "Production"
  order      = 1
}

resource "configcat_product_preferences" "preferences" {
  product_id = configcat_product.product.id

  key_generation_mode    = "kebabCase"
  mandatory_setting_hint = true
  show_variation_id      = false
  reason_required        = false
  reason_required_environments = {
    "${configcat_environment.test.id}"       = false,
    "${configcat_environment.production.id}" = true,
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `product_id` (String) The ID of the Product.

### Optional

- `key_generation_mode` (String) Determines the Feature Flag key generation mode. Available values: `camelCase`|`upperCase`|`lowerCase`|`pascalCase`|`kebabCase`. Default: `camelCase`.
- `mandatory_setting_hint` (Boolean) Indicates whether Feature flags and Settings must have a hint. Default: false.
- `reason_required` (Boolean) Indicates that a mandatory note is required for saving and publishing. Default: false.
- `reason_required_environments` (Map of Boolean) The environment specific mandatory note map block. Keys are the Environment IDs and the values indicate that a mandatory note is required for saving and publishing.
- `show_variation_id` (Boolean) Indicates whether variation IDs must be shown on the ConfigCat Dashboard. Default: false.

### Read-Only

- `id` (String) Internal ID of the resource. Do not use.

## Import

Import is supported using the following syntax:

```shell
# Product preferences can be imported using the ProductId. Get the ProductId using the [List Products API](https://api.configcat.com/docs/#tag/Products/operation/get-products) for example.

terraform import configcat_product_preferences.example 08d86d63-2726-47cd-8bfc-59608ecb91e2
```