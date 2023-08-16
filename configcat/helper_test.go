package configcat

import (
	"fmt"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const TEST_RESOURCE = "configcat_setting_value.test"

func checkTest1Value(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_STRING.Ptr(), "test1")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkTest2Value(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_STRING.Ptr(), "test2")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkTest1IntValue(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_INT.Ptr(), "1")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkTest2IntValue(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_INT.Ptr(), "2")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkTest1FloatValue(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_DOUBLE.Ptr(), "1.1")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkTest2FloatValue(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_DOUBLE.Ptr(), "2.1")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkTrueValue(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_BOOLEAN.Ptr(), "true")
	if err != nil {
		return err
	}
	return checkValue(s, &value)
}

func checkFalseValue(s *terraform.State) error {
	value, err := getSettingValue(sw.SETTINGTYPE_BOOLEAN.Ptr(), "false")
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

	if fmt.Sprintf("%v", settingValue.Value) != fmt.Sprintf("%v", *value) {
		return fmt.Errorf("%v != %v", settingValue.Value, *value)
	}

	return nil
}

func getValue(s *terraform.State) (*sw.SettingValueModel, error) {
	c := testAccProvider.Meta().(*Client)
	rs := s.RootModule().Resources["configcat_setting_value.test"]
	environmentID := rs.Primary.Attributes[ENVIRONMENT_ID]
	settingID, err := strconv.ParseInt(rs.Primary.Attributes[SETTING_ID], 10, 32)
	if err != nil {
		return nil, err
	}
	settingValue, err := c.GetSettingValue(environmentID, int32(settingID))
	if err != nil {
		return nil, err
	}
	return settingValue, nil
}
