package configcat

import (
	"context"
	"fmt"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatPermissionGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePermissionGroupCreate,
		ReadContext:   resourcePermissionGroupRead,
		UpdateContext: resourcePermissionGroupUpdate,
		DeleteContext: resourcePermissionGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
				ForceNew:     true,
			},

			PERMISSION_GROUP_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			PERMISSION_GROUP_CAN_MANAGE_MEMBERS: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_DELETE_CONFIG: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_TAG_SETTING: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_DELETE_SETTING: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_DELETE_TAG: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_MANAGE_WEBHOOK: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_USE_EXPORTIMPORT: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_VIEW_SDKKEY: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_ROTATE_SDKKEY: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_DELETE_SEGMENT: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_ACCESSTYPE: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  sw.ACCESSTYPE_CUSTOM,
			},
			PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  sw.ENVIRONMENTACCESSTYPE_NONE,
			},
			PERMISSION_GROUP_ENVIRONMENT_ACCESS: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ID: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateGUIDFunc,
						},
						PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ACCESS_TYPE: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  sw.ENVIRONMENTACCESSTYPE_NONE,
						},
					},
				},
			},
		},
	}
}

func resourcePermissionGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	productID := d.Get(PRODUCT_ID).(string)

	accessTypeString := d.Get(PERMISSION_GROUP_ACCESSTYPE).(string)
	accessType, accessTypeParseErr := sw.NewAccessTypeFromValue(accessTypeString)
	if accessTypeParseErr != nil {
		d.SetId("")
		return diag.FromErr(accessTypeParseErr)
	}

	newEnvironmentAccessTypeString := d.Get(PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE).(string)
	newEnvironmentAccessType, newEnvironmentAccessTypeParseErr := sw.NewEnvironmentAccessTypeFromValue(newEnvironmentAccessTypeString)
	if newEnvironmentAccessTypeParseErr != nil {
		d.SetId("")
		return diag.FromErr(newEnvironmentAccessTypeParseErr)
	}

	environmentAccesses, environmentAccessParseError := getEnvironmentAccesses(d.Get(PERMISSION_GROUP_ENVIRONMENT_ACCESS).([]interface{}), *accessType)
	if environmentAccessParseError != nil {
		return diag.FromErr(environmentAccessParseError)
	}

	canManageMembers := d.Get(PERMISSION_GROUP_CAN_MANAGE_MEMBERS).(bool)
	canCreateOrUpdateConfig := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG).(bool)
	canDeleteConfig := d.Get(PERMISSION_GROUP_CAN_DELETE_CONFIG).(bool)
	canCreateOrUpdateEnvironment := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT).(bool)
	canDeleteEnvironment := d.Get(PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT).(bool)
	canCreateOrUpdateSetting := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING).(bool)
	canTagSetting := d.Get(PERMISSION_GROUP_CAN_TAG_SETTING).(bool)
	canDeleteSetting := d.Get(PERMISSION_GROUP_CAN_DELETE_SETTING).(bool)
	canCreateOrUpdateTag := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG).(bool)
	canDeleteTag := d.Get(PERMISSION_GROUP_CAN_DELETE_TAG).(bool)
	canManageWebhook := d.Get(PERMISSION_GROUP_CAN_MANAGE_WEBHOOK).(bool)
	canUseExportImport := d.Get(PERMISSION_GROUP_CAN_USE_EXPORTIMPORT).(bool)
	canManageProductPreferences := d.Get(PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES).(bool)
	canManageIntegrations := d.Get(PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS).(bool)
	canViewSdkKey := d.Get(PERMISSION_GROUP_CAN_VIEW_SDKKEY).(bool)
	canRotateSdkKey := d.Get(PERMISSION_GROUP_CAN_ROTATE_SDKKEY).(bool)
	canCreateOrUpdateSegments := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT).(bool)
	canDeleteSegments := d.Get(PERMISSION_GROUP_CAN_DELETE_SEGMENT).(bool)
	canViewProductAuditLog := d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG).(bool)
	canViewProductStatistics := d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS).(bool)

	body := sw.CreatePermissionGroupRequest{
		Name:                         d.Get(ENVIRONMENT_NAME).(string),
		CanManageMembers:             &canManageMembers,
		CanCreateOrUpdateConfig:      &canCreateOrUpdateConfig,
		CanDeleteConfig:              &canDeleteConfig,
		CanCreateOrUpdateEnvironment: &canCreateOrUpdateEnvironment,
		CanDeleteEnvironment:         &canDeleteEnvironment,
		CanCreateOrUpdateSetting:     &canCreateOrUpdateSetting,
		CanTagSetting:                &canTagSetting,
		CanDeleteSetting:             &canDeleteSetting,
		CanCreateOrUpdateTag:         &canCreateOrUpdateTag,
		CanDeleteTag:                 &canDeleteTag,
		CanManageWebhook:             &canManageWebhook,
		CanUseExportImport:           &canUseExportImport,
		CanManageProductPreferences:  &canManageProductPreferences,
		CanManageIntegrations:        &canManageIntegrations,
		CanViewSdkKey:                &canViewSdkKey,
		CanRotateSdkKey:              &canRotateSdkKey,
		CanCreateOrUpdateSegments:    &canCreateOrUpdateSegments,
		CanDeleteSegments:            &canDeleteSegments,
		CanViewProductAuditLog:       &canViewProductAuditLog,
		CanViewProductStatistics:     &canViewProductStatistics,
		AccessType:                   accessType,
		NewEnvironmentAccessType:     newEnvironmentAccessType,
		EnvironmentAccesses:          *environmentAccesses,
	}

	permissionGroup, err := c.CreatePermissionGroup(productID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", *permissionGroup.PermissionGroupId))

	return resourcePermissionGroupRead(ctx, d, m)
}

func resourcePermissionGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	permissionGroupID, permissionGroupParseErr := strconv.ParseInt(d.Id(), 10, 32)
	if permissionGroupParseErr != nil {
		return diag.FromErr(permissionGroupParseErr)
	}

	permissionGroup, err := c.GetPermissionGroup(permissionGroupID)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(PRODUCT_ID, permissionGroup.Product.ProductId)
	d.Set(PERMISSION_GROUP_NAME, permissionGroup.Name.Get())
	d.Set(PERMISSION_GROUP_CAN_MANAGE_MEMBERS, permissionGroup.CanManageMembers)
	d.Set(PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG, permissionGroup.CanCreateOrUpdateConfig)
	d.Set(PERMISSION_GROUP_CAN_DELETE_CONFIG, permissionGroup.CanDeleteConfig)
	d.Set(PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT, permissionGroup.CanCreateOrUpdateEnvironment)
	d.Set(PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT, permissionGroup.CanDeleteEnvironment)
	d.Set(PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING, permissionGroup.CanCreateOrUpdateSetting)
	d.Set(PERMISSION_GROUP_CAN_TAG_SETTING, permissionGroup.CanTagSetting)
	d.Set(PERMISSION_GROUP_CAN_DELETE_SETTING, permissionGroup.CanDeleteSetting)
	d.Set(PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG, permissionGroup.CanCreateOrUpdateTag)
	d.Set(PERMISSION_GROUP_CAN_DELETE_TAG, permissionGroup.CanDeleteTag)
	d.Set(PERMISSION_GROUP_CAN_MANAGE_WEBHOOK, permissionGroup.CanManageWebhook)
	d.Set(PERMISSION_GROUP_CAN_USE_EXPORTIMPORT, permissionGroup.CanUseExportImport)
	d.Set(PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES, permissionGroup.CanManageProductPreferences)
	d.Set(PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS, permissionGroup.CanManageIntegrations)
	d.Set(PERMISSION_GROUP_CAN_VIEW_SDKKEY, permissionGroup.CanViewSdkKey)
	d.Set(PERMISSION_GROUP_CAN_ROTATE_SDKKEY, permissionGroup.CanRotateSdkKey)
	d.Set(PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT, permissionGroup.CanCreateOrUpdateSegments)
	d.Set(PERMISSION_GROUP_CAN_DELETE_SEGMENT, permissionGroup.CanDeleteSegments)
	d.Set(PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, permissionGroup.CanViewProductAuditLog)
	d.Set(PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, permissionGroup.CanViewProductStatistics)
	d.Set(PERMISSION_GROUP_ACCESSTYPE, permissionGroup.AccessType)
	d.Set(PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE, permissionGroup.NewEnvironmentAccessType)
	d.Set(PERMISSION_GROUP_ENVIRONMENT_ACCESS, flattenPermissionGroupEnvironmentAccessData(permissionGroup.EnvironmentAccesses, *permissionGroup.AccessType))

	return diags
}

func resourcePermissionGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.HasChanges(PERMISSION_GROUP_NAME,
		PERMISSION_GROUP_CAN_MANAGE_MEMBERS,
		PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG,
		PERMISSION_GROUP_CAN_DELETE_CONFIG,
		PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT,
		PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT,
		PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING,
		PERMISSION_GROUP_CAN_TAG_SETTING,
		PERMISSION_GROUP_CAN_DELETE_SETTING,
		PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG,
		PERMISSION_GROUP_CAN_DELETE_TAG,
		PERMISSION_GROUP_CAN_MANAGE_WEBHOOK,
		PERMISSION_GROUP_CAN_USE_EXPORTIMPORT,
		PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES,
		PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS,
		PERMISSION_GROUP_CAN_VIEW_SDKKEY,
		PERMISSION_GROUP_CAN_ROTATE_SDKKEY,
		PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT,
		PERMISSION_GROUP_CAN_DELETE_SEGMENT,
		PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG,
		PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS,
		PERMISSION_GROUP_ACCESSTYPE,
		PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE,
		PERMISSION_GROUP_ENVIRONMENT_ACCESS) {

		permimssionGroupName := d.Get(PERMISSION_GROUP_NAME).(string)

		accessTypeString := d.Get(PERMISSION_GROUP_ACCESSTYPE).(string)
		accessType, accessTypeParseErr := sw.NewAccessTypeFromValue(accessTypeString)
		if accessTypeParseErr != nil {
			return diag.FromErr(accessTypeParseErr)
		}

		newEnvironmentAccessTypeString := d.Get(PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE).(string)
		newEnvironmentAccessType, newEnvironmentAccessTypeParseErr := sw.NewEnvironmentAccessTypeFromValue(newEnvironmentAccessTypeString)
		if newEnvironmentAccessTypeParseErr != nil {
			return diag.FromErr(newEnvironmentAccessTypeParseErr)
		}

		environmentAccesses, environmentAccessParseError := getEnvironmentAccesses(d.Get(PERMISSION_GROUP_ENVIRONMENT_ACCESS).([]interface{}), *accessType)
		if environmentAccessParseError != nil {
			return diag.FromErr(environmentAccessParseError)
		}

		canManageMembers := d.Get(PERMISSION_GROUP_CAN_MANAGE_MEMBERS).(bool)
		canCreateOrUpdateConfig := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG).(bool)
		canDeleteConfig := d.Get(PERMISSION_GROUP_CAN_DELETE_CONFIG).(bool)
		canCreateOrUpdateEnvironment := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT).(bool)
		canDeleteEnvironment := d.Get(PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT).(bool)
		canCreateOrUpdateSetting := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING).(bool)
		canTagSetting := d.Get(PERMISSION_GROUP_CAN_TAG_SETTING).(bool)
		canDeleteSetting := d.Get(PERMISSION_GROUP_CAN_DELETE_SETTING).(bool)
		canCreateOrUpdateTag := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG).(bool)
		canDeleteTag := d.Get(PERMISSION_GROUP_CAN_DELETE_TAG).(bool)
		canManageWebhook := d.Get(PERMISSION_GROUP_CAN_MANAGE_WEBHOOK).(bool)
		canUseExportImport := d.Get(PERMISSION_GROUP_CAN_USE_EXPORTIMPORT).(bool)
		canManageProductPreferences := d.Get(PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES).(bool)
		canManageIntegrations := d.Get(PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS).(bool)
		canViewSdkKey := d.Get(PERMISSION_GROUP_CAN_VIEW_SDKKEY).(bool)
		canRotateSdkKey := d.Get(PERMISSION_GROUP_CAN_ROTATE_SDKKEY).(bool)
		canCreateOrUpdateSegments := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENT).(bool)
		canDeleteSegments := d.Get(PERMISSION_GROUP_CAN_DELETE_SEGMENT).(bool)
		canViewProductAuditLog := d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG).(bool)
		canViewProductStatistics := d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS).(bool)

		body := sw.UpdatePermissionGroupRequest{
			Name:                         *sw.NewNullableString(&permimssionGroupName),
			CanManageMembers:             *sw.NewNullableBool(&canManageMembers),
			CanCreateOrUpdateConfig:      *sw.NewNullableBool(&canCreateOrUpdateConfig),
			CanDeleteConfig:              *sw.NewNullableBool(&canDeleteConfig),
			CanCreateOrUpdateEnvironment: *sw.NewNullableBool(&canCreateOrUpdateEnvironment),
			CanDeleteEnvironment:         *sw.NewNullableBool(&canDeleteEnvironment),
			CanCreateOrUpdateSetting:     *sw.NewNullableBool(&canCreateOrUpdateSetting),
			CanTagSetting:                *sw.NewNullableBool(&canTagSetting),
			CanDeleteSetting:             *sw.NewNullableBool(&canDeleteSetting),
			CanCreateOrUpdateTag:         *sw.NewNullableBool(&canCreateOrUpdateTag),
			CanDeleteTag:                 *sw.NewNullableBool(&canDeleteTag),
			CanManageWebhook:             *sw.NewNullableBool(&canManageWebhook),
			CanUseExportImport:           *sw.NewNullableBool(&canUseExportImport),
			CanManageProductPreferences:  *sw.NewNullableBool(&canManageProductPreferences),
			CanManageIntegrations:        *sw.NewNullableBool(&canManageIntegrations),
			CanViewSdkKey:                *sw.NewNullableBool(&canViewSdkKey),
			CanRotateSdkKey:              *sw.NewNullableBool(&canRotateSdkKey),
			CanCreateOrUpdateSegments:    *sw.NewNullableBool(&canCreateOrUpdateSegments),
			CanDeleteSegments:            *sw.NewNullableBool(&canDeleteSegments),
			CanViewProductAuditLog:       *sw.NewNullableBool(&canViewProductAuditLog),
			CanViewProductStatistics:     *sw.NewNullableBool(&canViewProductStatistics),
			AccessType:                   accessType,
			NewEnvironmentAccessType:     newEnvironmentAccessType,
			EnvironmentAccesses:          *environmentAccesses,
		}

		permissionGroupID, permissionGroupParseErr := strconv.ParseInt(d.Id(), 10, 32)
		if permissionGroupParseErr != nil {
			return diag.FromErr(permissionGroupParseErr)
		}

		_, err := c.UpdatePermissionGroup(permissionGroupID, body)
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				d.SetId("")
				var diags diag.Diagnostics
				return diags
			}

			return diag.FromErr(err)
		}
	}

	return resourcePermissionGroupRead(ctx, d, m)
}

