package configcat

import (
	"context"
	"fmt"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sw "github.com/configcat/configcat-publicapi-go-client/v2"
)

var _ resource.Resource = &integrationResource{}
var _ resource.ResourceWithImportState = &integrationResource{}

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

type integrationResource struct {
	client *client.Client
}

type integrationResourceModel struct {
	ProductId types.String `tfsdk:"product_id"`

	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	IntegrationType types.String `tfsdk:"integration_type"`
	Parameters      types.Map    `tfsdk:"parameters"`
	Configs         types.Set    `tfsdk:"configs"`
	Environments    types.Set    `tfsdk:"environments"`
}

func (r *integrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *integrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages an **" + IntegrationResourceName + "**. [Read more about " + IntegrationResourceName + "s in ConfigCat](https://configcat.com/docs/integrations/overview/)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + IntegrationResourceName + ".",
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
				Description: "The name of the " + IntegrationResourceName + ".",
				Required:    true,
			},

			IntegrationType: schema.StringAttribute{
				Description: "The integration type of the " + IntegrationResourceName + ". Possible values: `dataDog`, `slack`, `amplitude`, `mixPanel`, `segment`, `pubNub`.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			IntegrationParameters: schema.MapAttribute{
				MarkdownDescription: "Parameters of the integration. The Parameters dictionary differs for each IntegrationType. See available options per integration type at the examples.",
				Computed:            true,
				Optional:            true,
				ElementType:         types.StringType,
			},

			IntegrationConfigs: schema.SetAttribute{
				Description: "List of Config IDs that are connected with this Integration. If the list is empty, all of the Configs are connected.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},

			IntegrationEnvironments: schema.SetAttribute{
				Description: "List of Environment IDs that are connected with this Integration. If the list is empty, all of the Environments are connected.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *integrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan integrationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationType, integrationTypeParseErr := sw.NewIntegrationTypeFromValue(plan.IntegrationType.ValueString())
	if integrationTypeParseErr != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Could not parse %s: %s. Error: %s", *integrationType, plan.IntegrationType.ValueString(), integrationTypeParseErr))
		return
	}

	parameters, diags := parseIntegrationParameters(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateParameters(*integrationType, parameters)...)

	configs, diags := parseIntegrationConfigs(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	environments, diags := parseIntegrationEnvironments(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sw.CreateIntegrationModel{
		Name:            plan.Name.ValueString(),
		IntegrationType: *integrationType,
		Parameters:      parameters,
		ConfigIds:       configs,
		EnvironmentIds:  environments,
	}

	model, err := r.client.CreateIntegration(plan.ProductId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+IntegrationResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(ctx, *model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func testParameters(parameters map[string]string, keys []string) diag.Diagnostics {
	var diag diag.Diagnostics
	for _, key := range keys {
		_, ok := parameters[key]
		if !ok {
			diag.AddAttributeError(path.Root(IntegrationParameters), key+" parameter is required", key+" parameter is required")
			return diag
		}
	}

	return diag
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state integrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetIntegration(state.ID.ValueString())
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+IntegrationResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(ctx, *model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan integrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.Equal(state.Name) && plan.IntegrationType.Equal(state.IntegrationType) && plan.Parameters.Equal(state.Parameters) && plan.Configs.Equal(state.Configs) && plan.Environments.Equal(state.Environments) {
		return
	}

	parameters, diags := parseIntegrationParameters(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	configs, diags := parseIntegrationConfigs(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	environments, diags := parseIntegrationEnvironments(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sw.ModifyIntegrationRequest{
		Name:           plan.Name.ValueString(),
		Parameters:     parameters,
		ConfigIds:      configs,
		EnvironmentIds: environments,
	}

	model, err := r.client.UpdateIntegration(plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+IntegrationResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(ctx, *model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIntegration(state.ID.ValueString())

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, IntegrationResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+IntegrationResourceName+", got error: %s", err))
		return
	}
}

func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *integrationResourceModel) UpdateFromApiModel(ctx context.Context, model sw.IntegrationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceModel.ID = types.StringPointerValue(model.IntegrationId)
	resourceModel.ProductId = types.StringPointerValue(model.Product.ProductId)
	resourceModel.Name = types.StringPointerValue(model.Name.Get())
	resourceModel.IntegrationType = types.StringValue(string(*model.IntegrationType))

	parameterMapValue, diags := types.MapValueFrom(ctx, types.StringType, model.Parameters)
	if diags.HasError() {
		return diags
	}
	resourceModel.Parameters = parameterMapValue

	configsValue, diags := types.SetValueFrom(ctx, types.StringType, model.ConfigIds)
	if diags.HasError() {
		return diags
	}
	resourceModel.Configs = configsValue

	environmentsValue, diags := types.SetValueFrom(ctx, types.StringType, model.EnvironmentIds)
	if diags.HasError() {
		return diags
	}
	resourceModel.Environments = environmentsValue

	return diags
}

func parseIntegrationParameters(ctx context.Context, plan integrationResourceModel) (map[string]string, diag.Diagnostics) {
	var result map[string]string
	var diags diag.Diagnostics

	if plan.Parameters.IsUnknown() || plan.Parameters.IsNull() {
		result = make(map[string]string, 0)
	} else {
		result = make(map[string]string, len(plan.Parameters.Elements()))
		diags.Append(plan.Parameters.ElementsAs(ctx, &result, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}
	return result, diags
}

func parseIntegrationConfigs(ctx context.Context, plan integrationResourceModel) ([]string, diag.Diagnostics) {
	var result []string
	var diags diag.Diagnostics

	if plan.Configs.IsUnknown() || plan.Configs.IsNull() {
		result = make([]string, 0)
	} else {
		result = make([]string, len(plan.Configs.Elements()))
		diags.Append(plan.Configs.ElementsAs(ctx, &result, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}
	return result, diags
}

func parseIntegrationEnvironments(ctx context.Context, plan integrationResourceModel) ([]string, diag.Diagnostics) {
	var result []string
	var diags diag.Diagnostics

	if plan.Environments.IsUnknown() || plan.Environments.IsNull() {
		result = make([]string, 0)
	} else {
		result = make([]string, len(plan.Environments.Elements()))
		diags.Append(plan.Environments.ElementsAs(ctx, &result, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}
	return result, diags
}

func validateParameters(integrationType sw.IntegrationType, parameters map[string]string) diag.Diagnostics {
	var diag diag.Diagnostics

	switch integrationType {
	case sw.INTEGRATIONTYPE_DATA_DOG:
		return testParameters(parameters, []string{"apikey"})
	case sw.INTEGRATIONTYPE_SLACK:
		return testParameters(parameters, []string{"incoming_webhook.url"})
	case sw.INTEGRATIONTYPE_AMPLITUDE:
		return testParameters(parameters, []string{"apiKey", "secretKey"})
	case sw.INTEGRATIONTYPE_MIX_PANEL:
		return testParameters(parameters, []string{"serviceAccountUserName", "serviceAccountSecret", "projectId"})
	case sw.INTEGRATIONTYPE_SEGMENT:
		return testParameters(parameters, []string{"writeKey"})
	case sw.INTEGRATIONTYPE_PUB_NUB:
		return testParameters(parameters, []string{"pubKey", "subKey", "channel"})
	}

	return diag
}
