package configcat

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	configId = "08d86d63-2731-4b8b-823a-56ddda9da038"
)

func TestConfig(t *testing.T) {
	const productValid = `
		data "configcat_config" "test" {
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
			config_id = "08d86d63-2731-4b8b-823a-56ddda9da038"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: productValid,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.configcat_config.test", "id", configId),
					resource.TestCheckResourceAttr("data.configcat_config.test", "product_id", productId),
					resource.TestCheckResourceAttr("data.configcat_config.test", "config_id", configId),
					resource.TestCheckResourceAttr("data.configcat_config.test", "name", "Main Config"),
				),
			},
		},
	})
}
