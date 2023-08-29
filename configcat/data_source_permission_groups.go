package configcat

import (
	"context"
	"regexp"
	"strconv"
	"time"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatPermissionGroups() *schema.Resource {
	return &schema.Resource{

		ReadContext: permissionGroupRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			PERMISSION_GROUP_NAME_FILTER_REGEX: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			PERMISSION_GROUPS: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PERMISSION_GROUP_ID: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						PERMISSION_GROUP_NAME: {
							Type:     schema.TypeString,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_MANAGE_MEMBERS: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_DELETE_CONFIG: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_TAG_SETTING: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_DELETE_SETTING: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_DELETE_TAG: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_MANAGE_WEBHOOK: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_USE_EXPORTIMPORT: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_VIEW_SDKKEY: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_ROTATE_SDKKEY: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_DELETE_SEGMENT: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PERMISSION_GROUP_ACCESSTYPE: {
							Type:     schema.TypeString,
							Computed: true,
						},
						PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE: {
							Type:     schema.TypeString,
							Computed: true,
						},
						PERMISSION_GROUP_ENVIRONMENT_ACCESSES: {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func permissionGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(PRODUCT_ID).(string)
	permissionGroupNameFilterRegex := d.Get(PERMISSION_GROUP_NAME_FILTER_REGEX).(string)

	permissionGroups, err := c.GetPermissionGroups(productID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredPermissionGroups := []sw.PermissionGroupModel{}
	if permissionGroupNameFilterRegex == "" {
		filteredPermissionGroups = permissionGroups
	} else {
		regex := regexp.MustCompile(permissionGroupNameFilterRegex)
		for i := range permissionGroups {
			if regex.MatchString(*permissionGroups[i].Name.Get()) {
				filteredPermissionGroups = append(filteredPermissionGroups, permissionGroups[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(PERMISSION_GROUPS, flattenPermissionGroupsData(&filteredPermissionGroups))

	var diags diag.Diagnostics
	return diags
}

func flattenPermissionGroupsData(permissionGroups *[]sw.PermissionGroupModel) []interface{} {
	if permissionGroups != nil {
		elements := make([]interface{}, len(*permissionGroups))

		for i, permissionGroup := range *permissionGroups {
			element := make(map[string]interface{})

			element[PERMISSION_GROUP_ID] = permissionGroup.PermissionGroupId
			element[PERMISSION_GROUP_NAME] = permissionGroup.Name.Get()
			element[PERMISSION_GROUP_CAN_MANAGE_MEMBERS] = permissionGroup.CanManageMembers
			element[PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG] = permissionGroup.CanCreateOrUpdateConfig
			element[PERMISSION_GROUP_CAN_DELETE_CONFIG] = permissionGroup.CanDeleteConfig
			element[PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT] = permissionGroup.CanCreateOrUpdateEnvironment
			element[PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT] = permissionGroup.CanDeleteEnvironment
			element[PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING] = permissionGroup.CanCreateOrUpdateSetting
			element[PERMISSION_GROUP_CAN_TAG_SETTING] = permissionGroup.CanTagSetting
			element[PERMISSION_GROUP_CAN_DELETE_SETTING] = permissionGroup.CanDeleteSetting
			element[PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG] = permissionGroup.CanCreateOrUpdateTag
			element[PERMISSION_GROUP_CAN_DELETE_TAG] = permissionGroup.CanDeleteTag
			element[PERMISSION_GROUP_CAN_MANAGE_WEBHOOK] = permissionGroup.CanManageWebhook
			element[PERMISSION_GROUP_CAN_USE_EXPORTIMPORT] = permissionGroup.CanUseExportImport
			element[PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES] = permissionGroup.CanManageProductPreferences
			element[PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS] = permissionGroup.CanManageIntegrations
			element[PERMISSION_GROUP_CAN_VIEW_SDKKEY] = permissionGroup.CanViewSdkKey
			element[PERMISSION_GROUP_CAN_ROTATE_SDKKEY] = permissionGroup.CanRotateSdkKey
			element[PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT] = permissionGroup.CanCreateOrUpdateSegments
			element[PERMISSION_GROUP_CAN_DELETE_SEGMENT] = permissionGroup.CanDeleteSegments
			element[PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG] = permissionGroup.CanViewProductAuditLog
			element[PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS] = permissionGroup.CanViewProductStatistics
			element[PERMISSION_GROUP_ACCESSTYPE] = *permissionGroup.AccessType
			element[PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE] = *permissionGroup.NewEnvironmentAccessType
			element[PERMISSION_GROUP_ENVIRONMENT_ACCESSES] = flattenPermissionGroupEnvironmentAccessData(permissionGroup.EnvironmentAccesses, *permissionGroup.AccessType)

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}

func flattenPermissionGroupEnvironmentAccessData(environmentAccesses []sw.EnvironmentAccessModel, accessType sw.AccessType) map[string]any {
	elements := make(map[string]any)
	if accessType != sw.ACCESSTYPE_CUSTOM {
		return elements
	}

	for _, environmentAccess := range environmentAccesses {
		if *environmentAccess.EnvironmentAccessType == sw.ENVIRONMENTACCESSTYPE_NONE {
			continue
		}

		elements[*environmentAccess.EnvironmentId] = *environmentAccess.EnvironmentAccessType
	}

	return elements
}
