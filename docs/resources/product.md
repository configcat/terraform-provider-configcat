# configcat_product Resource

Creates and manages a **Product**. [What is a Product in ConfigCat?](https://configcat.com/docs/main-concepts)

## Example Usage

```hcl
data "configcat_organizations" "my_organizations" {
  name_filter_regex = "ConfigCat"
}

resource "configcat_product" "my_product" {
  organization_id = data.configcat_organizations.my_organizations.organizations.0.organization_id
  name = "My product"
}


output "product_id" {
  value = configcat_product.my_product.id
}
```

## Argument Reference

* `organization_id` - (Required) The ID of the Organization.
* `name` - (Required) The name of the Product.

## Attribute Reference

* `id` - The unique Product ID.

## Import

Products can be imported using the ProductId. Get the ProductId using the [GetProducts API](https://api.configcat.com/docs/#operation/get-products) for example.

```
$ terraform import configcat_product.example 08d86d63-2726-47cd-8bfc-59608ecb91e2
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Product](https://api.configcat.com/docs/index.html#operation/get-product)
* [Create Product](https://api.configcat.com/docs/index.html#operation/create-product)
* [Update Product](https://api.configcat.com/docs/index.html#operation/update-product)
* [Delete Product](https://api.configcat.com/docs/index.html#operation/delete-product)
