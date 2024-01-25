package configcat

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var (
	_ datasource.DataSource              = &permissionGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &permissionGroupDataSource{}
)

func NewPermissionGroupDataSource() datasource.DataSource {
	return &permissionGroupDataSource{}
}

type permissionGroupDataSource struct {
	client *client.Client
}

type permissionGroupDataModel struct {
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

type permissionGroupDataSourceModel struct {
	ID              types.String               `tfsdk:"id"`
	ProductId       types.String               `tfsdk:"product_id"`
	NameFilterRegex types.String               `tfsdk:"name_filter_regex"`
	Data            []permissionGroupDataModel `tfsdk:"permission_groups"`
}

func (d *permissionGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission_groups"
}

func (d *permissionGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to access information about existing **" + PermissionGroupResourceName + "s**. [What is a " + PermissionGroupResourceName + " in ConfigCat?](https://configcat.com/docs/advanced/team-management/team-management-basics/#permissions--permission-groups-product-level)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "Internal ID of the data source. Do not use.",
				Computed:    true,
			},
			ProductId: schema.StringAttribute{
				Description: "The ID of the " + ProductResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
			},
			NameFilterRegex: schema.StringAttribute{
				Description: "Filter the " + PermissionGroupResourceName + "s by name.",
				Optional:    true,
				Validators:  []validator.String{IsRegex()},
			},
			PermissionGroups: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						PermissionGroupId: schema.StringAttribute{
							Description: "The unique " + PermissionGroupResourceName + " ID.",
							Computed:    true,
						},
						Name: schema.StringAttribute{
							Description: "The name of the " + PermissionGroupResourceName + ".",
							Computed:    true,
						},

						PermissionGroupCanManageMembers: schema.BoolAttribute{
							Description: "Group members can manage team members.",
							Computed:    true,
						},
						PermissionGroupCanCreateOrUpdateConfig: schema.BoolAttribute{
							Description: "Group members can create/update Configs.",
							Computed:    true,
						},
						PermissionGroupCanDeleteConfig: schema.BoolAttribute{
							Description: "Group members can delete Configs.",
							Computed:    true,
						},
						PermissionGroupCanCreateOrUpdateEnvironment: schema.BoolAttribute{
							Description: "Group members can create/update Environments.",
							Computed:    true,
						},
						PermissionGroupCanDeleteEnvironment: schema.BoolAttribute{
							Description: "Group members can delete Environments.",
							Computed:    true,
						},
						PermissionGroupCanCreateOrUpdateSetting: schema.BoolAttribute{
							Description: "Group members can create/update Feature Flags and Settings.",
							Computed:    true,
						},
						PermissionGroupCanTagSetting: schema.BoolAttribute{
							Description: "Group members can attach/detach Tags to Feature Flags and Settings.",
							Computed:    true,
						},
						PermissionGroupCanDeleteSetting: schema.BoolAttribute{
							Description: "Group members can delete Feature Flags and Settings.",
							Computed:    true,
						},
						PermissionGroupCanCreateOrUpdateTag: schema.BoolAttribute{
							Description: "Group members can create/update Tags.",
							Computed:    true,
						},
						PermissionGroupCanDeleteTag: schema.BoolAttribute{
							Description: "Group members can delete Tags.",
							Computed:    true,
						},
						PermissionGroupCanManageWebhook: schema.BoolAttribute{
							Description: "Group members can create/update/delete Webhooks.",
							Computed:    true,
						},
						PermissionGroupCanUseExportImport: schema.BoolAttribute{
							Description: "Group members can use the export/import feature.",
							Computed:    true,
						},
						PermissionGroupCanManageProductPreferences: schema.BoolAttribute{
							Description: "Group members can update Product preferences.",
							Computed:    true,
						},
						PermissionGroupCanManageIntegrations: schema.BoolAttribute{
							Description: "Group members can add and configure integrations.",
							Computed:    true,
						},
						PermissionGroupCanViewSdkKey: schema.BoolAttribute{
							Description: "Group members has access to SDK keys.",
							Computed:    true,
						},
						PermissionGroupCanRotateSdkKey: schema.BoolAttribute{
							Description: "Group members can rotate SDK keys.",
							Computed:    true,
						},
						PermissionGroupCanCreateOrUpdateSegment: schema.BoolAttribute{
							Description: "Group members can create/update Segments.",
							Computed:    true,
						},
						PermissionGroupCanDeleteSegment: schema.BoolAttribute{
							Description: "Group members can delete Segments.",
							Computed:    true,
						},
						PermissionGroupCanViewProductAuditLogs: schema.BoolAttribute{
							Description: "Group members has access to audit logs.",
							Computed:    true,
						},
						PermissionGroupCanViewProductStatistics: schema.BoolAttribute{
							Description: "Group members has access to product statistics.",
							Computed:    true,
						},
						PermissionGroupAccessType: schema.StringAttribute{
							Description: "Represent the Feature Management permission. Possible values: readOnly, full, custom",
							Computed:    true,
						},
						PermissionGroupNewEnvironmentAccessType: schema.StringAttribute{
							Description: "Represent the environment specific Feature Management permission for new Environments. Possible values: full, readOnly, none",
							Computed:    true,
						},
						PermissionGroupEnvironmentAccess: schema.MapAttribute{
							Description: "The environment specific permissions map block. Keys are the Environment IDs and the values represent the environment specific Feature Management permission. Possible values: full, readOnly",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *permissionGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *permissionGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state permissionGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resources, err := d.client.GetPermissionGroups(state.ProductId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+PermissionGroupResourceName+" data, got error: %s", err))
		return
	}

	state.ID = types.StringValue(state.ProductId.ValueString() + strconv.FormatInt(time.Now().Unix(), 10))

	filteredResources := []sw.PermissionGroupModel{}
	if !state.NameFilterRegex.IsUnknown() && !state.NameFilterRegex.IsNull() && state.NameFilterRegex.ValueString() != "" {
		regex := regexp.MustCompile(state.NameFilterRegex.ValueString())
		for i := range resources {
			if regex.MatchString(*resources[i].Name.Get()) {
				filteredResources = append(filteredResources, resources[i])
			}
		}
	} else {
		filteredResources = resources
	}

	state.Data = make([]permissionGroupDataModel, len(filteredResources))
	for i, resource := range filteredResources {
		permissionGroupIdString := fmt.Sprintf("%d", *resource.PermissionGroupId)

		environmentAccesses := make(map[string]string, len(resource.EnvironmentAccesses))
		for _, environmentAccess := range resource.EnvironmentAccesses {
			if *environmentAccess.EnvironmentAccessType == sw.ENVIRONMENTACCESSTYPE_NONE {
				continue
			}

			environmentAccesses[*environmentAccess.EnvironmentId] = (string)(*environmentAccess.EnvironmentAccessType)
		}

		environmentAccessesMapValue, diags := types.MapValueFrom(ctx, types.StringType, environmentAccesses)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		dataModel := &permissionGroupDataModel{
			ID:                           types.StringValue(permissionGroupIdString),
			Name:                         types.StringPointerValue(resource.Name.Get()),
			CanManageMembers:             types.BoolPointerValue(resource.CanManageMembers),
			CanCreateOrUpdateConfig:      types.BoolPointerValue(resource.CanCreateOrUpdateConfig),
			CanDeleteConfig:              types.BoolPointerValue(resource.CanDeleteConfig),
			CanCreateOrUpdateEnvironment: types.BoolPointerValue(resource.CanCreateOrUpdateEnvironment),
			CanDeleteEnvironment:         types.BoolPointerValue(resource.CanDeleteEnvironment),
			CanCreateOrUpdateSetting:     types.BoolPointerValue(resource.CanCreateOrUpdateSetting),
			CanTagSetting:                types.BoolPointerValue(resource.CanTagSetting),
			CanDeleteSetting:             types.BoolPointerValue(resource.CanDeleteSetting),
			CanCreateOrUpdateTag:         types.BoolPointerValue(resource.CanCreateOrUpdateTag),
			CanDeleteTag:                 types.BoolPointerValue(resource.CanDeleteTag),
			CanManageWebhook:             types.BoolPointerValue(resource.CanManageWebhook),
			CanUseExportImport:           types.BoolPointerValue(resource.CanUseExportImport),
			CanManageProductPreferences:  types.BoolPointerValue(resource.CanManageProductPreferences),
			CanManageIntegrations:        types.BoolPointerValue(resource.CanManageIntegrations),
			CanViewSdkKey:                types.BoolPointerValue(resource.CanViewSdkKey),
			CanRotateSdkKey:              types.BoolPointerValue(resource.CanRotateSdkKey),
			CanCreateOrUpdateSegment:     types.BoolPointerValue(resource.CanCreateOrUpdateSegments),
			CanDeleteSegment:             types.BoolPointerValue(resource.CanDeleteSegments),
			CanViewProductAuditLogs:      types.BoolPointerValue(resource.CanViewProductAuditLog),
			CanViewProductStatistics:     types.BoolPointerValue(resource.CanViewProductStatistics),
			AccessType:                   types.StringPointerValue((*string)(resource.AccessType)),
			NewEnvironmentAccessType:     types.StringPointerValue((*string)(resource.NewEnvironmentAccessType)),
			EnvironmentAccess:            environmentAccessesMapValue,
		}

		state.Data[i] = *dataModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
