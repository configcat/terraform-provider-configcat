package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductResource(t *testing.T) {
	const organizationId = "08d86d63-26dc-4276-86d6-eae122660e51"
	const testResourceName = "configcat_product.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"organization_id": config.StringVariable(organizationId),
					"name":            config.StringVariable("Resource name"),
					"order":           config.IntegerVariable(1),
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
					"organization_id": config.StringVariable(organizationId),
					"name":            config.StringVariable("Resource name updated"),
					"description":     config.StringVariable("Resource description"),
					"order":           config.IntegerVariable(10),
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
					"organization_id": config.StringVariable(organizationId),
					"name":            config.StringVariable("Resource name updated"),
					"description":     config.StringVariable("Resource description"),
					"order":           config.IntegerVariable(10),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
