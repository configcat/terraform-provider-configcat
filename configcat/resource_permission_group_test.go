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
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SEGMENT, "false"),
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
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SEGMENT, "false"),
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
						name = "TestPermissionGroup renamed 2"
						can_manage_members = true
						can_createorupdate_config = true
						can_delete_config = false
						can_createorupdate_environment = true
						can_delete_environment = false
						can_createorupdate_setting = true
						can_tag_setting = true
						can_delete_setting = false
						can_createorupdate_tag = true
						can_delete_tag = false
						can_manage_webhook = true
						can_use_exportimport = false
						can_manage_product_preferences = true
						can_manage_integrations = false
						can_view_sdkkey = true
						can_rotate_sdkkey = false
						can_createorupdate_segment = true
						can_delete_segment = false
						can_view_product_auditlog = true
						can_view_product_statistics = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPermissionGroupResourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_NAME, "TestPermissionGroup renamed 2"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_MEMBERS, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_CONFIG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_TAG_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SETTING, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_TAG, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_SDKKEY, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_ROTATE_SDKKEY, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SEGMENT, "false"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, "true"),
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
						name = "TestPermissionGroup renamed 2"
						can_manage_members = true
						can_createorupdate_config = true
						can_delete_config = true
						can_createorupdate_environment = true
						can_delete_environment = true
						can_createorupdate_setting = true
						can_tag_setting = true
						can_delete_setting = true
						can_createorupdate_tag = true
						can_delete_tag = true
						can_manage_webhook = true
						can_use_exportimport = true
						can_manage_product_preferences = true
						can_manage_integrations = true
						can_view_sdkkey = true
						can_rotate_sdkkey = true
						can_createorupdate_segment = true
						can_delete_segment = true
						can_view_product_auditlog = true
						can_view_product_statistics = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testPermissionGroupResourceName, "id"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_NAME, "TestPermissionGroup renamed 2"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_MEMBERS, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_CONFIG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_TAG_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SETTING, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_TAG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_SDKKEY, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_ROTATE_SDKKEY, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_DELETE_SEGMENT, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, "true"),
					resource.TestCheckResourceAttr(testPermissionGroupResourceName, PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, "true"),
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

func TestResourcePermissionGroupApiErrorFlow(t *testing.T) {
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
		},
	})
}
