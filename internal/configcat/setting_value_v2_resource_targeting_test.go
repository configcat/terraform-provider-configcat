package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2TargetingResource(t *testing.T) {
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
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "dateTimeBefore"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+DoubleValue, "1707752290.1"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			/*{
				TODO: This test case is valid, but the API should be fixed to fail for this. Uncomment when ready.
				ConfigFile: config.TestNameFile("comparison_value_type_mismatch1.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("must be a list"),
			},*/
			/*{
				TODO: This test case is valid, but the API should be fixed to fail for this. Uncomment when ready.
				ConfigFile: config.TestNameFile("comparison_value_type_mismatch2.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("must be a list"),
			},*/
			{
				ConfigFile: config.TestNameFile("comparison_value_type_mismatch3.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("must be a list"),
			},
			{
				ConfigFile: config.TestNameFile("comparison_value_missing.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("exactly one of the string_value, double_value or list_values attributes"),
			},
			{
				ConfigFile: config.TestNameFile("invalid_comparator.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("invalid value 'invalid' for UserComparator"),
			},
			{
				ConfigFile: config.TestNameFile("invalid_condition.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("exactly one of the user_condition, segment_condition or"),
			},
			{
				ConfigFile: config.TestNameFile("missing_targeting_rule_value_or_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("the value or the percentage_options attributes must be provided"),
			},
			{
				ConfigFile: config.TestNameFile("invalid_targeting_rule_value.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("Invalid term value."),
			},
			{
				ConfigFile: config.TestNameFile("invalid_percentage_option_value.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
				ExpectError: regexp.MustCompile("Invalid percentage option value."),
			},
		},
	})
}
