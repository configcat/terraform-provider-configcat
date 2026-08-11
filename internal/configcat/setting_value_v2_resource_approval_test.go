package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2ApprovalResource(t *testing.T) {
	const testResourceName = "configcat_setting_value_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Prepare: create product, config, environment, setting and set approve_required = true
				ConfigFile: config.TestNameFile("prepare.tf"),
			},
			{
				// Without bypass_approval, update should fail because approval is required
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"value": config.BoolVariable(true),
				},
				ExpectError: regexp.MustCompile("Approval required."),
			},
			{
				// With bypass_approval = true, update should succeed
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"value":           config.BoolVariable(true),
					"bypass_approval": config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
				),
			},
		},
	})
}
