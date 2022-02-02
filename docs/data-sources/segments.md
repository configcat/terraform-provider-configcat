# configcat_segments Resource

Use this data source to access information about existing **Segments**.

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

data "configcat_segments" "my_segments" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name_filter_regex = "Test"
}


output "segment_id" {
  value = data.configcat_segments.my_segments.segments.0.segment_id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name_filter_regex` - (Optional) Filter the Segments by name.

## Attribute Reference

* `segments` - A segment [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `segments` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `segment_id` - The unique Segment ID.
* `name` - The name of the Segment.
* `description` - The description of the Segment.

## Endpoints used
- [Get Segments](https://api.configcat.com/docs/index.html#operation/get-segments)
