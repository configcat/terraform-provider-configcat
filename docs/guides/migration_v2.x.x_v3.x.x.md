---
page_title: "Migration guide from v2.x.x to v3.x.x"
---

# Migration guide from v2.x.x to v3.x.x

## Breaking change in v3.x.x

We have introduced a new required property in the `configcat_product`, `configcat_config`, `configcat_environment` and the `configcat_setting` resources: `order`.
The `order` is a zero-based number that specifies the order of the resource when displayed on the ConfigCat Dashboard.
If multiple resources has the same `order`, they are displayed in alphabetical order.

Please specify explicitly the `order` in your  `configcat_product`, `configcat_config`, `configcat_environment` and the `configcat_setting` resources after upgrading to v3.x.x.