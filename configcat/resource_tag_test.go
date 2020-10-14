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
						name = "testName"
						color = "testColor"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_tag.test", "id"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_NAME, "testName"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_COLOR, "testColor"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_tag" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "testName2"
						color = "testColor2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_tag.test", "id"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_NAME, "testName2"),
					resource.TestCheckResourceAttr("configcat_tag.test", TAG_COLOR, "testColor2"),
				),
			},
			/*	{
				ResourceName:      "configcat_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},*/
		},
	})
}
