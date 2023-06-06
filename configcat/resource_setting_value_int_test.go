package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceSettingValueIntFreeze(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntFreeze"
			name = "testInt"
			setting_type = "int"
		}

		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntFreeze"
			name = "testInt"
			setting_type = "int"
		}
		
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "2"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TEST_RESOURCE, "id"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_TYPE, "int"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_VALUE, "1"),
					checkTest1IntValue,
				),
			},
			{
				Config: settingValueResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TEST_RESOURCE, "id"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_TYPE, "int"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_VALUE, "2"),
					checkTest1IntValue,
				),
			},
		},
	})
}

func TestResourceSettingValueIntNoFreeze(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntNoFreeze"
			name = "testInt"
			setting_type = "int"
		}

		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
			init_only = false
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntNoFreeze"
			name = "testInt"
			setting_type = "int"
		}

		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "2"
			init_only = false
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TEST_RESOURCE, "id"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_TYPE, "int"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_VALUE, "1"),
					checkTest1IntValue,
				),
			},
			{
				Config: settingValueResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TEST_RESOURCE, "id"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_TYPE, "int"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, SETTING_VALUE, "2"),
					checkTest2IntValue,
				),
			},
			{
				ResourceName:      TEST_RESOURCE,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceSettingValueIntRules(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntRules"
			name = "testInt"
			setting_type = "int"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
			init_only = false
		}
	`
	const settingValueResourceRule1 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntRules"
			name = "testInt"
			setting_type = "int"
		}	
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
			init_only = false
			rollout_rules {
				comparison_attribute = "email"
				comparator = "contains"
				comparison_value = "@configcat"
				value = "1"
			}
		}
	`

	const settingValueResourceRule2 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntRules"
			name = "testInt"
			setting_type = "int"
		}	
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
			init_only = false
			rollout_rules {
				comparison_attribute = "email"
				comparator = "contains"
				comparison_value = "@configcat"
				value = "1"
			}
			rollout_rules {
				comparison_attribute = "custom"
				comparator = "isOneOf"
				comparison_value = "red"
				value = "2"
			}
		}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".#", "0"),
				),
			},
			{
				Config: settingValueResourceRule1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".#", "1"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, "email"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARATOR, "contains"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_VALUE, "@configcat"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_VALUE, "1"),
				),
			},
			{
				Config: settingValueResourceRule2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".#", "2"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, "email"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARATOR, "contains"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_VALUE, "@configcat"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".0."+ROLLOUT_RULE_VALUE, "1"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, "custom"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARATOR, "isOneOf"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARISON_VALUE, "red"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_RULES+".1."+ROLLOUT_RULE_VALUE, "2"),
				),
			},
		},
	})
}

func TestResourceSettingValueIntPercentageItems(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntPercentageItems"
			name = "testInt"
			setting_type = "int"
		} 
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
			init_only = false
		}
	`
	const settingValueResourceItem1 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntPercentageItems"
			name = "testInt"
			setting_type = "int"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
			init_only = false
			percentage_items {
				percentage = 20
				value = "1"
			}
			percentage_items {
				percentage = 80
				value = "0"
			}
		}
	`

	const settingValueResourceItem2 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "TestResourceSettingValueIntPercentageItems"
			name = "testInt"
			setting_type = "int"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "1"
			init_only = false
			percentage_items {
				percentage = 50
				value = "1"
			}
			percentage_items {
				percentage = 50
				value = "0"
			}
		}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".#", "0"),
				),
			},
			{
				Config: settingValueResourceItem1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".#", "2"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "20"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "1"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "80"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "0"),
				),
			},
			{
				Config: settingValueResourceItem2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".#", "2"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "50"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "1"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "50"),
					resource.TestCheckResourceAttr(TEST_RESOURCE, ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "0"),
				),
			},
		},
	})
}
