package configcat

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccProductProductPreferencesResource(t *testing.T) {
	const testResourceName = "configcat_product_preferences.preferences"
	const testEnvironmentIdResourceName = "configcat_environment.test"
	const prodEnvironmentIdResourceName = "configcat_environment.prod"

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
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "0"),
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
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "0"),
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
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "0"),
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
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "0"),
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
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "0"),
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
			{
				ConfigFile: config.TestNameFile("reasonrequiredenvironments.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(false),
					"prod_required":                       config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "lowerCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "2"),
					testAccCheckEnvironmentReasonRequired(testEnvironmentIdResourceName, testResourceName, "false"),
					testAccCheckEnvironmentReasonRequired(prodEnvironmentIdResourceName, testResourceName, "false"),
				),
			},
			{
				ConfigFile:   config.TestNameFile("reasonrequiredenvironments.tf"),
				ResourceName: testResourceName,
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(false),
					"prod_required":                       config.BoolVariable(false),
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("reasonrequiredenvironments.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(true),
					"prod_required":                       config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "lowerCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "2"),
					testAccCheckEnvironmentReasonRequired(testEnvironmentIdResourceName, testResourceName, "true"),
					testAccCheckEnvironmentReasonRequired(prodEnvironmentIdResourceName, testResourceName, "false"),
				),
			},
			{
				ConfigFile:   config.TestNameFile("reasonrequiredenvironments.tf"),
				ResourceName: testResourceName,
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(true),
					"prod_required":                       config.BoolVariable(false),
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("reasonrequiredenvironments.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(true),
					"prod_required":                       config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "lowerCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "2"),
					testAccCheckEnvironmentReasonRequired(testEnvironmentIdResourceName, testResourceName, "true"),
					testAccCheckEnvironmentReasonRequired(prodEnvironmentIdResourceName, testResourceName, "true"),
				),
			},
			{
				ConfigFile:   config.TestNameFile("reasonrequiredenvironments.tf"),
				ResourceName: testResourceName,
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(true),
					"prod_required":                       config.BoolVariable(true),
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("reasonrequiredenvironments.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(false),
					"prod_required":                       config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceKeyGenerationMode, "lowerCase"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceShowVariationId, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceMandatorySettingHint, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequired, "false"),
					resource.TestCheckResourceAttr(testResourceName, ProductPreferenceReasonRequiredEnvironmentments+".%", "2"),
					testAccCheckEnvironmentReasonRequired(testEnvironmentIdResourceName, testResourceName, "false"),
					testAccCheckEnvironmentReasonRequired(prodEnvironmentIdResourceName, testResourceName, "true"),
				),
			},
			{
				ConfigFile:   config.TestNameFile("reasonrequiredenvironments.tf"),
				ResourceName: testResourceName,
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(false),
					"test_required":                       config.BoolVariable(false),
					"prod_required":                       config.BoolVariable(true),
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("reasonrequiredenvironments.tf"),
				ConfigVariables: config.Variables{
					ProductPreferenceKeyGenerationMode:    config.StringVariable("lowerCase"),
					ProductPreferenceShowVariationId:      config.BoolVariable(false),
					ProductPreferenceMandatorySettingHint: config.BoolVariable(false),
					ProductPreferenceReasonRequired:       config.BoolVariable(true),
					"test_required":                       config.BoolVariable(true),
					"prod_required":                       config.BoolVariable(true),
				},
				ExpectError: regexp.MustCompile("Please set reason_required to true to require mandatory notes globally"),
			},
		},
	})
}

func testAccCheckEnvironmentReasonRequired(environmentResourceName string, productPreferenceResourceName string, expectedReasonRequired string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		environmentResource, ok := s.RootModule().Resources[environmentResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", environmentResourceName)
		}

		productPreferenceResource, ok := s.RootModule().Resources[productPreferenceResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", productPreferenceResourceName)
		}

		reasonRequired := productPreferenceResource.Primary.Attributes[ProductPreferenceReasonRequiredEnvironmentments+"."+environmentResource.Primary.ID]
		if reasonRequired != expectedReasonRequired {
			return fmt.Errorf("Invalid ReasonRequired for %s. Expected %s, got %s", environmentResourceName, expectedReasonRequired, reasonRequired)
		}

		return nil
	}
}
