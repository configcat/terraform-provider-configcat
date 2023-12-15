package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceSettingValueMandatory(t *testing.T) {
	const environmentID = "08d8becf-d4d9-4c66-8b48-6ac74cd95fba"

	const settingResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueMandatory"
			name = "testMandatory"
			order = 30
		}
	`
	const settingValueResourceUpdatedWithoutMandatory = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueMandatory"
			name = "testMandatory"
			order = 30
		}

		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
		}
	`

	const settingValueResourceUpdatedWithMandatory = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testBool"
			name = "testBool"
			order = 30
		}
		
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			mandatory_notes = "mandatory note"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: settingResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting.testsetting", "id"),
				),
			},
			{
				Config:      settingValueResourceUpdatedWithoutMandatory,
				ExpectError: regexp.MustCompile(".*Reason required.*"),
			},
			{
				Config: settingValueResourceUpdatedWithMandatory,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TEST_RESOURCE, "id"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_TYPE, "boolean"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_VALUE, "true"),
					checkTrueValue,
				),
			},
		},
	})
}
