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
						name = "TestResourceConfigFlow"
						description = "testDescription"
						order = 10
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_config.test", "id"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_NAME, "TestResourceConfigFlow"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_DESCRIPTION, "testDescription"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_ORDER, "10"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_config" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName2"
						description = "testDescription2"
						order = 11
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_config.test", "id"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_NAME, "testName2"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_DESCRIPTION, "testDescription2"),
					resource.TestCheckResourceAttr("configcat_config.test", CONFIG_ORDER, "11"),
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
