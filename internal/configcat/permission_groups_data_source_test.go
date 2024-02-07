package configcat

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermissionGroupsDataSource(t *testing.T) {
	const productId = "08d86d63-2721-4da6-8c06-584521d516bc"
	const testResourceName = "data.configcat_permission_groups.test"

	const permissionGroup1Id = "219"
	const permissionGroup2Id = "31126"
	const listAttribute = PermissionGroups

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable(productId),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "2"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("notfound"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("Administrators"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupId, permissionGroup1Id),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "Administrators"),

					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageMembers, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateConfig, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteConfig, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateEnvironment, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteEnvironment, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateSetting, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanTagSetting, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteSetting, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateTag, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteTag, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageWebhook, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanUseExportImport, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageProductPreferences, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageIntegrations, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanViewSdkKey, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanRotateSdkKey, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateSegment, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteSegment, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanViewProductAuditLogs, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanViewProductStatistics, "true"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupAccessType, "full"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupNewEnvironmentAccessType, "full"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupEnvironmentAccess+".%", "0"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id":        config.StringVariable(productId),
					"name_filter_regex": config.StringVariable("Only test environment"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".#", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupId, permissionGroup2Id),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+Name, "Only test environment"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageMembers, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateConfig, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteConfig, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateEnvironment, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteEnvironment, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateSetting, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanTagSetting, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteSetting, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateTag, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteTag, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageWebhook, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanUseExportImport, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageProductPreferences, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanManageIntegrations, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanViewSdkKey, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanRotateSdkKey, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanCreateOrUpdateSegment, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanDeleteSegment, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanViewProductAuditLogs, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupCanViewProductStatistics, "false"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupAccessType, "custom"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupNewEnvironmentAccessType, "readOnly"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupEnvironmentAccess+".%", "1"),
					resource.TestCheckResourceAttr(testResourceName, listAttribute+".0."+PermissionGroupEnvironmentAccess+".08d86d63-2726-47cd-8bfc-59608ecb91e2", "full"),
				),
			},
			{
				ConfigFile: config.TestNameFile("main.tf"),
				ConfigVariables: config.Variables{
					"product_id": config.StringVariable("invalidguid"),
				},
				ExpectError: regexp.MustCompile(`Attribute product_id value must be a valid GUID, got: invalidguid`),
			},
		},
	})
}
