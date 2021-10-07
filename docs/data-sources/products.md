# configcat_products Resource

Use this data source to access information about existing **Products**. [What is a Product in ConfigCat?](https://configcat.com/docs/main-concepts)

## Example Usage

```hcl
data "configcat_products" "my_products" {
  name_filter_regex = "ConfigCat's product"
}


output "product_id" {
  value = data.configcat_products.my_products.products.0.product_id
}
```

## Argument Reference

* `name_filter_regex` - (Optional) Filter the Products by name.

## Attribute Reference

* `products` - A product [list](https://www.terraform.io/docs/configuration/types.html#list-) block defined as below.

### The `products` [list](https://www.terraform.io/docs/configuration/types.html#list-) block

* `product_id` - The unique Product ID.
* `name` - The name of the Product.
* `description` - The description of the Product.

## Endpoints used
- [Get Products](https://api.configcat.com/docs/index.html#operation/get-products)