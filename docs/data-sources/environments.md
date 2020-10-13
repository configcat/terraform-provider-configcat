# configcat_environments Resource

Use this data source to access information about an existing Environment. [Used API](https://api.configcat.com/docs/index.html#operation/get-environments)

## Example Usage

```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_environments" "environments" {
  product_id = data.configcat_product.products.products.0.product_id
  name_filter_regex = "Test"
}


output "environment_id" {
  value = data.configcat_environments.environments.environments.0.environment_id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product
* `name_filter_regex` - (Optional) Filter the Environments by name.

## Attribute Reference

* `environments` - An environment [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `environments` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `environment_id` - The unique Environment ID.
* `name` - The name of the Environment.