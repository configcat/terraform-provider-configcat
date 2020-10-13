package configcat

import (
	"regexp"
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

func TestResourceSettingBoolean(t *testing.T) {
	testResourceSettingForSettingType(t, "boolean")
}

func TestResourceSettingString(t *testing.T) {
	testResourceSettingForSettingType(t, "string")
}

func TestResourceSettingInt(t *testing.T) {
	testResourceSettingForSettingType(t, "int")
}

func TestResourceSettingDouble(t *testing.T) {
	testResourceSettingForSettingType(t, "double")
}

func testResourceSettingForSettingType(t *testing.T, settingType string) {
	var settingResource = `
	resource "configcat_setting" "testBoolean" {
		config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
		key = "test` + settingType + `"
		name = "test"
		setting_type = "` + settingType + `"
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
					resource.TestCheckResourceAttrSet("configcat_setting.testBoolean", "id"),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_KEY, "test"+settingType),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_NAME, "test"),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_HINT, ""),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_TYPE, settingType),
				),
			},
			resource.TestStep{
				ResourceName:      "configcat_setting.testBoolean",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceSettingInvalidSettingType(t *testing.T) {
	const settingResource = `
		resource "configcat_setting" "test2" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testKey2"
			name = "testName"
			setting_type = "invalid"
		}
	`
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      settingResource,
				ExpectError: regexp.MustCompile(`setting_type parse failed: invalid. Valid values: boolean/string/int/double`),
			},
		},
	})
}

func TestResourceSettingDuplicatedKey(t *testing.T) {
	const settingResource = `
		resource "configcat_setting" "test3" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "isAwesomeFeatureEnabled"
			name = "testName"
		}
	`
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      settingResource,
				ExpectError: regexp.MustCompile(`.*This key is already in use\. Please, choose another.*`),
			},
		},
	})
}
