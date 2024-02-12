package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2PrerequisiteResource(t *testing.T) {
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, DefaultValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+BoolValue, "true"),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparator, "equals"),
					resource.TestCheckResourceAttrPair(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionSettingId, "configcat_setting.bool", ID),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparisonValue+"."+BoolValue, "true"),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".1."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparator, "doesNotEqual"),
					resource.TestCheckResourceAttrPair(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".1."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionSettingId, "configcat_setting.string", ID),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".1."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparisonValue+"."+StringValue, "test"),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleValue+"."+BoolValue, "false"),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparator, "equals"),
					resource.TestCheckResourceAttrPair(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionSettingId, "configcat_setting.int", ID),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparisonValue+"."+IntValue, "1"),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".1."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparator, "doesNotEqual"),
					resource.TestCheckResourceAttrPair(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".1."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionSettingId, "configcat_setting.double", ID),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".1."+TargetingRulePrerequisiteFlagCondition+"."+TargetingRulePrerequisiteFlagConditionComparisonValue+"."+DoubleValue, "1.1"),
				),
			},

			{
				ConfigFile: config.TestNameFile("bool_error.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("Invalid prerequisite value."),
			},

			{
				ConfigFile: config.TestNameFile("string_error.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("Invalid prerequisite value."),
			},

			{
				ConfigFile: config.TestNameFile("int_error.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("Invalid prerequisite value."),
			},

			{
				ConfigFile: config.TestNameFile("double_error.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("Invalid prerequisite value."),
			},

			{
				ConfigFile: config.TestNameFile("comparator_error.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("invalid value 'invalid' for PrerequisiteComparator"),
			},

			{
				ConfigFile: config.TestNameFile("setting_id_error.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("Prerequisite setting not found."),
			},

			{
				ConfigFile: config.TestNameFile("cleanup.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
			},
		},
	})
}
