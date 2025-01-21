variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

variable "url" {
  type = string
}

variable "http_method" {
  type    = string
  default = null
}

variable "content" {
  type    = string
  default = null
}

variable "webhook_headers" {
  type = list(object({
    key   = string
    value = string
  }))
  default = null
}

variable "secure_webhook_headers" {
  type = list(object({
    key   = string
    value = string
  }))
  default = null
}

resource "configcat_webhook" "test" {
  config_id      = var.config_id
  environment_id = var.environment_id

  url                    = var.url
  http_method            = var.http_method
  content                = var.content
  webhook_headers        = var.webhook_headers
  secure_webhook_headers = var.secure_webhook_headers
}
