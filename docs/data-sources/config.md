# configcat_config Resource

Use this data source to access information about an existing Config. [Used API](https://api.configcat.com/docs/index.html#operation/get-configs)

## Example Usage

```hcl
data "configcat_product" "product" {
  name = "ConfigCat's product"
}

data "configcat_config" "product" {
  product_id = configcat_product.product.id
  name = "Main Config"
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product
* `name` - (Required) The name of the Config

## Attribute Reference

* `id` - The unique Config ID.