package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceTagFlow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_tag" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestResourceTagFlow"
						color = "panther"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_tag.test", "id"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_NAME, "TestResourceTagFlow"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_COLOR, "panther"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_tag" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestResourceTagFlow2"
						color = "koala"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_tag.test", "id"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_NAME, "TestResourceTagFlow2"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_COLOR, "koala"),
				),
			},
			{
				ResourceName:      "configcat_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
