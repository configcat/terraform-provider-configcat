package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingValueV2PredefinedVariationResource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const environmentId = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const boolSettingResourceName = "configcat_setting.bool_setting"
	const stringSettingResourceName = "configcat_setting.string_setting"
	const intSettingResourceName = "configcat_setting.int_setting"
	const doubleSettingResourceName = "configcat_setting.double_setting"
	const boolSettingValueResourceName = "configcat_setting_value_v2.bool_setting_value"
	const stringSettingValueResourceName = "configcat_setting_value_v2.string_setting_value"
	const intSettingValueResourceName = "configcat_setting_value_v2.int_setting_value"
	const doubleSettingValueResourceName = "configcat_setting_value_v2.double_setting_value"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("step_1.tf"),
				ConfigVariables: config.Variables{
					"product_id":     config.StringVariable(productId),
					"environment_id": config.StringVariable(environmentId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(boolSettingResourceName, PredefinedVariations+".#", "2"),
					resource.TestCheckResourceAttrSet(boolSettingResourceName, PredefinedVariations+".0."+PredefinedVariationId),
					resource.TestCheckResourceAttr(boolSettingResourceName, PredefinedVariations+".0."+PredefinedVariationValue+"."+BoolValue, "false"),
					resource.TestCheckResourceAttr(boolSettingResourceName, PredefinedVariations+".0."+PredefinedVariationName, "Off name"),
					resource.TestCheckNoResourceAttr(boolSettingResourceName, PredefinedVariations+".0."+PredefinedVariationHint),
					resource.TestCheckResourceAttrSet(boolSettingResourceName, PredefinedVariations+".1."+PredefinedVariationId),
					resource.TestCheckResourceAttr(boolSettingResourceName, PredefinedVariations+".1."+PredefinedVariationValue+"."+BoolValue, "true"),
					resource.TestCheckNoResourceAttr(boolSettingResourceName, PredefinedVariations+".1."+PredefinedVariationName),
					resource.TestCheckResourceAttr(boolSettingResourceName, PredefinedVariations+".1."+PredefinedVariationHint, "On hint"),

					resource.TestCheckResourceAttr(stringSettingResourceName, PredefinedVariations+".#", "2"),
					resource.TestCheckResourceAttrSet(stringSettingResourceName, PredefinedVariations+".0."+PredefinedVariationId),
					resource.TestCheckResourceAttr(stringSettingResourceName, PredefinedVariations+".0."+PredefinedVariationValue+"."+StringValue, "Variation A"),
					resource.TestCheckResourceAttr(stringSettingResourceName, PredefinedVariations+".0."+PredefinedVariationName, "Variation A name"),
					resource.TestCheckNoResourceAttr(stringSettingResourceName, PredefinedVariations+".0."+PredefinedVariationHint),
					resource.TestCheckResourceAttrSet(stringSettingResourceName, PredefinedVariations+".1."+PredefinedVariationId),
					resource.TestCheckResourceAttr(stringSettingResourceName, PredefinedVariations+".1."+PredefinedVariationValue+"."+StringValue, "Variation B"),
					resource.TestCheckNoResourceAttr(stringSettingResourceName, PredefinedVariations+".1."+PredefinedVariationName),
					resource.TestCheckResourceAttr(stringSettingResourceName, PredefinedVariations+".1."+PredefinedVariationHint, "Variation B hint"),

					resource.TestCheckResourceAttr(intSettingResourceName, PredefinedVariations+".#", "2"),
					resource.TestCheckResourceAttrSet(intSettingResourceName, PredefinedVariations+".0."+PredefinedVariationId),
					resource.TestCheckResourceAttr(intSettingResourceName, PredefinedVariations+".0."+PredefinedVariationValue+"."+IntValue, "10"),
					resource.TestCheckResourceAttr(intSettingResourceName, PredefinedVariations+".0."+PredefinedVariationName, "10 name"),
					resource.TestCheckNoResourceAttr(intSettingResourceName, PredefinedVariations+".0."+PredefinedVariationHint),
					resource.TestCheckResourceAttrSet(intSettingResourceName, PredefinedVariations+".1."+PredefinedVariationId),
					resource.TestCheckResourceAttr(intSettingResourceName, PredefinedVariations+".1."+PredefinedVariationValue+"."+IntValue, "20"),
					resource.TestCheckNoResourceAttr(intSettingResourceName, PredefinedVariations+".1."+PredefinedVariationName),
					resource.TestCheckResourceAttr(intSettingResourceName, PredefinedVariations+".1."+PredefinedVariationHint, "20 hint"),

					resource.TestCheckResourceAttr(doubleSettingResourceName, PredefinedVariations+".#", "2"),
					resource.TestCheckResourceAttrSet(doubleSettingResourceName, PredefinedVariations+".0."+PredefinedVariationId),
					resource.TestCheckResourceAttr(doubleSettingResourceName, PredefinedVariations+".0."+PredefinedVariationValue+"."+DoubleValue, "10"),
					resource.TestCheckResourceAttr(doubleSettingResourceName, PredefinedVariations+".0."+PredefinedVariationName, "10.0 name"),
					resource.TestCheckNoResourceAttr(doubleSettingResourceName, PredefinedVariations+".0."+PredefinedVariationHint),
					resource.TestCheckResourceAttrSet(doubleSettingResourceName, PredefinedVariations+".1."+PredefinedVariationId),
					resource.TestCheckResourceAttr(doubleSettingResourceName, PredefinedVariations+".1."+PredefinedVariationValue+"."+DoubleValue, "20.2"),
					resource.TestCheckNoResourceAttr(doubleSettingResourceName, PredefinedVariations+".1."+PredefinedVariationName),
					resource.TestCheckResourceAttr(doubleSettingResourceName, PredefinedVariations+".1."+PredefinedVariationHint, "20.2 hint"),
				),
			},
			{
				ConfigFile: config.TestNameFile("step_2.tf"),
				ConfigVariables: config.Variables{
					"product_id":     config.StringVariable(productId),
					"environment_id": config.StringVariable(environmentId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
			{
				ConfigFile: config.TestNameFile("step_3.tf"),
				ConfigVariables: config.Variables{
					"product_id":     config.StringVariable(productId),
					"environment_id": config.StringVariable(environmentId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
			{
				ConfigFile: config.TestNameFile("step_3.tf"),
				ConfigVariables: config.Variables{
					"product_id":     config.StringVariable(productId),
					"environment_id": config.StringVariable(environmentId),
				},
				ResourceName:            boolSettingValueResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("step_3.tf"),
				ConfigVariables: config.Variables{
					"product_id":     config.StringVariable(productId),
					"environment_id": config.StringVariable(environmentId),
				},
				ResourceName:            stringSettingValueResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("step_3.tf"),
				ConfigVariables: config.Variables{
					"product_id":     config.StringVariable(productId),
					"environment_id": config.StringVariable(environmentId),
				},
				ResourceName:            intSettingValueResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
			{
				ConfigFile: config.TestNameFile("step_3.tf"),
				ConfigVariables: config.Variables{
					"product_id":     config.StringVariable(productId),
					"environment_id": config.StringVariable(environmentId),
				},
				ResourceName:            doubleSettingValueResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{InitOnly},
			},
		},
	})
}
