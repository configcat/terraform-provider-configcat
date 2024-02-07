package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSdkKeysDataSource(t *testing.T) {
	const testResourceName = "data.configcat_sdkkeys.test"

	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const environment1Id = "08d86d63-272c-4355-8027-4b52787bc1bd"
	const environment2Id = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const environment3Id = "00000000-0000-0000-0000-000000000000" // not found

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environment1Id),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, PrimarySdkKey, "Y23YCCEnpk2MBlhFIdUWvA/nJmlkmynTE-t3GUMoJjOAQ"),
					resource.TestCheckNoResourceAttr(testResourceName, SecondarySdkKey),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environment2Id),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, PrimarySdkKey, "Y23YCCEnpk2MBlhFIdUWvA/k-nG8bhE10K41sa8C2rYOQ"),
					resource.TestCheckResourceAttr(testResourceName, SecondarySdkKey, "Y23YCCEnpk2MBlhFIdUWvA/q7dBIG43LkuX8NJbCeAzdg"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environment3Id),
				},
				ExpectError: regexp.MustCompile(`Unable to read SDK Keys, got error: 404 Not Found`),
			},
		},
	})
}
