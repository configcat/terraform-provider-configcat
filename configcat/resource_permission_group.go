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
			PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			PERMISSION_GROUP_CAN_DELETE_SEGMENTS: {
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
				Type:    schema.TypeString,
				Default: sw.ACCESSTYPE_CUSTOM,
			},
			PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE: {
				Type:    schema.TypeString,
				Default: sw.ENVIRONMENTACCESSTYPE_NONE,
			},
			PERMISSION_GROUP_ENVIRONMENT_ACCESSES: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ID: {
							Type:     schema.TypeString,
							Required: true,
						},
						PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ACCESS_TYPE: {
							Type:     schema.TypeString,
							Required: true,
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

	environmentAccesses, environmentAccessParseError := getEnvironmentAccesses(d.Get(PERMISSION_GROUP_ENVIRONMENT_ACCESSES).([]interface{}))
	if environmentAccessParseError != nil {
		return diag.FromErr(environmentAccessParseError)
	}

	body := sw.CreatePermissionGroupRequest{
		Name:                         d.Get(ENVIRONMENT_NAME).(string),
		CanManageMembers:             d.Get(PERMISSION_GROUP_CAN_MANAGE_MEMBERS).(*bool),
		CanCreateOrUpdateConfig:      d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG).(*bool),
		CanDeleteConfig:              d.Get(PERMISSION_GROUP_CAN_DELETE_CONFIG).(*bool),
		CanCreateOrUpdateEnvironment: d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT).(*bool),
		CanDeleteEnvironment:         d.Get(PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT).(*bool),
		CanCreateOrUpdateSetting:     d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING).(*bool),
		CanTagSetting:                d.Get(PERMISSION_GROUP_CAN_TAG_SETTING).(*bool),
		CanDeleteSetting:             d.Get(PERMISSION_GROUP_CAN_DELETE_SETTING).(*bool),
		CanCreateOrUpdateTag:         d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG).(*bool),
		CanDeleteTag:                 d.Get(PERMISSION_GROUP_CAN_DELETE_TAG).(*bool),
		CanManageWebhook:             d.Get(PERMISSION_GROUP_CAN_MANAGE_WEBHOOK).(*bool),
		CanUseExportImport:           d.Get(PERMISSION_GROUP_CAN_USE_EXPORTIMPORT).(*bool),
		CanManageProductPreferences:  d.Get(PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES).(*bool),
		CanManageIntegrations:        d.Get(PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS).(*bool),
		CanViewSdkKey:                d.Get(PERMISSION_GROUP_CAN_VIEW_SDKKEY).(*bool),
		CanRotateSdkKey:              d.Get(PERMISSION_GROUP_CAN_ROTATE_SDKKEY).(*bool),
		CanCreateOrUpdateSegments:    d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS).(*bool),
		CanDeleteSegments:            d.Get(PERMISSION_GROUP_CAN_DELETE_SEGMENTS).(*bool),
		CanViewProductAuditLog:       d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG).(*bool),
		CanViewProductStatistics:     d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS).(*bool),
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
	d.Set(PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS, permissionGroup.CanCreateOrUpdateSegments)
	d.Set(PERMISSION_GROUP_CAN_DELETE_SEGMENTS, permissionGroup.CanDeleteSegments)
	d.Set(PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG, permissionGroup.CanViewProductAuditLog)
	d.Set(PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS, permissionGroup.CanViewProductStatistics)
	d.Set(PERMISSION_GROUP_ACCESSTYPE, permissionGroup.AccessType)
	d.Set(PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE, permissionGroup.NewEnvironmentAccessType)
	d.Set(PERMISSION_GROUP_ENVIRONMENT_ACCESSES, flattenPermissionGroupEnvironmentAccessData(permissionGroup.EnvironmentAccesses))

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
		PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS,
		PERMISSION_GROUP_CAN_DELETE_SEGMENTS,
		PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG,
		PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS,
		PERMISSION_GROUP_ACCESSTYPE,
		PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE,
		PERMISSION_GROUP_ENVIRONMENT_ACCESSES) {

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

		environmentAccesses, environmentAccessParseError := getEnvironmentAccesses(d.Get(PERMISSION_GROUP_ENVIRONMENT_ACCESSES).([]interface{}))
		if environmentAccessParseError != nil {
			return diag.FromErr(environmentAccessParseError)
		}

		body := sw.UpdatePermissionGroupRequest{
			Name:                         *sw.NewNullableString(&permimssionGroupName),
			CanManageMembers:             *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_MANAGE_MEMBERS).(*bool)),
			CanCreateOrUpdateConfig:      *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG).(*bool)),
			CanDeleteConfig:              *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_DELETE_CONFIG).(*bool)),
			CanCreateOrUpdateEnvironment: *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT).(*bool)),
			CanDeleteEnvironment:         *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT).(*bool)),
			CanCreateOrUpdateSetting:     *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING).(*bool)),
			CanTagSetting:                *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_TAG_SETTING).(*bool)),
			CanDeleteSetting:             *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_DELETE_SETTING).(*bool)),
			CanCreateOrUpdateTag:         *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG).(*bool)),
			CanDeleteTag:                 *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_DELETE_TAG).(*bool)),
			CanManageWebhook:             *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_MANAGE_WEBHOOK).(*bool)),
			CanUseExportImport:           *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_USE_EXPORTIMPORT).(*bool)),
			CanManageProductPreferences:  *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES).(*bool)),
			CanManageIntegrations:        *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS).(*bool)),
			CanViewSdkKey:                *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_VIEW_SDKKEY).(*bool)),
			CanRotateSdkKey:              *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_ROTATE_SDKKEY).(*bool)),
			CanCreateOrUpdateSegments:    *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS).(*bool)),
			CanDeleteSegments:            *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_DELETE_SEGMENTS).(*bool)),
			CanViewProductAuditLog:       *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG).(*bool)),
			CanViewProductStatistics:     *sw.NewNullableBool(d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS).(*bool)),
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

func getEnvironmentAccesses(environmentAccesses []interface{}) (*[]sw.CreateOrUpdateEnvironmentAccessModel, error) {
	if environmentAccesses != nil {
		elements := make([]sw.CreateOrUpdateEnvironmentAccessModel, len(environmentAccesses))

		for i, environmentAccess := range environmentAccesses {
			item := environmentAccess.(map[string]interface{})

			environmentAccessTypeString := item[PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ACCESS_TYPE].(string)
			environmentAccessType, environmentAccessTypeParseError := sw.NewEnvironmentAccessTypeFromValue(environmentAccessTypeString)
			if environmentAccessTypeParseError != nil {
				return nil, environmentAccessTypeParseError
			}

			element := sw.CreateOrUpdateEnvironmentAccessModel{
				EnvironmentId:         item[PERMISSION_GROUP_ENVIRONMENT_ACCESS_ENVIRONMENT_ID].(*string),
				EnvironmentAccessType: environmentAccessType,
			}

			elements[i] = element
		}

		return &elements, nil
	}
	empty := make([]sw.CreateOrUpdateEnvironmentAccessModel, 0)
	return &empty, nil
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
