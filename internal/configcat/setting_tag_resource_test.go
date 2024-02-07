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
		},
	})
}
