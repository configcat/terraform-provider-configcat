# configcat_segment Resource

Creates and manages a **Segment**.

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}

resource "configcat_segment" "my_segment" {
  product_id = data.configcat_products.my_products.products.0.product_id
  name = "Beta users"
  description = "Beta users' description"
  comparison_attribute = "email"
  comparator = "isOneOfSensitive"
  comparison_value = "betauser1@example.com,betauser2@example.com"
}


output "segment_id" {
  value = configcat_segment.my_segment.id
}
```

## Argument Reference

* `product_id` - (Required) The ID of the Product.
* `name` - (Required) The name of the Segment.
* `description` - (Optional) The description of the Segment.
* `comparison_attribute` - (Required) The [comparison attribute](https://configcat.com/docs/advanced/targeting/#attribute).
* `comparator` - (Required) The [comparator](https://configcat.com/docs/advanced/targeting/#comparator).
* `comparison_value` - (Required) The [comparison value](https://configcat.com/docs/advanced/targeting/#comparison-value).

## Attribute Reference

* `id` - The unique Segment ID.

## Import

Segments can be imported using the SegmentId. Get the SegmentId using the [GetSegments API](https://api.configcat.com/docs/#operation/get-segments) for example.

```
$ terraform import configcat_segment.example 08d86d63-2726-47cd-8bfc-59608ecb91e2
```

[Read more](https://learn.hashicorp.com/tutorials/terraform/state-import) about importing.

## Endpoints used
* [Get Segment](https://api.configcat.com/docs/index.html#operation/get-segment)
* [Create Segment](https://api.configcat.com/docs/index.html#operation/create-segment)
* [Update Segment](https://api.configcat.com/docs/index.html#operation/update-segment)
* [Delete Segment](https://api.configcat.com/docs/index.html#operation/delete-segment)
