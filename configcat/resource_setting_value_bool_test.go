package configcat

import (
	"fmt"
	"strconv"
	"testing"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					checkTrueValue,
				),
			},
			resource.TestStep{
				Config: settingValueResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting_value.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "boolean"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "false"),
					checkTrueValue,
				),
			},
		},
	})
}

func TestResourceSettingValueExistingNoFreeze(t *testing.T) {
	const settingID = "67639"
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			value = "true"
			init_only = false
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			value = "false"
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
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "boolean"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "true"),
					checkTrueValue,
				),
			},
			resource.TestStep{
				Config: settingValueResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting_value.test", "id"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_TYPE, "boolean"),
					resource.TestCheckResourceAttr("configcat_setting_value.test", SETTING_VALUE, "false"),
					checkFalseValue,
				),
			},
		},
	})
}

func TestResourceSettingValueRules(t *testing.T) {
	const settingID = "67639"
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			value = "true"
			init_only = false
		}
	`
	const settingValueResourceRule1 = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
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
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
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

func TestResourceSettingValuePercentageItems(t *testing.T) {
	const settingID = "67639"
	const environmentID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	const settingValueResource = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			value = "true"
			init_only = false
		}
	`
	const settingValueResourceItem1 = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
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
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
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

func checkTrueValue(s *terraform.State) error {
	value, err := getSettingValue("boolean", "true")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkFalseValue(s *terraform.State) error {
	value, err := getSettingValue("boolean", "false")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkValue(s *terraform.State, value *interface{}) error {
	settingValue, err := getValue(s)
	if err != nil {
		return err
	}

	if *settingValue.Value != *value {
		return fmt.Errorf("%v != %v", *settingValue.Value, *value)
	}

	return nil
}

func getValue(s *terraform.State) (*sw.SettingValueSimpleModel, error) {
	c := testAccProvider.Meta().(*Client)
	rs := s.RootModule().Resources["configcat_setting_value.test"]
	environmentID := rs.Primary.Attributes[ENVIRONMENT_ID]
	settingID, err := strconv.ParseInt(rs.Primary.Attributes[SETTING_ID], 10, 32)
	if err != nil {
		return nil, err
	}
	settingValue, err := c.GetSettingValueSimple(environmentID, int32(settingID))
	if err != nil {
		return nil, err
	}
	return &settingValue, nil
}
