package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestConfigValid(t *testing.T) {
	const dataSource = `
		data "configcat_config" "test" {
			name = "Main Config"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const productID = "08d86d63-2721-4da6-8c06-584521d516bc"
	const configID = "08d86d63-2731-4b8b-823a-56ddda9da038"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.configcat_config.test", "id", configID),
					resource.TestCheckResourceAttr("data.configcat_config.test", PRODUCT_ID, productID),
					resource.TestCheckResourceAttr("data.configcat_config.test", CONFIG_ID, configID),
					resource.TestCheckResourceAttr("data.configcat_config.test", CONFIG_NAME, "Main Config"),
				),
			},
		},
	})
}

func TestConfigInvalid(t *testing.T) {
	const dataSource = `
		data "configcat_config" "test" {
			name = "notfound"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      dataSource,
				ExpectError: regexp.MustCompile("could not find Config. product_id: 08d86d63-2721-4da6-8c06-584521d516bc name: notfound"),
			},
		},
	})
}
