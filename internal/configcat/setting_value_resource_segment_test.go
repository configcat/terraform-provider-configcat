package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueSegmentResource(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const testResourceName = "configcat_setting_value.test"
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
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentComparator, "isIn"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentId, segmentBetaUsersID),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, ""),
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
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentComparator, "isNotIn"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentId, segmentOldApplicationVersionID),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, ""),
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
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentComparator, "isIn"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentId, segmentBetaUsersID),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, ""),

					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleSegmentComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleSegmentId, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleValue, "true"),

					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleSegmentComparator, "isNotIn"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleSegmentId, segmentOldApplicationVersionID),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleComparisonAttribute, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleComparisonValue, ""),
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
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentComparator, "isNotIn"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleSegmentId, segmentOldApplicationVersionID),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, ""),

					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleSegmentComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleSegmentId, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleValue, "true"),

					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleSegmentComparator, "isIn"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleSegmentId, segmentBetaUsersID),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleComparisonAttribute, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleComparator, ""),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".2."+RolloutRuleComparisonValue, ""),
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
