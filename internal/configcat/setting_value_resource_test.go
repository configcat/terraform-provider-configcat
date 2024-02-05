package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueBoolResource(t *testing.T) {
	const configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const testResourceName = "configcat_setting_value.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("false"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "false"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("false"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"percentage":     config.StringVariable("10"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemPercentage, "10"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemPercentage, "90"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"percentage":     config.StringVariable("10"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"percentage":     config.StringVariable("20"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemPercentage, "20"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemPercentage, "80"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemValue, "false"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
					"percentage":     config.StringVariable("20"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("with_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "true"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("with_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
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
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("with_rules_and_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("true"),
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
					"value":          config.StringVariable("true"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
