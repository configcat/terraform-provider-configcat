# configcat_config Resource

Creates and manages a **Config**. [What is a Config in ConfigCat?](https://configcat.com/docs/main-concepts)

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

resource "configcat_config" "my_config" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "My config"
  description = "My config description"
}


output "config_id" {
  value = configcat_config.my_config.id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name` - (Required) The name of the Config.
* `description` - (Optional) The description of the Config.

## Attribute Reference

* `id` - The unique Config ID.

## Import

Configs can be imported using the ConfigId. Get the ConfigId using the [List Configs API](https://api.configcat.com/docs/#tag/Configs/operation/get-configs) for example.

```
$ terraform import configcat_config.example 08d86d63-2726-47cd-8bfc-59608ecb91e2
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Config](https://api.configcat.com/docs/#tag/Configs/operation/get-config)
* [Create Config](https://api.configcat.com/docs/#tag/Configs/operation/create-config)
* [Update Config](https://api.configcat.com/docs/#tag/Configs/operation/update-config)
* [Delete Config](https://api.configcat.com/docs/#tag/Configs/operation/delete-config)
