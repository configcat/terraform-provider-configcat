# configcat_tag Resource

Creates and manages a **Tag**.  

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

resource "configcat_tag" "my_tag" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "Created by Terraform"
}


output "tag_id" {
  value = configcat_tag.my_tag.id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name` - (Required) The name of the Tag.
* `color` - (Optional) Default: `panther`. The color of the Tag. Valid values: `panther`, `whale`, `salmon`, `lizard`, `canary`, `koala`.

## Attribute Reference

* `id` - The unique Tag ID.

## Import

Tags can be imported using the TagId. Get the TagId using e.g. the [List Tags API](https://api.configcat.com/docs/#tag/Tags/operation/get-tags).

```
$ terraform import configcat_tag.example 1234
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Tag](https://api.configcat.com/docs/#tag/Tags/operation/get-tag)
* [Create Tag](https://api.configcat.com/docs/#tag/Tags/operation/create-tag)
* [Update Tag](https://api.configcat.com/docs/#tag/Tags/operation/update-tag)
* [Delete Tag](https://api.configcat.com/docs/#tag/Tags/operation/delete-tag)
