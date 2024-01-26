package configcat

import (
	"context"
	"fmt"
	"strconv"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var _ resource.Resource = &permissionGroupResource{}
var _ resource.ResourceWithImportState = &permissionGroupResource{}

func NewPermissionGroupResource() resource.Resource {
	return &permissionGroupResource{}
}

type permissionGroupResource struct {
	client *client.Client
}

type permissionGroupResourceModel struct {
	ProductId types.String `tfsdk:"product_id"`

	ID                           types.String `tfsdk:"permission_group_id"`
	Name                         types.String `tfsdk:"name"`
	CanManageMembers             types.Bool   `tfsdk:"can_manage_members"`
	CanCreateOrUpdateConfig      types.Bool   `tfsdk:"can_createorupdate_config"`
	CanDeleteConfig              types.Bool   `tfsdk:"can_delete_config"`
	CanCreateOrUpdateEnvironment types.Bool   `tfsdk:"can_createorupdate_environment"`
	CanDeleteEnvironment         types.Bool   `tfsdk:"can_delete_environment"`
	CanCreateOrUpdateSetting     types.Bool   `tfsdk:"can_createorupdate_setting"`
	CanTagSetting                types.Bool   `tfsdk:"can_tag_setting"`
	CanDeleteSetting             types.Bool   `tfsdk:"can_delete_setting"`
	CanCreateOrUpdateTag         types.Bool   `tfsdk:"can_createorupdate_tag"`
	CanDeleteTag                 types.Bool   `tfsdk:"can_delete_tag"`
	CanManageWebhook             types.Bool   `tfsdk:"can_manage_webhook"`
	CanUseExportImport           types.Bool   `tfsdk:"can_use_exportimport"`
	CanManageProductPreferences  types.Bool   `tfsdk:"can_manage_product_preferences"`
	CanManageIntegrations        types.Bool   `tfsdk:"can_manage_integrations"`
	CanViewSdkKey                types.Bool   `tfsdk:"can_view_sdkkey"`
	CanRotateSdkKey              types.Bool   `tfsdk:"can_rotate_sdkkey"`
	CanCreateOrUpdateSegment     types.Bool   `tfsdk:"can_createorupdate_segment"`
	CanDeleteSegment             types.Bool   `tfsdk:"can_delete_segment"`
	CanViewProductAuditLogs      types.Bool   `tfsdk:"can_view_product_auditlog"`
	CanViewProductStatistics     types.Bool   `tfsdk:"can_view_product_statistics"`
	AccessType                   types.String `tfsdk:"accesstype"`
	NewEnvironmentAccessType     types.String `tfsdk:"new_environment_accesstype"`
	EnvironmentAccess            types.Map    `tfsdk:"environment_accesses"`
}

func (r *permissionGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission_group"
}

func (r *permissionGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages a **" + PermissionGroupResourceName + "**. [What is a " + PermissionGroupResourceName + " in ConfigCat?](https://configcat.com/docs/advanced/team-management/team-management-basics/#permissions--permission-groups-product-level)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + PermissionGroupResourceName + ".",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			ProductId: schema.StringAttribute{
				Description: "The ID of the " + ProductResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			Name: schema.StringAttribute{
				Description: "The name of the " + PermissionGroupResourceName + ".",
				Required:    true,
			},

			PermissionGroupCanManageMembers: schema.BoolAttribute{
				Description: "Group members can manage team members.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanCreateOrUpdateConfig: schema.BoolAttribute{
				Description: "Group members can create/update Configs.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanDeleteConfig: schema.BoolAttribute{
				Description: "Group members can delete Configs.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanCreateOrUpdateEnvironment: schema.BoolAttribute{
				Description: "Group members can create/update Environments.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanDeleteEnvironment: schema.BoolAttribute{
				Description: "Group members can delete Environments.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanCreateOrUpdateSetting: schema.BoolAttribute{
				Description: "Group members can create/update Feature Flags and Settings.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanTagSetting: schema.BoolAttribute{
				Description: "Group members can attach/detach Tags to Feature Flags and Settings.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanDeleteSetting: schema.BoolAttribute{
				Description: "Group members can delete Feature Flags and Settings.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanCreateOrUpdateTag: schema.BoolAttribute{
				Description: "Group members can create/update Tags.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanDeleteTag: schema.BoolAttribute{
				Description: "Group members can delete Tags.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanManageWebhook: schema.BoolAttribute{
				Description: "Group members can create/update/delete Webhooks.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanUseExportImport: schema.BoolAttribute{
				Description: "Group members can use the export/import feature.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanManageProductPreferences: schema.BoolAttribute{
				Description: "Group members can update Product preferences.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanManageIntegrations: schema.BoolAttribute{
				Description: "Group members can add and configure integrations.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanViewSdkKey: schema.BoolAttribute{
				Description: "Group members has access to SDK keys.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanRotateSdkKey: schema.BoolAttribute{
				Description: "Group members can rotate SDK keys.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanCreateOrUpdateSegment: schema.BoolAttribute{
				Description: "Group members can create/update Segments.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanDeleteSegment: schema.BoolAttribute{
				Description: "Group members can delete Segments.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanViewProductAuditLogs: schema.BoolAttribute{
				Description: "Group members has access to audit logs.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupCanViewProductStatistics: schema.BoolAttribute{
				Description: "Group members has access to product statistics.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			PermissionGroupAccessType: schema.StringAttribute{
				Description: "Represent the Feature Management permission. Possible values: readOnly, full, custom",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(string(sw.ACCESSTYPE_CUSTOM)),
			},
			PermissionGroupNewEnvironmentAccessType: schema.StringAttribute{
				Description: "Represent the environment specific Feature Management permission for new Environments. Possible values: full, readOnly, none",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(string(sw.ENVIRONMENTACCESSTYPE_NONE)),
			},
			PermissionGroupEnvironmentAccess: schema.MapAttribute{
				Description: "The environment specific permissions map block. Keys are the Environment IDs and the values represent the environment specific Feature Management permission. Possible values: full, readOnly",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *permissionGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *permissionGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	accessTypeString := plan.AccessType.ValueString()
	accessType, accessTypeParseErr := sw.NewAccessTypeFromValue(accessTypeString)
	if accessTypeParseErr != nil {
		resp.Diagnostics.AddAttributeError(path.Root(PermissionGroupAccessType), "invalid accesstype", accessTypeParseErr.Error())
		return
	}

	newEnvironmentAccessTypeString := plan.NewEnvironmentAccessType.ValueString()
	newEnvironmentAccessType, newEnvironmentAccessTypeParseErr := sw.NewEnvironmentAccessTypeFromValue(newEnvironmentAccessTypeString)
	if newEnvironmentAccessTypeParseErr != nil {
		resp.Diagnostics.AddAttributeError(path.Root(PermissionGroupNewEnvironmentAccessType), "invalid new_environment_accesstype", newEnvironmentAccessTypeParseErr.Error())
		return
	}

	environmentAccessMap := make(map[string]types.String, len(plan.EnvironmentAccess.Elements()))
	resp.Diagnostics.Append(plan.EnvironmentAccess.ElementsAs(ctx, &environmentAccessMap, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environmentAccesses, environmentAccessParseError := getEnvironmentAccesses(environmentAccessMap, nil, *accessType)
	if environmentAccessParseError != nil {
		resp.Diagnostics.AddAttributeError(path.Root(PermissionGroupEnvironmentAccess), "invalid environment_accesses", environmentAccessParseError.Error())
		return
	}

	body := sw.CreatePermissionGroupRequest{
		Name:                         plan.Name.ValueString(),
		CanManageMembers:             plan.CanManageMembers.ValueBoolPointer(),
		CanCreateOrUpdateConfig:      plan.CanCreateOrUpdateConfig.ValueBoolPointer(),
		CanDeleteConfig:              plan.CanDeleteConfig.ValueBoolPointer(),
		CanCreateOrUpdateEnvironment: plan.CanCreateOrUpdateEnvironment.ValueBoolPointer(),
		CanDeleteEnvironment:         plan.CanDeleteEnvironment.ValueBoolPointer(),
		CanCreateOrUpdateSetting:     plan.CanCreateOrUpdateSetting.ValueBoolPointer(),
		CanTagSetting:                plan.CanTagSetting.ValueBoolPointer(),
		CanDeleteSetting:             plan.CanDeleteSetting.ValueBoolPointer(),
		CanCreateOrUpdateTag:         plan.CanCreateOrUpdateTag.ValueBoolPointer(),
		CanDeleteTag:                 plan.CanDeleteTag.ValueBoolPointer(),
		CanManageWebhook:             plan.CanManageWebhook.ValueBoolPointer(),
		CanUseExportImport:           plan.CanUseExportImport.ValueBoolPointer(),
		CanManageProductPreferences:  plan.CanManageProductPreferences.ValueBoolPointer(),
		CanManageIntegrations:        plan.CanManageIntegrations.ValueBoolPointer(),
		CanViewSdkKey:                plan.CanViewSdkKey.ValueBoolPointer(),
		CanRotateSdkKey:              plan.CanRotateSdkKey.ValueBoolPointer(),
		CanCreateOrUpdateSegments:    plan.CanCreateOrUpdateSegment.ValueBoolPointer(),
		CanDeleteSegments:            plan.CanDeleteSegment.ValueBoolPointer(),
		CanViewProductAuditLog:       plan.CanViewProductAuditLogs.ValueBoolPointer(),
		CanViewProductStatistics:     plan.CanViewProductStatistics.ValueBoolPointer(),
		AccessType:                   accessType,
		NewEnvironmentAccessType:     newEnvironmentAccessType,
		EnvironmentAccesses:          *environmentAccesses,
	}

	model, err := r.client.CreatePermissionGroup(plan.ProductId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+PermissionGroupResourceName+", got error: %s", err))
		return
	}

	resp.Diagnostics.Append(plan.UpdateFromApiModel(ctx, *model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *permissionGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissionGroupId, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Permission Group ID", convErr.Error())
		return
	}

	model, err := r.client.GetPermissionGroup(permissionGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+PermissionGroupResourceName+", got error: %s", err))
		return
	}

	resp.Diagnostics.Append(state.UpdateFromApiModel(ctx, *model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *permissionGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan permissionGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.Equal(state.Name) &&
		plan.CanManageMembers.Equal(state.CanManageMembers) &&
		plan.CanCreateOrUpdateConfig.Equal(state.CanCreateOrUpdateConfig) &&
		plan.CanDeleteConfig.Equal(state.CanDeleteConfig) &&
		plan.CanCreateOrUpdateEnvironment.Equal(state.CanCreateOrUpdateEnvironment) &&
		plan.CanDeleteEnvironment.Equal(state.CanDeleteEnvironment) &&
		plan.CanCreateOrUpdateSetting.Equal(state.CanCreateOrUpdateSetting) &&
		plan.CanTagSetting.Equal(state.CanTagSetting) &&
		plan.CanDeleteSetting.Equal(state.CanDeleteSetting) &&
		plan.CanCreateOrUpdateTag.Equal(state.CanCreateOrUpdateTag) &&
		plan.CanDeleteTag.Equal(state.CanDeleteTag) &&
		plan.CanManageWebhook.Equal(state.CanManageWebhook) &&
		plan.CanUseExportImport.Equal(state.CanUseExportImport) &&
		plan.CanManageProductPreferences.Equal(state.CanManageProductPreferences) &&
		plan.CanManageIntegrations.Equal(state.CanManageIntegrations) &&
		plan.CanViewSdkKey.Equal(state.CanViewSdkKey) &&
		plan.CanRotateSdkKey.Equal(state.CanRotateSdkKey) &&
		plan.CanCreateOrUpdateSegment.Equal(state.CanCreateOrUpdateSegment) &&
		plan.CanDeleteSegment.Equal(state.CanDeleteSegment) &&
		plan.CanViewProductAuditLogs.Equal(state.CanViewProductAuditLogs) &&
		plan.CanViewProductStatistics.Equal(state.CanViewProductStatistics) &&
		plan.AccessType.Equal(state.AccessType) &&
		plan.NewEnvironmentAccessType.Equal(state.NewEnvironmentAccessType) &&
		plan.EnvironmentAccess.Equal(state.EnvironmentAccess) {
		return
	}

	permissionGroupId, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Permission Group ID", convErr.Error())
		return
	}

	accessTypeString := plan.AccessType.ValueString()
	accessType, accessTypeParseErr := sw.NewAccessTypeFromValue(accessTypeString)
	if accessTypeParseErr != nil {
		resp.Diagnostics.AddAttributeError(path.Root(PermissionGroupAccessType), "invalid accesstype", accessTypeParseErr.Error())
		return
	}

	newEnvironmentAccessTypeString := plan.NewEnvironmentAccessType.ValueString()
	newEnvironmentAccessType, newEnvironmentAccessTypeParseErr := sw.NewEnvironmentAccessTypeFromValue(newEnvironmentAccessTypeString)
	if newEnvironmentAccessTypeParseErr != nil {
		resp.Diagnostics.AddAttributeError(path.Root(PermissionGroupNewEnvironmentAccessType), "invalid new_environment_accesstype", newEnvironmentAccessTypeParseErr.Error())
		return
	}

	environmentAccessMap := make(map[string]types.String, len(plan.EnvironmentAccess.Elements()))
	resp.Diagnostics.Append(plan.EnvironmentAccess.ElementsAs(ctx, &environmentAccessMap, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environmentAccesses, environmentAccessParseError := getEnvironmentAccesses(environmentAccessMap, nil, *accessType)
	if environmentAccessParseError != nil {
		resp.Diagnostics.AddAttributeError(path.Root(PermissionGroupEnvironmentAccess), "invalid environment_accesses", environmentAccessParseError.Error())
		return
	}

	body := sw.UpdatePermissionGroupRequest{
		Name:                         *sw.NewNullableString(plan.Name.ValueStringPointer()),
		CanManageMembers:             *sw.NewNullableBool(plan.CanManageMembers.ValueBoolPointer()),
		CanCreateOrUpdateConfig:      *sw.NewNullableBool(plan.CanCreateOrUpdateConfig.ValueBoolPointer()),
		CanDeleteConfig:              *sw.NewNullableBool(plan.CanDeleteConfig.ValueBoolPointer()),
		CanCreateOrUpdateEnvironment: *sw.NewNullableBool(plan.CanCreateOrUpdateEnvironment.ValueBoolPointer()),
		CanDeleteEnvironment:         *sw.NewNullableBool(plan.CanDeleteEnvironment.ValueBoolPointer()),
		CanCreateOrUpdateSetting:     *sw.NewNullableBool(plan.CanCreateOrUpdateSetting.ValueBoolPointer()),
		CanTagSetting:                *sw.NewNullableBool(plan.CanTagSetting.ValueBoolPointer()),
		CanDeleteSetting:             *sw.NewNullableBool(plan.CanDeleteSetting.ValueBoolPointer()),
		CanCreateOrUpdateTag:         *sw.NewNullableBool(plan.CanCreateOrUpdateTag.ValueBoolPointer()),
		CanDeleteTag:                 *sw.NewNullableBool(plan.CanDeleteTag.ValueBoolPointer()),
		CanManageWebhook:             *sw.NewNullableBool(plan.CanManageWebhook.ValueBoolPointer()),
		CanUseExportImport:           *sw.NewNullableBool(plan.CanUseExportImport.ValueBoolPointer()),
		CanManageProductPreferences:  *sw.NewNullableBool(plan.CanManageProductPreferences.ValueBoolPointer()),
		CanManageIntegrations:        *sw.NewNullableBool(plan.CanManageIntegrations.ValueBoolPointer()),
		CanViewSdkKey:                *sw.NewNullableBool(plan.CanViewSdkKey.ValueBoolPointer()),
		CanRotateSdkKey:              *sw.NewNullableBool(plan.CanRotateSdkKey.ValueBoolPointer()),
		CanCreateOrUpdateSegments:    *sw.NewNullableBool(plan.CanCreateOrUpdateSegment.ValueBoolPointer()),
		CanDeleteSegments:            *sw.NewNullableBool(plan.CanDeleteSegment.ValueBoolPointer()),
		CanViewProductAuditLog:       *sw.NewNullableBool(plan.CanViewProductAuditLogs.ValueBoolPointer()),
		CanViewProductStatistics:     *sw.NewNullableBool(plan.CanViewProductStatistics.ValueBoolPointer()),
		AccessType:                   accessType,
		NewEnvironmentAccessType:     newEnvironmentAccessType,
		EnvironmentAccesses:          *environmentAccesses,
	}

	model, err := r.client.UpdatePermissionGroup(permissionGroupId, body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+PermissionGroupResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(ctx, *model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *permissionGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state permissionGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissionGroupId, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Permission Group ID", convErr.Error())
		return
	}

	err := r.client.DeletePermissionGroup(permissionGroupId)

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, PermissionGroupResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+PermissionGroupResourceName+", got error: %s", err))
		return
	}
}

func (r *permissionGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *permissionGroupResourceModel) UpdateFromApiModel(ctx context.Context, model sw.PermissionGroupModel) diag.Diagnostics {
	var diags diag.Diagnostics

	environmentAccesses := make(map[string]string, len(model.EnvironmentAccesses))
	for _, environmentAccess := range model.EnvironmentAccesses {
		if *environmentAccess.EnvironmentAccessType == sw.ENVIRONMENTACCESSTYPE_NONE {
			continue
		}

		environmentAccesses[*environmentAccess.EnvironmentId] = (string)(*environmentAccess.EnvironmentAccessType)
	}

	environmentAccessesMapValue, diags := types.MapValueFrom(ctx, types.StringType, environmentAccesses)
	if diags.HasError() {
		return diags
	}

	resourceModel.ID = types.StringValue(strconv.FormatInt(*model.PermissionGroupId, 10))
	resourceModel.Name = types.StringPointerValue(model.Name.Get())
	resourceModel.CanManageMembers = types.BoolPointerValue(model.CanManageMembers)
	resourceModel.CanCreateOrUpdateConfig = types.BoolPointerValue(model.CanCreateOrUpdateConfig)
	resourceModel.CanDeleteConfig = types.BoolPointerValue(model.CanDeleteConfig)
	resourceModel.CanCreateOrUpdateEnvironment = types.BoolPointerValue(model.CanCreateOrUpdateEnvironment)
	resourceModel.CanDeleteEnvironment = types.BoolPointerValue(model.CanDeleteEnvironment)
	resourceModel.CanCreateOrUpdateSetting = types.BoolPointerValue(model.CanCreateOrUpdateSetting)
	resourceModel.CanTagSetting = types.BoolPointerValue(model.CanTagSetting)
	resourceModel.CanDeleteSetting = types.BoolPointerValue(model.CanDeleteSetting)
	resourceModel.CanCreateOrUpdateTag = types.BoolPointerValue(model.CanCreateOrUpdateTag)
	resourceModel.CanDeleteTag = types.BoolPointerValue(model.CanDeleteTag)
	resourceModel.CanManageWebhook = types.BoolPointerValue(model.CanManageWebhook)
	resourceModel.CanUseExportImport = types.BoolPointerValue(model.CanUseExportImport)
	resourceModel.CanManageProductPreferences = types.BoolPointerValue(model.CanManageProductPreferences)
	resourceModel.CanManageIntegrations = types.BoolPointerValue(model.CanManageIntegrations)
	resourceModel.CanViewSdkKey = types.BoolPointerValue(model.CanViewSdkKey)
	resourceModel.CanRotateSdkKey = types.BoolPointerValue(model.CanRotateSdkKey)
	resourceModel.CanCreateOrUpdateSegment = types.BoolPointerValue(model.CanCreateOrUpdateSegments)
	resourceModel.CanDeleteSegment = types.BoolPointerValue(model.CanDeleteSegments)
	resourceModel.CanViewProductAuditLogs = types.BoolPointerValue(model.CanViewProductAuditLog)
	resourceModel.CanViewProductStatistics = types.BoolPointerValue(model.CanViewProductStatistics)
	resourceModel.AccessType = types.StringPointerValue((*string)(model.AccessType))
	resourceModel.NewEnvironmentAccessType = types.StringPointerValue((*string)(model.NewEnvironmentAccessType))
	resourceModel.EnvironmentAccess = environmentAccessesMapValue

	return diags
}

func getEnvironmentAccesses(newEnvironmentAccesses map[string]types.String, oldEnvironmentAccesses map[string]types.String, accessType sw.AccessType) (*[]sw.CreateOrUpdateEnvironmentAccessModel, error) {
	elements := make([]sw.CreateOrUpdateEnvironmentAccessModel, 0)

	if accessType != sw.ACCESSTYPE_CUSTOM && len(newEnvironmentAccesses) > 0 {
		return nil, fmt.Errorf("Error: environment_accesses can only be set if the accesstype is custom")
	}

	if accessType != sw.ACCESSTYPE_CUSTOM {
		return &elements, nil
	}

	for environmentIdKey, environmentAccessType := range newEnvironmentAccesses {
		environmentId := environmentIdKey
		environmentAccessType, environmentAccessTypeParseError := sw.NewEnvironmentAccessTypeFromValue(environmentAccessType.ValueString())
		if environmentAccessTypeParseError != nil || *environmentAccessType == sw.ENVIRONMENTACCESSTYPE_NONE {
			return nil, fmt.Errorf("Error: invalid value '" + (string)(*environmentAccessType) + "' for EnvironmentAccessType: valid values are [full readOnly]")
		}
		element := sw.CreateOrUpdateEnvironmentAccessModel{
			EnvironmentId:         &environmentId,
			EnvironmentAccessType: environmentAccessType,
		}

		elements = append(elements, element)
	}

	// We should set none to those environment accesses that were deleted
	for environmentIdKey := range oldEnvironmentAccesses {
		environmentId := environmentIdKey
		_, ok := newEnvironmentAccesses[environmentId]
		if !ok {
			element := sw.CreateOrUpdateEnvironmentAccessModel{
				EnvironmentId:         &environmentId,
				EnvironmentAccessType: sw.ENVIRONMENTACCESSTYPE_NONE.Ptr(),
			}
			elements = append(elements, element)
		}
	}

	return &elements, nil
}
