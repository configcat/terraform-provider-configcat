package configcat

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestSetting(t *testing.T) {
	const setting1 = `
		resource "configcat_setting" "setting1" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "setting1Key"
			name = "setting1Name"
			hint = "setting1Hint"
		}
	`

	const setting1Updated = `
		resource "configcat_setting" "setting1" {
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
			key = "setting1Key"
			name = "setting1NameUpdated"
			hint = "setting1HintUpdated"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: setting1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting.setting1", "key", "setting1Key"),
					resource.TestCheckResourceAttr("configcat_setting.setting1", "name", "setting1Name"),
					resource.TestCheckResourceAttr("configcat_setting.setting1", "hint", "setting1Hint"),
				),
				ExpectNonEmptyPlan: true,
			},
			resource.TestStep{
				Config: setting1Updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("configcat_setting.setting1", "key", "setting1Key"),
					resource.TestCheckResourceAttr("configcat_setting.setting1", "name", "setting1NameUpdated"),
					resource.TestCheckResourceAttr("configcat_setting.setting1", "hint", "setting1HintUpdated"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
