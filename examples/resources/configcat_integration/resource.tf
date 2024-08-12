variable "product_id" {
  type = string
}

variable "config_id" {
  type = string
}

variable "environment_id" {
  type = string
}

# Configure connected configs and environments
resource "configcat_integration" "slack_integration_with_configs_and_environments" {
  product_id       = var.product_id
  integration_type = "slack"
  name             = "Slack integration"
  parameters = {
    "incoming_webhook.url" = "https://mycompany.slack.com/services/B0000000000"
  }

  configs      = [var.config_id]
  environments = [var.environment_id]
}

# Slack integration (https://configcat.com/docs/integrations/slack/)
resource "configcat_integration" "slack_integration" {
  product_id       = var.product_id
  integration_type = "slack"
  name             = "Slack integration"
  parameters = {
    "incoming_webhook.url" = "https://mycompany.slack.com/services/B0000000000" # The incoming webhook URL where the integration should post messages. Read more at https://api.slack.com/messaging/webhooks.
  }
}

# Datadog integration (https://configcat.com/docs/integrations/datadog/)
resource "configcat_integration" "datadog_integration" {
  product_id       = var.product_id
  integration_type = "dataDog"
  name             = "Datadog integration"
  parameters = {
    "apikey" = ""   # Datadog API key. Read more at https://docs.datadoghq.com/account_management/api-app-keys/#api-keys
    "site"   = "Us" # Optional. Datadog site. Available values: Us, Eu, Us1Fed, Us3, Us5. Default: Us. Read more at https://docs.datadoghq.com/getting_started/site/.
  }
}

# Amplitude integration (https://configcat.com/docs/integrations/amplitude/#annotations)
resource "configcat_integration" "amplitude_integration" {
  product_id       = var.product_id
  integration_type = "amplitude"
  name             = "Amplitude integration"
  parameters = {
    "apikey"    = "" # Amplitude API Key. Read more at https://amplitude.com/docs/apis/authentication.
    "secretKey" = "" # Amplitude Secret Key. Read more at https://amplitude.com/docs/apis/authentication.
  }
}

# Mixpanel integration (https://configcat.com/docs/integrations/mixpanel/#annotations)
resource "configcat_integration" "mixpanel_integration" {
  product_id       = var.product_id
  integration_type = "mixPanel"
  name             = "Mixpanel integration"
  parameters = {
    "serviceAccountUserName" = ""               # Mixpanel Service Account Username.
    "serviceAccountSecret"   = ""               # Mixpanel Service Account Secret.
    "projectId"              = ""               # Mixpanel Project ID.
    "server"                 = "StandardServer" # Mixpanel Server. Available values: StandardServer, EUResidencyServer. Default: StandardServer. Read more at https://docs.mixpanel.com/docs/privacy/eu-residency.
  }
}

# Segment integration (https://configcat.com/docs/integrations/segment/#changeevents)
resource "configcat_integration" "mixpanel_integration" {
  product_id       = var.product_id
  integration_type = "segment"
  name             = "Mixpanel integration"
  parameters = {
    "writeKey" = ""               # Twilio Segment Write Key.
    "server"   = "StandardServer" # Twilio Segment Server. Available values: Us, Eu. Default: Us. Read more at https://segment.com/docs/guides/regional-segment/.
  }
}

output "slack_integration_id" {
  value = configcat_integration.slack_integration.id
}
