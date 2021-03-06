# configcat_environment Resource

Creates and manages an **Environment**. [What is an Environment in ConfigCat?](https://configcat.com/docs/main-concepts)

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

resource "configcat_environment" "my_environment" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "Staging"
}


output "environment_id" {
  value = configcat_environment.my_environment.id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name` - (Required) The name of the Environment.

## Attribute Reference

* `id` - The unique Environment ID.

## Import

Environments can be imported using the EnvironmentId. Get the EnvironmentId using the [GetEnvironments API](https://api.configcat.com/docs/#operation/get-environments) for example.

```
$ terraform import configcat_environment.example 08d86d63-2726-47cd-8bfc-59608ecb91e2
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Environment](https://api.configcat.com/docs/index.html#operation/get-environment)
* [Create Environment](https://api.configcat.com/docs/index.html#operation/create-environment)
* [Update Environment](https://api.configcat.com/docs/index.html#operation/update-environment)
* [Delete Environment](https://api.configcat.com/docs/index.html#operation/delete-environment)
