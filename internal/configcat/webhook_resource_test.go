package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookResource(t *testing.T) {
	const config_id = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const environment_id = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const test_url = "https://test.example.com"
	const wrong_url = "asdasd"
	const testResourceName = "configcat_webhook.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(wrong_url),
				},
				ExpectError: regexp.MustCompile("Url is invalid"),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url),
					"http_method":    config.StringVariable("patch"),
				},
				ExpectError: regexp.MustCompile("Invalid http_method"),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, WebhookUrl, test_url),
					resource.TestCheckResourceAttr(testResourceName, WebhookHttpMethod, "get"),
					resource.TestCheckNoResourceAttr(testResourceName, WebhookContent),
					resource.TestCheckNoResourceAttr(testResourceName, WebhookHeaders),
					resource.TestCheckNoResourceAttr(testResourceName, SecureWebhookHeaders),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, WebhookUrl, test_url),
					resource.TestCheckResourceAttr(testResourceName, WebhookHttpMethod, "get"),
					resource.TestCheckNoResourceAttr(testResourceName, WebhookContent),
					resource.TestCheckNoResourceAttr(testResourceName, WebhookHeaders),
					resource.TestCheckNoResourceAttr(testResourceName, SecureWebhookHeaders),
				),
			},
		},
	})
}
