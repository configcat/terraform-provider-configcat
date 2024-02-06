package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentsDataSource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "data.configcat_environments.test"

	const environment1ID = "08d8becf-d4d9-4c66-8b48-6ac74cd95fba"
	const environment2ID = "08d86d63-272c-4355-8027-4b52787bc1bd"
	const environment3ID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"
	const listAttribute = Environments

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
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "3"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+EnvironmentId, environment1ID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Name, "Mandatory"),
					resource.TestCheckNoResourceAttr(testResourceName, listAttribute+"[0]."+Description),
					resource.TestCheckNoResourceAttr(testResourceName, listAttribute+"[0]."+Color),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Order, "0"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[1]."+EnvironmentId, environment2ID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[1]."+Name, "Production"),
					resource.TestCheckNoResourceAttr(testResourceName, listAttribute+"[1]."+Description),
					resource.TestCheckNoResourceAttr(testResourceName, listAttribute+"[1]."+Color),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[1]."+Order, "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[2]."+EnvironmentId, environment3ID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[2]."+Name, "Test"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[2]."+Description, "Test Env Description"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[2]."+Color, "#5c6bc0"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[2]."+Order, "2"),
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
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+EnvironmentId, environment3ID),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Name, "Test"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Description, "Test Env Description"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Color, "#5c6bc0"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+"[0]."+Order, "2"),
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
