package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueDoubleResource(t *testing.T) {
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
					"value":          config.StringVariable("2.31"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "2.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("2"),
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
					"value":          config.StringVariable("3.31"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "3.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("3.31"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("4.31"),
					"percentage1":      config.StringVariable("5"),
					"percentage2":      config.StringVariable("6"),
					"percentage3":      config.StringVariable("89"),
					"percentage1value": config.StringVariable("6.31"),
					"percentage2value": config.StringVariable("7.31"),
					"percentage3value": config.StringVariable("8.31"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "4.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemPercentage, "5"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemValue, "6.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemPercentage, "6"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemValue, "7.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".2."+RolloutPercentageItemPercentage, "89"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".2."+RolloutPercentageItemValue, "8.31"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("4.31"),
					"percentage1":      config.StringVariable("5"),
					"percentage2":      config.StringVariable("6"),
					"percentage3":      config.StringVariable("89"),
					"percentage1value": config.StringVariable("6.31"),
					"percentage2value": config.StringVariable("7.31"),
					"percentage3value": config.StringVariable("8.31"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("5.31"),
					"percentage1":      config.StringVariable("10"),
					"percentage2":      config.StringVariable("20"),
					"percentage3":      config.StringVariable("70"),
					"percentage1value": config.StringVariable("10.31"),
					"percentage2value": config.StringVariable("11.31"),
					"percentage3value": config.StringVariable("12.31"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "5.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemPercentage, "10"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemValue, "10.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemPercentage, "20"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemValue, "11.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".2."+RolloutPercentageItemPercentage, "70"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".2."+RolloutPercentageItemValue, "12.31"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("5.31"),
					"percentage1":      config.StringVariable("10"),
					"percentage2":      config.StringVariable("20"),
					"percentage3":      config.StringVariable("70"),
					"percentage1value": config.StringVariable("10.31"),
					"percentage2value": config.StringVariable("11.31"),
					"percentage3value": config.StringVariable("12.31"),
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "20.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "30.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "4.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "5.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonAttribute, "color"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparator, "isOneOf"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonValue, "red"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleValue, "6.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "40.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".0."+RolloutRuleValue, "50.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonAttribute, "color"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparator, "isOneOf"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleComparisonValue, "red"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".1."+RolloutRuleValue, "60.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemPercentage, "10"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".0."+RolloutPercentageItemValue, "70.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemPercentage, "20"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".1."+RolloutPercentageItemValue, "80.31"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".2."+RolloutPercentageItemPercentage, "70"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".2."+RolloutPercentageItemValue, "90.31"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules_and_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("lastvalue"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
