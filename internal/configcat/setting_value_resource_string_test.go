package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueStringResource(t *testing.T) {
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
					"value":          config.StringVariable("Vvalue"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "Vvalue"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("Vvalue"),
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
					"value":          config.StringVariable("VvalueMod"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "VvalueMod"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("value_only.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("VvalueMod"),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("Vvalue"),
					"percentage1":      config.StringVariable("5"),
					"percentage2":      config.StringVariable("6"),
					"percentage3":      config.StringVariable("89"),
					"percentage1value": config.StringVariable("p1value"),
					"percentage2value": config.StringVariable("p2value"),
					"percentage3value": config.StringVariable("p3value"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "Vvalue"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[0]."+RolloutPercentageItemPercentage, "5"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[0]."+RolloutPercentageItemValue, "p1value"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[1]."+RolloutPercentageItemPercentage, "6"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[1]."+RolloutPercentageItemValue, "p2value"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[2]."+RolloutPercentageItemPercentage, "89"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[2]."+RolloutPercentageItemValue, "p3value"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("Vvalue"),
					"percentage1":      config.StringVariable("5"),
					"percentage2":      config.StringVariable("6"),
					"percentage3":      config.StringVariable("89"),
					"percentage1value": config.StringVariable("p1value"),
					"percentage2value": config.StringVariable("p2value"),
					"percentage3value": config.StringVariable("p3value"),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("Vvaluemod"),
					"percentage1":      config.StringVariable("10"),
					"percentage2":      config.StringVariable("20"),
					"percentage3":      config.StringVariable("70"),
					"percentage1value": config.StringVariable("p1valuemod"),
					"percentage2value": config.StringVariable("p2valuemod"),
					"percentage3value": config.StringVariable("p3valuemod"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "Vvaluemod"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "0"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[0]."+RolloutPercentageItemPercentage, "10"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[0]."+RolloutPercentageItemValue, "p1valuemod"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[1]."+RolloutPercentageItemPercentage, "20"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[1]."+RolloutPercentageItemValue, "p2valuemod"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[2]."+RolloutPercentageItemPercentage, "70"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[2]."+RolloutPercentageItemValue, "p3valuemod"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":        config.StringVariable(configId),
					"environment_id":   config.StringVariable(environmentId),
					"value":            config.StringVariable("Vvaluemod"),
					"percentage1":      config.StringVariable("10"),
					"percentage2":      config.StringVariable("20"),
					"percentage3":      config.StringVariable("70"),
					"percentage1value": config.StringVariable("p1valuemod"),
					"percentage2value": config.StringVariable("p2valuemod"),
					"percentage3value": config.StringVariable("p3valuemod"),
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
					"value":          config.StringVariable("Vvalue"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "Vvalue"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleValue, "RRvalue"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rule.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("Vvalue"),
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
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "Vvalue"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleValue, "RR1value"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleComparisonAttribute, "color"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleComparator, "isOneOf"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleComparisonValue, "red"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleValue, "RR2value"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "0"),
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, SettingValue, "Vvalue"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+".#", "2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparisonAttribute, "email"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparator, "contains"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleComparisonValue, "@configcat.com"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[0]."+RolloutRuleValue, "RRValue1"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleComparisonAttribute, "color"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleComparator, "isOneOf"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleComparisonValue, "red"),
					resource.TestCheckResourceAttr(testResourceName, RolloutRules+"[1]."+RolloutRuleValue, "RRValue2"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[0]."+RolloutPercentageItemPercentage, "10"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[0]."+RolloutPercentageItemValue, "P1value"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[1]."+RolloutPercentageItemPercentage, "20"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[1]."+RolloutPercentageItemValue, "P2value"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[2]."+RolloutPercentageItemPercentage, "70"),
					resource.TestCheckResourceAttr(testResourceName, RolloutPercentageItems+"[2]."+RolloutPercentageItemValue, "P3value"),
				),
			},
			{
				ConfigFile: config.TestNameFile("with_rules_and_percentage.tf"),
				ConfigVariables: config.Variables{
					"config_id":      config.StringVariable(configId),
					"environment_id": config.StringVariable(environmentId),
					"value":          config.StringVariable("lastvalue"),
				},
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
		},
	})
}
