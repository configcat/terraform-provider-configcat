# configcat_product Resource

Use this data source to access information about an existing Product. [Used API](https://api.configcat.com/docs/index.html#operation/get-products)

## Example Usage

```hcl
data "configcat_product" "product" {
  name = "ConfigCat's product"
}
```

## Argument Reference

* `name` - (Required) The name of the Product

## Attribute Reference

* `id` - The unique Product ID.