func getEnvironmentAccesses(environmentAccesses []interface{}, accessType sw.AccessType) (*[]sw.CreateOrUpdateEnvironmentAccessModel, error) {
	elements := make([]sw.CreateOrUpdateEnvironmentAccessModel, 0)

	if accessType != sw.ACCESSTYPE_CUSTOM && environmentAccesses != nil && len(environmentAccesses) > 0 {
		return nil, fmt.Errorf("Error: environment_access can only be set if the accesstype is custom")
	}

	for _, environmentAccess := range environmentAccesses {
		item := environmentAccess.(map[string]interface{})

		environmentAccessTypeString := item[PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ACCESS_TYPE].(string)
		environmentAccessType, environmentAccessTypeParseError := sw.NewEnvironmentAccessTypeFromValue(environmentAccessTypeString)
		if environmentAccessTypeParseError != nil {
			return nil, environmentAccessTypeParseError
		}

		environmentID := item[PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ID].(string)
		element := sw.CreateOrUpdateEnvironmentAccessModel{
			EnvironmentId:         &environmentID,
			EnvironmentAccessType: environmentAccessType,
		}

		elements = append(elements, element)
	}

	return &elements, nil
}

func resourcePermissionGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	permissionGroupID, permissionGroupParseErr := strconv.ParseInt(d.Id(), 10, 32)
	if permissionGroupParseErr != nil {
		return diag.FromErr(permissionGroupParseErr)
	}

	err := c.DeletePermissionGroup(permissionGroupID)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
