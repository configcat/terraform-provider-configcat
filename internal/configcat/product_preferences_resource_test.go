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
		},
	})
}
