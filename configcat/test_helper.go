package configcat

import (
	"fmt"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func checkTest1Value(s *terraform.State) error {
	value, err := getSettingValue("string", "test1")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkTest2Value(s *terraform.State) error {
	value, err := getSettingValue("string", "test2")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
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
