package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSettingsValid(t *testing.T) {
	const dataSource = `
		data "configcat_settings" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "isAwesomeFeatureEnabled"
		}
	`
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.configcat_settings.test", "id", "67639"),
					resource.TestCheckResourceAttr("data.configcat_settings.test", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("data.configcat_settings.test", SETTING_ID, "67639"),
					resource.TestCheckResourceAttr("data.configcat_settings.test", SETTING_KEY, "isAwesomeFeatureEnabled"),
					resource.TestCheckResourceAttr("data.configcat_settings.test", SETTING_NAME, "My awesome feature flag"),
					resource.TestCheckResourceAttr("data.configcat_settings.test", SETTING_HINT, "This is the hint for my awesome feature flag"),
					resource.TestCheckResourceAttr("data.configcat_settings.test", SETTING_TYPE, "boolean"),
				),
			},
		},
	})
}
