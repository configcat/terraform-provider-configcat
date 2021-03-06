package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceConfigFlow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_config" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_config.test", "id"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_NAME, "testName"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_config" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_config.test", "id"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_NAME, "testName2"),
				),
			},
			{
				ResourceName:      "configcat_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
