package configcat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOrganizationValid(t *testing.T) {
	const dataSource = `
		data "configcat_organizations" "test" {
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_organizations.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_organizations.test", ORGANIZATIONS+".#", "1"),
				),
			},
		},
	})
}

func TestOrganizationValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_organizations" "test" {
			name_filter_regex = "ConfigCat"
		}
	`
	const organizationID = "08d86d63-26dc-4276-86d6-eae122660e51"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.configcat_organizations.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_organizations.test", ORGANIZATIONS+".#", "1"),
					resource.TestCheckResourceAttr("data.configcat_organizations.test", ORGANIZATIONS+".0."+ORGANIZATION_ID, organizationID),
					resource.TestCheckResourceAttr("data.configcat_organizations.test", ORGANIZATIONS+".0."+ORGANIZATION_NAME, "ConfigCat"),
				),
			},
		},
	})
}

func TestOrganizationNotFound(t *testing.T) {
	const dataSource = `
		data "configcat_organizations" "test" {
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
					resource.TestCheckResourceAttrSet("data.configcat_organizations.test", "id"),
					resource.TestCheckResourceAttr("data.configcat_organizations.test", ORGANIZATIONS+".#", "0"),
				),
			},
		},
	})
}
