package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductsDataSource(t *testing.T) {
	const testResourceName = "data.configcat_products.test"

	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const listAttribute = Products

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+ProductId, productId),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "ConfigCat's product"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Description, "ConfigCat's product description"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Order, "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"name_filter_regex": config.StringVariable("ConfigCat's product"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+ProductId, productId),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "ConfigCat's product"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Description, "ConfigCat's product description"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Order, "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"name_filter_regex": config.StringVariable("notfound"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "0"),
				),
			},
		},
	})
}
