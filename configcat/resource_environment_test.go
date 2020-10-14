package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceEnvironmentFlow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_environment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "testName"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_environment" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_environment.test", "id"),
					resource.TestCheckResourceAttr("configcat_environment.test", ENVIRONMENT_NAME, "testName2"),
				),
			},
			{
				ResourceName:      "configcat_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
