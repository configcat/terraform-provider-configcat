package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2DoubleResource(t *testing.T) {
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
					"value":          config.FloatVariable(1.1),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "1.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.FloatVariable(1.1),
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
					"value":          config.FloatVariable(2.1),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "2.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.FloatVariable(2.1),
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
					"value":          config.FloatVariable(3.1),
					"percentage":     config.IntegerVariable(10),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "3.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionPercentage, "10"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionValue+"."+DoubleValue, "10.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionPercentage, "90"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionValue+"."+DoubleValue, "11.1"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.FloatVariable(3.1),
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
					"value":          config.FloatVariable(4.1),
					"percentage":     config.IntegerVariable(20),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "4.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionPercentage, "20"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionValue+"."+DoubleValue, "10.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionPercentage, "80"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionValue+"."+DoubleValue, "11.1"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.FloatVariable(4.1),
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "20.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+DoubleValue, "21.1"),
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "40.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+DoubleValue, "41.1"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "sensitiveTextEquals"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+StringValue, "jane@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleValue+"."+DoubleValue, "42.1"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "containsAnyOf"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueHint, "the greatest company of the world"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueValue, "@example.com"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueHint),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
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
					"percentage":     config.IntegerVariable(27),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "30.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+DoubleValue, "31.1"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "sensitiveTextEquals"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+StringValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleValue+"."+DoubleValue, "32.1"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "color"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "isOneOf"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueValue, "#000000"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueHint, "black"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueValue, "red"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueHint),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionPercentage, "27"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionValue+"."+DoubleValue, "33.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionPercentage, "73"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionValue+"."+DoubleValue, "34.1"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules_and_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"percentage":     config.IntegerVariable(27),
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
					"percentage":     config.IntegerVariable(37),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+DoubleValue, "30.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+DoubleValue, "31.1"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "sensitiveTextEquals"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+StringValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRulePercentageOptions+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleValue+"."+DoubleValue, "32.1"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "color"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "isOneOf"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueValue, "#000000"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".0."+ListValueHint, "black"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueValue, "red"),
					resource.TestCheckNoResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+ListValues+".1."+ListValueHint),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionPercentage, "37"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".0."+TargetingRulePercentageOptionValue+"."+DoubleValue, "33.1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionPercentage, "63"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".2."+TargetingRulePercentageOptions+".1."+TargetingRulePercentageOptionValue+"."+DoubleValue, "34.1"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules_and_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"percentage":     config.IntegerVariable(37),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
		},
	})
}
