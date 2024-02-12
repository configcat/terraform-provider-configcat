package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2SegmentResource(t *testing.T) {
	const configId = "08dc1bfa-b8b0-45f0-8127-fac0de7a37ac"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const testResourceName = "configcat_setting_value_v2.test"
	const segmentBetaUsersID = "08d9e65b-a9f2-48ec-8423-9b2224771639"
	const segmentOldApplicationVersionID = "08d9e65b-c95e-4cb1-8a24-a326d1da676f"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// prepare
				ConfigFile: config.TestNameFile("base.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
				},
			},
			{
				ConfigFile: config.TestNameFile("one_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator":     config.StringVariable("isIn"),
					"segment_id":     config.StringVariable(segmentBetaUsersID),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+BoolValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentComparator, "isIn"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentId, segmentBetaUsersID),
				),
			},
			{
				ConfigFile: config.TestNameFile("one_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator":     config.StringVariable("isIn"),
					"segment_id":     config.StringVariable(segmentBetaUsersID),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("one_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator":     config.StringVariable("isNotIn"),
					"segment_id":     config.StringVariable(segmentOldApplicationVersionID),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+BoolValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentComparator, "isNotIn"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentId, segmentOldApplicationVersionID),
				),
			},
			{
				ConfigFile: config.TestNameFile("one_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator":     config.StringVariable("isNotIn"),
					"segment_id":     config.StringVariable(segmentOldApplicationVersionID),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("one_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator":     config.StringVariable("invalid"),
					"segment_id":     config.StringVariable(segmentOldApplicationVersionID),
				},
				ExpectError: regexp.MustCompile("invalid value 'invalid' for SegmentComparator"),
			},
			{
				ConfigFile: config.TestNameFile("multiple_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator_1":   config.StringVariable("isIn"),
					"segment_id_1":   config.StringVariable(segmentBetaUsersID),
					"comparator_2":   config.StringVariable("isNotIn"),
					"segment_id_2":   config.StringVariable(segmentOldApplicationVersionID),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+BoolValue, "false"),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentComparator, "isIn"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentId, segmentBetaUsersID),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "textEquals"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+StringValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".1."+TargetingRuleSegmentCondition+"."+SegmentComparator, "isNotIn"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".1."+TargetingRuleSegmentCondition+"."+SegmentId, segmentOldApplicationVersionID),
				),
			},
			{
				ConfigFile: config.TestNameFile("multiple_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator_1":   config.StringVariable("isIn"),
					"segment_id_1":   config.StringVariable(segmentBetaUsersID),
					"comparator_2":   config.StringVariable("isNotIn"),
					"segment_id_2":   config.StringVariable(segmentOldApplicationVersionID),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},

			{
				ConfigFile: config.TestNameFile("multiple_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator_1":   config.StringVariable("isNotIn"),
					"segment_id_1":   config.StringVariable(segmentOldApplicationVersionID),
					"comparator_2":   config.StringVariable("isIn"),
					"segment_id_2":   config.StringVariable(segmentBetaUsersID),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleValue+"."+BoolValue, "false"),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentComparator, "isNotIn"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".0."+TargetingRuleConditions+".0."+TargetingRuleSegmentCondition+"."+SegmentId, segmentOldApplicationVersionID),

					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleValue+"."+BoolValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparator, "textEquals"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".0."+TargetingRuleUserCondition+"."+TargetingRuleUserConditionComparisonValue+"."+StringValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".1."+TargetingRuleSegmentCondition+"."+SegmentComparator, "isIn"),
					resource.TestCheckResourceAttr(testResourceName, TargetingRules+".1."+TargetingRuleConditions+".1."+TargetingRuleSegmentCondition+"."+SegmentId, segmentBetaUsersID),
				),
			},
			{
				ConfigFile: config.TestNameFile("multiple_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"comparator_1":   config.StringVariable("isIn"),
					"segment_id_1":   config.StringVariable(segmentBetaUsersID),
					"comparator_2":   config.StringVariable("isNotIn"),
					"segment_id_2":   config.StringVariable(segmentOldApplicationVersionID),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
		},
	})
}
