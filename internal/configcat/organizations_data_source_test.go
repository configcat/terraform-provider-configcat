package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationsDataSource(t *testing.T) {
	const testResourceName = "data.configcat_organizations.test"

	const organizationId = "08d86d63-26dc-4276-86d6-eae122660e51"
	const listAttribute = Organizations

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+OrganizationId, organizationId),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "ConfigCat"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"name_filter_regex": config.StringVariable("ConfigCat"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+OrganizationId, organizationId),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "ConfigCat"),
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
