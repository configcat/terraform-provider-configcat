package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testPermissionGroupsDataSourceName = "data.configcat_permission_groups.test_permission_group"

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
					resource.TestCheckResourceAttrSet(testPermissionGroupsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".#", "2"),
				),
			},
		},
	})
}

func TestAdministratorsPermissionGroupValidFilter(t *testing.T) {
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
					resource.TestCheckResourceAttrSet(testPermissionGroupsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".#", "1"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ID, permissionGroupID),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NAME, "Administrators"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_MEMBERS, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_CONFIG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_TAG_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_TAG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_SDKKEY, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_ROTATE_SDKKEY, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_SEGMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ACCESSTYPE, "full"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE, "full"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ENVIRONMENT_ACCESS+".#", "0"),
				),
			},
		},
	})
}

func TestOnlyTestEnvironmentPermissionGroupValidFilter(t *testing.T) {
	const dataSource = `
		data "configcat_permission_groups" "test_permission_group" {
			name_filter_regex = "Only test environment"
			product_id = "08d86d63-2721-4da6-8c06-584521d516bc"
		}
	`
	const permissionGroupID = "29859"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPermissionGroupsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".#", "1"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ID, permissionGroupID),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NAME, "Only test environment"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_MEMBERS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_CONFIG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_TAG_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_TAG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_SDKKEY, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_ROTATE_SDKKEY, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_DELETE_SEGMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ACCESSTYPE, "custom"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE, "readOnly"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ENVIRONMENT_ACCESS+".#", "1"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ENVIRONMENT_ACCESS+".0."+PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ID, "08d86d63-2726-47cd-8bfc-59608ecb91e2"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".0."+PERMISSION_GROUP_ENVIRONMENT_ACCESS+".0."+PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ACCESS_TYPE, "full"),
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
					resource.TestCheckResourceAttrSet(testPermissionGroupsDataSourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupsDataSourceName, PERMISSION_GROUPS+".#", "0"),
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
