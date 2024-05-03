# Webhooks can be imported using the WebhookId. Get the WebhookId using the [List Webhooks API](https://api.configcat.com/docs/index.html#tag/Webhooks/operation/get-webhooks) for example.
# It is important to note that webhooks containing secure webhook headers cannot be imported via `terraform import`.
# If you want to manage your webhooks that already contain secure webhook headers, you should create brand new configcat_webhook resources in Terraform without importing them. After they are created successfully and managed by Terraform, you can safely delete the old, non Terraform managed webhook from the ConfigCat Dashboard.

terraform import configcat_webhook.example 1234