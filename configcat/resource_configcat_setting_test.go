package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceSettingValid(t *testing.T) {
	const settingResource = `
		resource "configcat_setting" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testKey"
			name = "testName"
		}
	`
	const settingResourceUpdated = `
		resource "configcat_setting" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testKey"
			name = "testNameUpdated"
			hint = "testHintUpdated"
		}
	`
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: settingResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting.test", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_KEY, "testKey"),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_NAME, "testName"),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_HINT, ""),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_TYPE, "boolean"),
				),
			},
			resource.TestStep{
				Config: settingResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting.test", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_KEY, "testKey"),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_NAME, "testNameUpdated"),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_HINT, "testHintUpdated"),
					resource.TestCheckResourceAttr("configcat_setting.test", SETTING_TYPE, "boolean"),
				),
			},
		},
	})
}
