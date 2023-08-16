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
				Type:     schema.TypeString,
				Required: true,
			},
			PERMISSION_GROUP_NEW_ENVIRONMENT_ACCESSTYPE: {
				Type:     schema.TypeString,
				Optional: true,
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

	canManageMembers := d.Get(PERMISSION_GROUP_CAN_MANAGE_MEMBERS).(*bool)
	canCreateOrUpdateConfig := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_CONFIG).(*bool)
	canDeleteConfig := d.Get(PERMISSION_GROUP_CAN_DELETE_CONFIG).(*bool)
	canCreateOrUpdateEnvironment := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_ENVIRONMENT).(*bool)
	canDeleteEnvironment := d.Get(PERMISSION_GROUP_CAN_DELETE_ENVIRONMENT).(*bool)
	canCreateOrUpdateSetting := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SETTING).(*bool)
	canTagSetting := d.Get(PERMISSION_GROUP_CAN_TAG_SETTING).(*bool)
	canDeleteSetting := d.Get(PERMISSION_GROUP_CAN_DELETE_SETTING).(*bool)
	canCreateOrUpdateTag := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_TAG).(*bool)
	canDeleteTag := d.Get(PERMISSION_GROUP_CAN_DELETE_TAG).(*bool)
	canManageWebhook := d.Get(PERMISSION_GROUP_CAN_MANAGE_WEBHOOK).(*bool)
	canUseExportImport := d.Get(PERMISSION_GROUP_CAN_USE_EXPORTIMPORT).(*bool)
	canManageProductPreferences := d.Get(PERMISSION_GROUP_CAN_MANAGE_PRODUCT_PREFERENCES).(*bool)
	canManageIntegrations := d.Get(PERMISSION_GROUP_CAN_MANAGE_INTEGRATIONS).(*bool)
	canViewSdkKey := d.Get(PERMISSION_GROUP_CAN_VIEW_SDKKEY).(*bool)
	canRotateSdkKey := d.Get(PERMISSION_GROUP_CAN_ROTATE_SDKKEY).(*bool)
	canCreateOrUpdateSegments := d.Get(PERMISSION_GROUP_CAN_CREATEORUPDATE_SEGMENTS).(*bool)
	canDeleteSegments := d.Get(PERMISSION_GROUP_CAN_DELETE_SEGMENTS).(*bool)
	canViewProductAuditLog := d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_AUDITLOG).(*bool)
	canViewProductStatistics := d.Get(PERMISSION_GROUP_CAN_VIEW_PRODUCT_STATISTICS).(*bool)

	accessTypeString := d.Get(PERMISSION_GROUP_ACCESSTYPE).(string)
	accessType, err := sw.NewAccessTypeFromValue(accessTypeString)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	body := sw.CreatePermissionGroupRequest{
		Name:                         d.Get(ENVIRONMENT_NAME).(string),
		CanManageMembers:             canManageMembers,
		CanCreateOrUpdateConfig:      canCreateOrUpdateConfig,
		CanDeleteConfig:              canDeleteConfig,
		CanCreateOrUpdateEnvironment: canCreateOrUpdateEnvironment,
		CanDeleteEnvironment:         canDeleteEnvironment,
		CanCreateOrUpdateSetting:     canCreateOrUpdateSetting,
		CanTagSetting:                canTagSetting,
		CanDeleteSetting:             canDeleteSetting,
		CanCreateOrUpdateTag:         canCreateOrUpdateTag,
		CanDeleteTag:                 canDeleteTag,
		CanManageWebhook:             canManageWebhook,
		CanUseExportImport:           canUseExportImport,
		CanManageProductPreferences:  canManageProductPreferences,
		CanManageIntegrations:        canManageIntegrations,
		CanViewSdkKey:                canViewSdkKey,
		CanRotateSdkKey:              canRotateSdkKey,
		CanCreateOrUpdateSegments:    canCreateOrUpdateSegments,
		CanDeleteSegments:            canDeleteSegments,
		CanViewProductAuditLog:       canViewProductAuditLog,
		CanViewProductStatistics:     canViewProductStatistics,
		AccessType:                   accessType,
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

	environment, err := c.GetPermissionGroup(permissionGroupID)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(PRODUCT_ID, environment.Product.ProductId)
	d.Set(ENVIRONMENT_NAME, environment.Name.Get())

	return diags
}

func resourcePermissionGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.HasChanges(ENVIRONMENT_NAME, ENVIRONMENT_DESCRIPTION, ENVIRONMENT_COLOR) {
		environmentName := d.Get(ENVIRONMENT_NAME).(string)

		body := sw.UpdatePermissionGroupRequest{
			Name: *sw.NewNullableString(&environmentName),
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
