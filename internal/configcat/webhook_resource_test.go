package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccWebhookResource(t *testing.T) {
	const config_id = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const environment_id = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const test_url = "https://test.example.com"
	const test_url2 = "https://test2.example.com"
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
					"url":            config.StringVariable(test_url2),
					"http_method":    config.StringVariable("post"),
					"content":        config.StringVariable("test content"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, WebhookUrl, test_url2),
					resource.TestCheckResourceAttr(testResourceName, WebhookHttpMethod, "post"),
					resource.TestCheckResourceAttr(testResourceName, WebhookContent, "test content"),
					resource.TestCheckNoResourceAttr(testResourceName, WebhookHeaders),
					resource.TestCheckNoResourceAttr(testResourceName, SecureWebhookHeaders),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url2),
					"http_method":    config.StringVariable("post"),
					"content":        config.StringVariable("test content"),
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
					"url":            config.StringVariable(test_url2),
					"http_method":    config.StringVariable("post"),
					"content":        config.StringVariable("test content"),
					"webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey1"),
							"value": config.StringVariable("whvalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey2"),
							"value": config.StringVariable("whvalue2"),
						}),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, WebhookUrl, test_url2),
					resource.TestCheckResourceAttr(testResourceName, WebhookHttpMethod, "post"),
					resource.TestCheckResourceAttr(testResourceName, WebhookContent, "test content"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".0."+WebhookHeaderKey, "whkey1"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".0."+WebhookHeaderValue, "whvalue1"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".1."+WebhookHeaderKey, "whkey2"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".1."+WebhookHeaderValue, "whvalue2"),
					resource.TestCheckNoResourceAttr(testResourceName, SecureWebhookHeaders),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url2),
					"http_method":    config.StringVariable("post"),
					"content":        config.StringVariable("test content"),
					"webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey1"),
							"value": config.StringVariable("whvalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey2"),
							"value": config.StringVariable("whvalue2"),
						}),
					),
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
					"url":            config.StringVariable(test_url2),
					"http_method":    config.StringVariable("post"),
					"content":        config.StringVariable("test content"),
					"webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey1"),
							"value": config.StringVariable("whvalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey2"),
							"value": config.StringVariable("whvalue2"),
						}),
					),
					"secure_webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whsecurekey1"),
							"value": config.StringVariable("whsecurevalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whsecurekey2"),
							"value": config.StringVariable("whsecurevalue2"),
						}),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, WebhookUrl, test_url2),
					resource.TestCheckResourceAttr(testResourceName, WebhookHttpMethod, "post"),
					resource.TestCheckResourceAttr(testResourceName, WebhookContent, "test content"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".0."+WebhookHeaderKey, "whkey1"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".0."+WebhookHeaderValue, "whvalue1"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".1."+WebhookHeaderKey, "whkey2"),
					resource.TestCheckResourceAttr(testResourceName, WebhookHeaders+".1."+WebhookHeaderValue, "whvalue2"),
					resource.TestCheckResourceAttr(testResourceName, SecureWebhookHeaders+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, SecureWebhookHeaders+".0."+WebhookHeaderKey, "whsecurekey1"),
					resource.TestCheckResourceAttr(testResourceName, SecureWebhookHeaders+".0."+WebhookHeaderValue, "whsecurevalue1"),
					resource.TestCheckResourceAttr(testResourceName, SecureWebhookHeaders+".1."+WebhookHeaderKey, "whsecurekey2"),
					resource.TestCheckResourceAttr(testResourceName, SecureWebhookHeaders+".1."+WebhookHeaderValue, "whsecurevalue2"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url2),
					"http_method":    config.StringVariable("post"),
					"content":        config.StringVariable("test content"),
					"webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey1"),
							"value": config.StringVariable("whvalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey2"),
							"value": config.StringVariable("whvalue2"),
						}),
					),
					"secure_webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whsecurekey1"),
							"value": config.StringVariable("whsecurevalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whsecurekey2"),
							"value": config.StringVariable("whsecurevalue2"),
						}),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(config_id),
					"environment_id": config.StringVariable(environment_id),
					"url":            config.StringVariable(test_url2),
					"http_method":    config.StringVariable("post"),
					"content":        config.StringVariable("test content"),
					"webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey1"),
							"value": config.StringVariable("whvalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whkey2"),
							"value": config.StringVariable("whvalue2"),
						}),
					),
					"secure_webhook_headers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whsecurekey1"),
							"value": config.StringVariable("whsecurevalue1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"key":   config.StringVariable("whsecurekey2"),
							"value": config.StringVariable("whsecurevalue2"),
						}),
					),
				},
				ResourceName: testResourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile("Unable to read secure webhook headers."),
			},
		},
	})
}
