package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testPermissionGroupResourceName = "configcat_permission_group.test"

func TestResourcePermissionGroupFlow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_permission_group" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestPermissionGroup"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPermissionGroupResourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_NAME, "TestPermissionGroup"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_MEMBERS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_CONFIG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_TAG_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_TAG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_SDKKEY, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_ROTATE_SDKKEY, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SEGMENTS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_ACCESSTYPE, "custom"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE, "none"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_ENVIRONMENT_ACCESSES+".#", "0"),
				),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_permission_group" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestPermissionGroup"
						can_delete_config = true
					}
				`,
				ExpectError: regexp.MustCompile(`CanDeleteConfig is not allowed without the CanCreateOrUpdateConfig and CanDeleteSetting permissions.`),
			},
			{
				Config: `
					data "configcat_products" "products" {
					}
					resource "configcat_permission_group" "test" {
						product_id = data.configcat_products.products.products.0.product_id
						name = "TestPermissionGroup renamed"
						can_createorupdate_config = true
						can_createorupdate_environment = true
						can_delete_environment = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPermissionGroupResourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_NAME, "TestPermissionGroup renamed"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_MEMBERS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_CONFIG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_TAG_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_TAG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_SDKKEY, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_ROTATE_SDKKEY, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SEGMENTS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_ACCESSTYPE, "custom"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE, "none"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_ENVIRONMENT_ACCESSES+".#", "0"),
				),
			},
			{
				ResourceName:      testPermissionGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
