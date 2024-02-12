package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigV2Resource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "configcat_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":         config.StringVariable(productId),
					"name":               config.StringVariable("Resource name"),
					"order":              config.IntegerVariable(1),
					"evaluation_version": config.StringVariable("v2"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name"),
					resource.TestCheckResourceAttr(testResourceName, Order, "1"),
					resource.TestCheckResourceAttr(testResourceName, Description, ""),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":         config.StringVariable(productId),
					"name":               config.StringVariable("Resource name updated"),
					"description":        config.StringVariable("Resource description"),
					"order":              config.IntegerVariable(10),
					"evaluation_version": config.StringVariable("v2"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name updated"),
					resource.TestCheckResourceAttr(testResourceName, Order, "10"),
					resource.TestCheckResourceAttr(testResourceName, Description, "Resource description"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":         config.StringVariable(productId),
					"name":               config.StringVariable("Resource name updated"),
					"description":        config.StringVariable("Resource description"),
					"order":              config.IntegerVariable(10),
					"evaluation_version": config.StringVariable("v2"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":         config.StringVariable(productId),
					"name":               config.StringVariable("Resource name updated"),
					"description":        config.StringVariable("Resource description"),
					"order":              config.IntegerVariable(10),
					"evaluation_version": config.StringVariable("v1"),
				},
				ExpectError: regexp.MustCompile("evaluation_version cannot be changed."),
			},
		},
	})
}
