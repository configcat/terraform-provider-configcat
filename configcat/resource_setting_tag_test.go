package configcat

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceSettingTag(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					data "configcat_configs" "configs" {
						product_id = data.configcat_products.products.products.0.product_id
					}
					data "configcat_settings" "settings" {
						config_id = data.configcat_configs.configs.configs.0.config_id
					}
					data "configcat_tags" "tags" {
						product_id = data.configcat_products.products.products.0.product_id
					}
					resource "configcat_setting_tag" "setting_tag" {
						setting_id = data.configcat_settings.settings.settings.0.setting_id
						tag_id = data.configcat_tags.tags.tags.0.tag_id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_setting_tag.setting_tag", "id"),
					resource.TestCheckResourceAttrSet("configcat_setting_tag.setting_tag", "setting_id"),
					resource.TestCheckResourceAttrSet("configcat_setting_tag.setting_tag", "tag_id"),
					checkTag,
				),
			},
			{
				ResourceName:      "configcat_setting_tag.setting_tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkTag(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)
	rs := s.RootModule().Resources["configcat_setting_tag.setting_tag"]
	settingID, err := strconv.ParseInt(rs.Primary.Attributes[SETTING_ID], 10, 32)
	if err != nil {
		return err
	}
	tagID, err := strconv.ParseInt(rs.Primary.Attributes[TAG_ID], 10, 64)
	if err != nil {
		return err
	}

	setting, err := c.GetSetting(int32(settingID))
	if err != nil {
		return err
	}

	for _, tag := range setting.Tags {
		if *tag.TagId == tagID {
			return nil
		}
	}

	return fmt.Errorf("Could not find setting tag")
}
