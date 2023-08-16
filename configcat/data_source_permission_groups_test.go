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
	const permissionGroupID = "219"

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
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NAME, "Administrators"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_MEMBERS, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_CONFIG, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_TAG_SETTING, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_SETTING, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_TAG, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_SDKKEY, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_ROTATE_SDKKEY, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_SEGMENTS, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, "true"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ACCESSTYPE, "full"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE, "full"),
					resource.TestCheckResourceAttr(test_resource_name, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ENVIRONMENT_ACCESSES+".#", "0"),
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
