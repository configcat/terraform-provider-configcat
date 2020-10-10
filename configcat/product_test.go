package configcat

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	productId = "08d86d63-2721-4da6-8c06-584521d516bc"
)

func TestProduct(t *testing.T) {
	const productValid = `
		data "configcat_product" "test" {
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: productValid,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.configcat_product.test", "id", productId),
					resource.TestCheckResourceAttr("data.configcat_product.test", "product_id", productId),
					resource.TestCheckResourceAttr("data.configcat_product.test", "name", "Configcat's product"),
				),
			},
		},
	})
}
