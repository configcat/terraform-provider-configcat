# configcat_configs Resource

Use this data source to access information about existing **Configs**. [What is a Config in ConfigCat?](https://configcat.com/docs/main-concepts)

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_configs" "my_configs" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Main Config"
}


output "config_id" {
  value = data.configcat_configs.my_configs.configs.0.config_id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name_filter_regex` - (Optional) Filter the Configs by name.

## Attribute Reference

* `configs` - A config [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `configs` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `config_id` - The unique Config ID.
* `name` - The name of the Config.

## Endpoints used
[Get Configs](https://api.configcat.com/docs/index.html#operation/get-configs)