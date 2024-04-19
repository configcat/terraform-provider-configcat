package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductProductPreferencesResource(t *testing.T) {
	const testResourceName = "configcat_product_preferences.preferences"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("empty.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "camelCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".#", "0"),
				),
			},
			{
				ConfigFile:        config.TestNameFile("empty.tf"),
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("empty.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "camelCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".#", "0"),
				),
			},
			{
				ConfigFile:        config.TestNameFile("empty.tf"),
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("withsettings.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("camelCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "camelCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".#", "0"),
				),
			},
			{
				ConfigFile:   config.TestNameFile("withsettings.tf"),
				ResourceName: testResourceName,
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("camelCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("withsettings.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("pascalCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(true),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "pascalCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "true"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".#", "0"),
				),
			},
			{
				ConfigFile:   config.TestNameFile("withsettings.tf"),
				ResourceName: testResourceName,
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("pascalCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(true),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("withsettings.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(true),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "lowerCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "true"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "true"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".#", "0"),
				),
			},
			{
				ConfigFile:   config.TestNameFile("withsettings.tf"),
				ResourceName: testResourceName,
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(true),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(true),
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
