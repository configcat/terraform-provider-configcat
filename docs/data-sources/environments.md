# configcat_environments Resource

Use this data source to access information about existing **Environments**. [What is an Environment in ConfigCat?](https://configcat.com/docs/main-concepts)

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_environments" "my_environments" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Test"
}


output "environment_id" {
  value = data.configcat_environments.my_environments.environments.0.environment_id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name_filter_regex` - (Optional) Filter the Environments by name.

## Attribute Reference

* `environments` - An environment [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `environments` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `environment_id` - The unique Environment ID.
* `name` - The name of the Environment.
* `description` - The description of the Environment.
* `color` - The color of the Environment.

## Endpoints used
- [Get Environments](https://api.configcat.com/docs/index.html#operation/get-environments)
