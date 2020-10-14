package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceSettingValueStringFreeze(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}

		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "test1"
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}
		
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "test2"
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
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "string"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "test1"),
					checkTest1Value,
				),
			},
			resource.TestStep{
				Config: settingValueResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting_value.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "string"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "test2"),
					checkTest1Value,
				),
			},
		},
	})
}

func TestResourceSettingValueStringNoFreeze(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}

		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "test1"
			init_only = false
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}

		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "test2"
			init_only = false
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
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "string"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "test1"),
					checkTest1Value,
				),
			},
			resource.TestStep{
				Config: settingValueResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting_value.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "string"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "test2"),
					checkTest2Value,
				),
			},
			resource.TestStep{
				ResourceName:      "configcat_setting_value.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceSettingValueStringRules(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
		}
	`
	const settingValueResourceRule1 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}	
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
			rollout_rules {
				comparison_attribute = "email"
				comparator = "contains"
				comparison_value = "@configcat"
				value = "true"
			}
		}
	`

	const settingValueResourceRule2 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}	
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
			rollout_rules {
				comparison_attribute = "email"
				comparator = "contains"
				comparison_value = "@configcat"
				value = "true"
			}
			rollout_rules {
				comparison_attribute = "custom"
				comparator = "isOneOf"
				comparison_value = "red"
				value = "false"
			}
		}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".#", "0"),
				),
			},
			resource.TestStep{
				Config: settingValueResourceRule1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".#", "1"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, "email"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARATOR, "contains"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_VALUE, "@configcat"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_VALUE, "true"),
				),
			},
			resource.TestStep{
				Config: settingValueResourceRule2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".#", "2"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, "email"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARATOR, "contains"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_COMPARISON_VALUE, "@configcat"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".0."+ROLLOUT_RULE_VALUE, "true"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARISON_ATTRIBUTE, "custom"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARATOR, "isOneOf"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".1."+ROLLOUT_RULE_COMPARISON_VALUE, "red"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_RULES+".1."+ROLLOUT_RULE_VALUE, "false"),
				),
			},
		},
	})
}

func TestResourceSettingValueStringPercentageItems(t *testing.T) {
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		} 
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
		}
	`
	const settingValueResourceItem1 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
			percentage_items {
				percentage = 20
				value = "true"
			}
			percentage_items {
				percentage = 80
				value = "false"
			}
		}
	`

	const settingValueResourceItem2 = `
		resource "configcat_setting" "testsetting" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "testString"
			name = "testString"
			setting_type = "string"
		}
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = configcat_setting.testsetting.id
			value = "true"
			init_only = false
			percentage_items {
				percentage = 50
				value = "true"
			}
			percentage_items {
				percentage = 50
				value = "false"
			}
		}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: settingValueResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".#", "0"),
				),
			},
			resource.TestStep{
				Config: settingValueResourceItem1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".#", "2"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "20"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "true"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "80"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "false"),
				),
			},
			resource.TestStep{
				Config: settingValueResourceItem2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".#", "2"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "50"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".0."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "true"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE, "50"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", ROLLOUT_PERCENTAGE_ITEMS+".1."+ROLLOUT_PERCENTAGE_ITEM_VALUE, "false"),
				),
			},
		},
	})
}
