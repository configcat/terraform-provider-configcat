package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testSettingsDataSourceName = "data.configcat_settings.test"

func TestSettingValid(t *testing.T) {
	const dataSource = `
		data "configcat_settings" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testSettingsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".#", "1"),
				),
			},
		},
	})
}

func TestSettingValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_settings" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key_filter_regex = "isAwesomeFeatureEnabled"
		}
	`
	const settingID = "67639"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testSettingsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".#", "1"),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".0."+SETTING_ID, settingID),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".0."+SETTING_KEY, "isAwesomeFeatureEnabled"),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".0."+SETTING_NAME, "My awesome feature flag"),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".0."+SETTING_HINT, "This is the hint for my awesome feature flag"),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".0."+SETTING_TYPE, "boolean"),
				),
			},
		},
	})
}

func TestSettingInvalid(t *testing.T) {
	const dataSource = `
		data "configcat_settings" "test" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key_filter_regex = "notfound"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testSettingsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testSettingsDataSourceName, SETTINGS+".#", "0"),
				),
			},
		},
	})
}

func TestSettingInvalidGuid(t *testing.T) {
	const dataSource = `
		data "configcat_settings" "test" {
			config_id = "invalidGuid"
			key_filter_regex = "notfound"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      dataSource,
				ExpectError: regexp.MustCompile(`"config_id": invalid GUID`),
			},
		},
	})
}
