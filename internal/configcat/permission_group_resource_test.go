package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermissionGroupResource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "configcat_permission_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable(productId),
					"name":       config.StringVariable("Resource name"),
					"accesstype": config.StringVariable("full"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, ID),
					resource.TestCheckResourceAttr(testResourceName, Name, "Resource name"),
					resource.TestCheckResourceAttr(testResourceName, PermissionGroupAccessType, "full"),
					resource.TestCheckResourceAttr(testResourceName, PermissionGroupNewEnvironmentAccessType, "none"),
					resource.TestCheckResourceAttr(testResourceName, PermissionGroupEnvironmentAccess+".%", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable(productId),
					"name":       config.StringVariable("Resource name"),
					"accesstype": config.StringVariable("full"),
				},
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
