# configcat_tags Resource

Use this data source to access information about existing **Tags**.
## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_tags" "my_tags" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Test"
}


output "tag_id" {
  value = data.configcat_tags.my_tags.tags.0.tag_id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name_filter_regex` - (Optional) Filter the Tags by name.

## Attribute Reference

* `tags` - A tag [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `tags` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `tag_id` - The unique Tag ID.
* `name` - The name of the Tag.
* `color` - The color of the Tag.

## Used APIs
- [Get Tags](https://api.configcat.com/docs/index.html#operation/get-tags)
