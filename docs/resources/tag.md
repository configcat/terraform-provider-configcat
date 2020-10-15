# configcat_tag Resource

Creates and manages a **Tag**.  

## Example Usage

```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

resource "configcat_tag" "tag" {
  product_id = data.configcat_products.products.products.0.product_id
  name = "Created by Terraform"
}


output "tag_id" {
  value = configcat_tag.tag.id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name` - (Required) The name of the Tag.

## Attribute Reference

* `id` - The unique Tag ID.

## Import

Tags can be imported using the TagId. Get the TagId using e.g. the [GetTags API](https://api.configcat.com/docs/#operation/get-tags).

```
$ terraform import configcat_tag.example 1234
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Used APIs
* [Get Tag](https://api.configcat.com/docs/index.html#operation/get-tag)
* [Create Tag](https://api.configcat.com/docs/index.html#operation/create-tag)
* [Update Tag](https://api.configcat.com/docs/index.html#operation/update-tag)
* [Delete Tag](https://api.configcat.com/docs/index.html#operation/delete-tag)
