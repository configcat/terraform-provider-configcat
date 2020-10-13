package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProductValid(t *testing.T) {
	const dataSource = `
		data "configcat_product" "test" {
			name_filter_regex = "Configcat's product"
		}
	`
	const productID = "08d86d63-2721-4da6-8c06-584521d516bc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_product.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_product.test.products", "#", "1"),
					resource.TestCheckResourceAttr("data.configcat_product.test.products.0", PRODUCT_ID, productID),
					resource.TestCheckResourceAttr("data.configcat_product.test.products.0", PRODUCT_NAME, "Configcat's product"),
				),
			},
		},
	})
}

func TestProductNotFound(t *testing.T) {
	const dataSource = `
		data "configcat_product" "test" {
			name_filter_regex = "notfound"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_product.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_product.test.products", "#", "0"),
				),
			},
		},
	})
}
