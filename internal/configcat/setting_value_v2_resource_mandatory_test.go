package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2MandatoryResource(t *testing.T) {
	const configId = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const mandatoryEnvironmentId = "08d8becf-d4d9-4c66-8b48-6ac74cd95fba"
	const testResourceName = "configcat_setting_value_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// prepare
				ConfigFile: config.TestNameFile("resource.tf"),
				ConfigVariables: config.Variables{
					"config_id": config.StringVariable(configId),
				},
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
				),
			},

			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(mandatoryEnvironmentId),
					"value":          config.BoolVariable(true),
				},
				ExpectError: regexp.MustCompile(".*Reason required.*"),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":       config.StringVariable(configId),
					"environment_id":  config.StringVariable(mandatoryEnvironmentId),
					"value":           config.BoolVariable(true),
					"mandatory_notes": config.StringVariable("mandatory"),
				},

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
				),
			},
		},
	})
}
