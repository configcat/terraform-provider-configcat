package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueInitOnlyResource(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const testResourceName = "configcat_setting_value.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("false"),
					"init_only":      config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"init_only":      config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "true"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("false"),
					"init_only":      config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"init_only":      config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "true"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"init_only":      config.BoolVariable(true),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"init_only":      config.BoolVariable(false),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
		},
	})
}
