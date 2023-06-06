package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceProductFlow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_organizations" "organizations" {
					}
					resource "configcat_product" "test" {
						organization_id = data.configcat_organizations.organizations.organizations.0.organization_id
						name = "TestResourceProductFlow"
						description = "testDescription"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_product.test", "id"),
					resource.TestCheckResourceAttr("configcat_product.test", PRODUCT_NAME, "TestResourceProductFlow"),
					resource.TestCheckResourceAttr("configcat_product.test", PRODUCT_DESCRIPTION, "testDescription"),
				),
			},
			{
				Config: `
					data "configcat_organizations" "organizations" {
					}
					resource "configcat_product" "test" {
						organization_id = data.configcat_organizations.organizations.organizations.0.organization_id
						name = "TestResourceProductFlow2"
						description = "testDescription2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("configcat_product.test", "id"),
					resource.TestCheckResourceAttr("configcat_product.test", PRODUCT_NAME, "TestResourceProductFlow2"),
					resource.TestCheckResourceAttr("configcat_product.test", PRODUCT_DESCRIPTION, "testDescription2"),
				),
			},
			{
				ResourceName:      "configcat_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
