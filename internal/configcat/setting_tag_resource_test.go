package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccSettingTagResource(t *testing.T) {
	const setting1Id = "67639"
	const setting2Id = "167364"
	const tagId = "46"
	const testResourceName = "configcat_setting_tag.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"setting_id": config.StringVariable(setting1Id),
					"tag_id":     config.StringVariable(tagId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, ID, setting1Id+":"+tagId),
					resource.TestCheckResourceAttr(testResourceName, SettingId, setting1Id),
					resource.TestCheckResourceAttr(testResourceName, TagId, tagId),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"setting_id": config.StringVariable(setting1Id),
					"tag_id":     config.StringVariable(tagId),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Attribute change should recreate resource
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"setting_id": config.StringVariable(setting2Id),
					"tag_id":     config.StringVariable(tagId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, ID, setting2Id+":"+tagId),
					resource.TestCheckResourceAttr(testResourceName, SettingId, setting2Id),
					resource.TestCheckResourceAttr(testResourceName, TagId, tagId),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testResourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
			{
				// Removing the resource should work fine.
				ConfigFile: config.TestNameFile("cleanup.tf"),
			},
		},
	})
}

func TestAccSettingTagMultipleResource(t *testing.T) {
	const configId = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResource1Name = "configcat_setting_tag.settingTag1"
	const testResource2Name = "configcat_setting_tag.settingTag2"
	const testResource3Name = "configcat_setting_tag.settingTag3"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("init.tf"),
				ConfigVariables: config.Variables{
					"config_id":  config.StringVariable(configId),
					"product_id": config.StringVariable(productId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResource1Name, ID),
					resource.TestCheckResourceAttrSet(testResource1Name, SettingId),
					resource.TestCheckResourceAttrSet(testResource1Name, TagId),
					resource.TestCheckResourceAttrSet(testResource2Name, ID),
					resource.TestCheckResourceAttrSet(testResource2Name, SettingId),
					resource.TestCheckResourceAttrSet(testResource2Name, TagId),
					resource.TestCheckResourceAttrSet(testResource3Name, ID),
					resource.TestCheckResourceAttrSet(testResource3Name, SettingId),
					resource.TestCheckResourceAttrSet(testResource3Name, TagId),
				),
			},
			{
				ConfigFile: config.TestNameFile("removeonetag.tf"),
				ConfigVariables: config.Variables{
					"config_id":  config.StringVariable(configId),
					"product_id": config.StringVariable(productId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResource1Name, ID),
					resource.TestCheckResourceAttrSet(testResource1Name, SettingId),
					resource.TestCheckResourceAttrSet(testResource1Name, TagId),
					resource.TestCheckResourceAttrSet(testResource3Name, ID),
					resource.TestCheckResourceAttrSet(testResource3Name, SettingId),
					resource.TestCheckResourceAttrSet(testResource3Name, TagId),
				),
			},
			{
				ConfigFile: config.TestNameFile("removeeverything.tf"),
			},
		},
	})
}
