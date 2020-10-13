package configcat

import (
	"fmt"
	"strconv"
	"testing"

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
			setting_type = "boolean"
			value = "true"
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			setting_type = "boolean"
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
			setting_type = "boolean"
			value = "true"
			freeze_after_init = false
		}
	`
	const settingValueResourceUpdated = `
		resource "configcat_setting_value" "test" {
			environment_id = "` + environmentID + `"
			setting_id = "` + settingID + `"
			setting_type = "boolean"
			value = "false"
			freeze_after_init = false
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
	c := testAccProvider.Meta().(*Client)
	rs := s.RootModule().Resources["configcat_setting_value.test"]
	environmentID := rs.Primary.Attributes[ENVIRONMENT_ID]
	settingID, err := strconv.ParseInt(rs.Primary.Attributes[SETTING_ID], 10, 32)
	if err != nil {
		return err
	}

	settingValue, err := c.GetSettingValueSimple(environmentID, int32(settingID))
	if err != nil {
		return err
	}

	if *settingValue.Value != *value {
		return fmt.Errorf("%v != %v", *settingValue.Value, *value)
	}
	return nil
}
