# configcat_config Resource

Creates and manages a Config.  

## Example Usage

```hcl
data "configcat_products" "products" {
  name_filter_regex = "ConfigCat's product"
}

resource "configcat_config" "config" {
  product_id = data.configcat_products.products.products.0.product_id
  name = "My config"
}


output "config_id" {
  value = configcat_config.config.id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name` - (Required) The name of the Config.

## Attribute Reference

* `id` - The unique Config ID.

## Import

Configs can be imported using the ConfigId. Get the ConfigId using e.g. the [GetConfigs API](https://api.configcat.com/docs/#operation/get-configs).

```
$ terraform import configcat_config.example 08d86d63-2726-47cd-8bfc-59608ecb91e2
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Used APIs
* [Get Config](https://api.configcat.com/docs/index.html#operation/get-config)
* [Create Config](https://api.configcat.com/docs/index.html#operation/create-config)
* [Update Config](https://api.configcat.com/docs/index.html#operation/update-config)
* [Delete Config](https://api.configcat.com/docs/index.html#operation/delete-config)
