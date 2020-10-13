package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceSettingValueExistingFreeze(t *testing.T) {
	const settingID = "67639"
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			value = "true"
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			value = "false"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting_value.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "boolean"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "true"),
				),
			},
			resource.TestStep{
				Config: settingValueResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting_value.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "boolean"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "true"),
				),
			},
		},
	})
}
