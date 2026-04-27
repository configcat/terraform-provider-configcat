package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2ListValueResource(t *testing.T) {
	const configId = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const testResourceName = "configcat_setting_value_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"list_value":     config.StringVariable("TEST_BEFORE"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueValue, "TEST_BEFORE"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"list_value":     config.StringVariable("TEST_AFTER"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueValue, "TEST_AFTER"),
				),
			},
		},
	})
}
