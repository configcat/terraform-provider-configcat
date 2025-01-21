package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookSigningKeysDataSource(t *testing.T) {
	const testResourceName = "data.configcat_webhook_signing_keys.test"

	const webhook1Id = 627
	const webhook2Id = 632

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"webhook_id": config.IntegerVariable(webhook1Id),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, WebhookSigningKeyKey1, "configcat_whsk_g+cdMQSCS9yS9bFWSQtKBTGMzICcfOmbzYoNZJxYr6E="),
					resource.TestCheckNoResourceAttr(testResourceName, WebhookSigningKeyKey2),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"webhook_id": config.IntegerVariable(webhook2Id),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, WebhookSigningKeyKey1, "configcat_whsk_K6W9wnJ7LAnuOvfm9fXZP5XHoM0EsOey6hyZAmT3k9Q="),
					resource.TestCheckResourceAttr(testResourceName, WebhookSigningKeyKey2, "configcat_whsk_aJsqdCgiVAHySBtQxJlgUauQ3P+2ffZIGJOkf9iCvNE="),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"webhook_id": config.IntegerVariable(-1),
				},
				ExpectError: regexp.MustCompile(`Unable to read Webhook Signing Keys, got error: 404 Not Found`),
			},
		},
	})
}
