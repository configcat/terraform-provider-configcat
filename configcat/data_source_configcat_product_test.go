package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProductValid(t *testing.T) {
	const dataSource = `
		data "configcat_product" "test" {
			name = "Configcat's product"
		}
	`
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.configcat_product.test", "id", productId),
					resource.TestCheckResourceAttr("data.configcat_product.test", "product_id", productId),
					resource.TestCheckResourceAttr("data.configcat_product.test", "name", "Configcat's product"),
				),
			},
		},
	})
}

func TestProductInvalid(t *testing.T) {
	const dataSource = `
		data "configcat_product" "test" {
			name = "notfound"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      dataSource,
				ExpectError: regexp.MustCompile("could not find Product with name: notfound"),
			},
		},
	})
}
