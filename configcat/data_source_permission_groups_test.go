package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const test_resource_name = "data.configcat_permission_groups.test_permission_group"

func TestPermissionGroupValid(t *testing.T) {
	const dataSource = `
		data "configcat_permission_groups" "test_permission_group" {
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(test_resource_name, "id"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".#", "1"),
				),
			},
		},
	})
}

func TestPermissionGroupValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_permission_groups" "test_permission_group" {
			name_filter_regex = "Administrators"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const permissionGroupID = "08d86d63-2726-47cd-8bfc-59608ecb91e2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(test_resource_name, "id"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".#", "1"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ID, permissionGroupID),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NAME, "Test"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NAME, "Test"),
					// TODO
				),
			},
		},
	})
}

func TestPermissionGroupNotFoundFilter(t *testing.T) {
	const dataSource = `
		data "configcat_permission_groups" "test_permission_group" {
			name_filter_regex = "invalid"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(test_resource_name, "id"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".#", "0"),
				),
			},
		},
	})
}

func TestPermissionGroupInvalidGuid(t *testing.T) {
	const dataSource = `
		data "configcat_permission_groups" "test_permission_group" {
			name_filter_regex = "notfound"
			product_id = "invalidGuid"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      dataSource,
				ExpectError: regexp.MustCompile(`"product_id": invalid GUID`),
			},
		},
	})
}
