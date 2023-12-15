package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testProductsDataSourceName = "data.configcat_products.test"

func TestProductValid(t *testing.T) {
	const dataSource = `
		data "configcat_products" "test" {
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testProductsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testProductsDataSourceName, PRODUCTS+".#", "1"),
				),
			},
		},
	})
}

func TestProductValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_products" "test" {
			name_filter_regex = "ConfigCat's product"
		}
	`
	const productID = "08d86d63-2721-4da6-8c06-584521d516bc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testProductsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testProductsDataSourceName, PRODUCTS+".#", "1"),
					resource.TestCheckResourceAttr(testProductsDataSourceName, PRODUCTS+".0."+PRODUCT_ID, productID),
					resource.TestCheckResourceAttr(testProductsDataSourceName, PRODUCTS+".0."+PRODUCT_NAME, "ConfigCat's product"),
					resource.TestCheckResourceAttr(testProductsDataSourceName, PRODUCTS+".0."+PRODUCT_DESCRIPTION, "ConfigCat's product description"),
					resource.TestCheckResourceAttr(testProductsDataSourceName, PRODUCTS+".0."+PRODUCT_ORDER, "0"),
				),
			},
		},
	})
}

func TestProductNotFound(t *testing.T) {
	const dataSource = `
		data "configcat_products" "test" {
			name_filter_regex = "notfound"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testProductsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testProductsDataSourceName, PRODUCTS+".#", "0"),
				),
			},
		},
	})
}
