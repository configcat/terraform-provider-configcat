package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingsDataSource(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const settingId = "67639"
	const testResourceName = "data.configcat_settings.test"

	const listAttribute = Settings

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id": config.StringVariable(configId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingId, settingId),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "My awesome feature flag"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingKey, "isAwesomeFeatureEnabled"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingHint, "This is the hint for my awesome feature flag"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingType, "boolean"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Order, "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"key_filter_regex": config.StringVariable("isAwesomeFeatureEnabled"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingId, settingId),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "My awesome feature flag"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingKey, "isAwesomeFeatureEnabled"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingHint, "This is the hint for my awesome feature flag"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+SettingType, "boolean"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Order, "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"key_filter_regex": config.StringVariable("notfound"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id": config.StringVariable("invalidguid"),
				},
				ExpectError: regexp.MustCompile(`Attribute config_id value must be a valid GUID, got: invalidguid`),
			},
		},
	})
}
