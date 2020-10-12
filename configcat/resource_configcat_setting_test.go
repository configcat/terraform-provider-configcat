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

func TestResourceSettingValidSettingTypes(t *testing.T) {
	const settingsResource = `
	resource "configcat_setting" "testBoolean" {
		config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
		key = "testBoolKey"
		name = "testBoolName"
		setting_type = "boolean"
	}

	resource "configcat_setting" "testString" {
		config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
		key = "testStringKey"
		name = "testStringName"
		setting_type = "string"
	}

	resource "configcat_setting" "testInt" {
		config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
		key = "testIntKey"
		name = "testIntName"
		setting_type = "int"
	}
	
	resource "configcat_setting" "testDouble" {
		config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
		key = "testDoubleKey"
		name = "testDoubleName"
		setting_type = "double"
	}
`
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: settingsResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting.testBoolean", "id"),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_KEY, "testBoolKey"),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_NAME, "testBoolName"),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_HINT, ""),
					resource.TestCheckResourceAttr("configcat_setting.testBoolean", SETTING_TYPE, "boolean"),

					resource.TestCheckResourceAttrSet("configcat_setting.testString", "id"),
					resource.TestCheckResourceAttr("configcat_setting.testString", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("configcat_setting.testString", SETTING_KEY, "testStringKey"),
					resource.TestCheckResourceAttr("configcat_setting.testString", SETTING_NAME, "testStringName"),
					resource.TestCheckResourceAttr("configcat_setting.testString", SETTING_HINT, ""),
					resource.TestCheckResourceAttr("configcat_setting.testString", SETTING_TYPE, "string"),

					resource.TestCheckResourceAttrSet("configcat_setting.testInt", "id"),
					resource.TestCheckResourceAttr("configcat_setting.testInt", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("configcat_setting.testInt", SETTING_KEY, "testIntKey"),
					resource.TestCheckResourceAttr("configcat_setting.testInt", SETTING_NAME, "testIntName"),
					resource.TestCheckResourceAttr("configcat_setting.testInt", SETTING_HINT, ""),
					resource.TestCheckResourceAttr("configcat_setting.testInt", SETTING_TYPE, "int"),

					resource.TestCheckResourceAttrSet("configcat_setting.testDouble", "id"),
					resource.TestCheckResourceAttr("configcat_setting.testDouble", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("configcat_setting.testDouble", SETTING_KEY, "testDoubleKey"),
					resource.TestCheckResourceAttr("configcat_setting.testDouble", SETTING_NAME, "testDoubleName"),
					resource.TestCheckResourceAttr("configcat_setting.testDouble", SETTING_HINT, ""),
					resource.TestCheckResourceAttr("configcat_setting.testDouble", SETTING_TYPE, "double"),
				),
			},
		},
	})
}
