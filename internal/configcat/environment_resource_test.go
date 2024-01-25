package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentResource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "configcat_environment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable(productId),
					"name":       config.StringVariable("Resource name"),
					"order":      config.IntegerVariable(1),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name"),
					resource.TestCheckResourceAttr(testResourceName, Order, "1"),
					resource.TestCheckNoResourceAttr(testResourceName, Description),
					resource.TestCheckNoResourceAttr(testResourceName, Color),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":  config.StringVariable(productId),
					"name":        config.StringVariable("Resource name updated"),
					"description": config.StringVariable("Resource description"),
					"color":       config.StringVariable("#5c6bc0"),
					"order":       config.IntegerVariable(10),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name updated"),
					resource.TestCheckResourceAttr(testResourceName, Order, "10"),
					resource.TestCheckResourceAttr(testResourceName, Description, "Resource description"),
					resource.TestCheckResourceAttr(testResourceName, Color, "#5c6bc0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":  config.StringVariable(productId),
					"name":        config.StringVariable("Resource name updated"),
					"description": config.StringVariable("Resource description"),
					"color":       config.StringVariable("#5c6bc0"),
					"order":       config.IntegerVariable(10),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}