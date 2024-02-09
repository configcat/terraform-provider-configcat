package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2BoolResource(t *testing.T) {
	const configId = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const testResourceName = "configcat_setting_value_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(false),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
					"percentage":     config.IntegerVariable(10),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionPercentage, "10"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionPercentage, "90"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionValue+"."+BoolValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
					"percentage":     config.IntegerVariable(10),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
					"percentage":     config.IntegerVariable(20),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionPercentage, "20"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionPercentage, "80"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionValue+"."+BoolValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
					"percentage":     config.IntegerVariable(20),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("with_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+BoolValue, "true"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "sensitiveTextEquals"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+StringValue, "@configcat.com"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("with_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+BoolValue, "true"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "sensitiveTextEquals"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+StringValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleValue+"."+BoolValue, "true"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "containsAnyOf"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueHint, "the greatest company of the world"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueValue, "@example.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueHint, ""),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("with_rules_and_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonAttribute, "color"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparator, "isOneOf"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonValue, "red"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemPercentage, "30"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemPercentage, "70"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules_and_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.BoolVariable(true),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
		},
	})
}
