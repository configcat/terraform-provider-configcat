package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagsDataSource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const tagID = "46"
	const testResourceName = "data.configcat_tags.test"

	const listAttribute = Tags

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable(productId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+TagId, tagID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Name, "Test"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Color, "panther"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("Test"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+TagId, tagID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Name, "Test"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Color, "panther"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("notfound"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable("invalidguid"),
				},
				ExpectError: regexp.MustCompile(`Attribute product_id value must be a valid GUID, got: invalidguid`),
			},
		},
	})
}
