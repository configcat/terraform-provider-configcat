package configcat

import (
	"context"
	"fmt"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var _ resource.Resource = &configResource{}
var _ resource.ResourceWithImportState = &configResource{}

func NewConfigResource() resource.Resource {
	return &configResource{}
}

type configResource struct {
	client *client.Client
}

type configResourceModel struct {
	ProductId types.String `tfsdk:"product_id"`

	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Order             types.Int64  `tfsdk:"order"`
	EvaluationVersion types.String `tfsdk:"evaluation_version"`
}

func (r *configResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *configResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "Creates and manages a **" + ConfigResourceName + "**. [What is a " + ConfigResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)"

	resp.Schema = schema.Schema{
		MarkdownDescription: description,

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + ConfigResourceName + ".",
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
				Description: "The name of the " + ConfigResourceName + ".",
				Required:    true,
			},
			Description: schema.StringAttribute{
				Description: "The description of the " + ConfigResourceName + ".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			EvaluationVersion: schema.StringAttribute{
				Description: "Determines the evaluation version of a Config. Possible values: v1, v2. Default value: v2. Using v2 enables the new features of [Config V2](https://configcat.com/docs/advanced/config-v2).",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(string(sw.EVALUATIONVERSION_V1)),
			},
			Order: schema.Int64Attribute{
				Description: "The order of the " + ConfigResourceName + " within a " + ProductResourceName + " (zero-based). If multiple " + ConfigResourceName + "s has the same order, they are displayed in alphabetical order.",
				Required:    true,
			},
		},
	}
}

func (r *configResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *configResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan configResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	evaluationVersion, evaluationVersionParseErr := sw.NewEvaluationVersionFromValue(plan.EvaluationVersion.ValueString())
	if evaluationVersionParseErr != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Could not parse %s: %s. Error: %s", EvaluationVersion, plan.EvaluationVersion.ValueString(), evaluationVersionParseErr))
		return
	}

	order := int32(plan.Order.ValueInt64())
	body := sw.CreateConfigRequest{
		Name:              plan.Name.ValueString(),
		Description:       *sw.NewNullableString(plan.Description.ValueStringPointer()),
		Order:             *sw.NewNullableInt32(&order),
		EvaluationVersion: evaluationVersion,
	}

	model, err := r.client.CreateConfig(plan.ProductId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+ConfigResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state configResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetConfig(state.ID.ValueString())
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+ConfigResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *configResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan configResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.EvaluationVersion.Equal(state.EvaluationVersion) {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("%s cannot be changed. Please create a new configcat_config resource with the new %s.", EvaluationVersion, EvaluationVersion))
		return
	}

	if plan.Name.Equal(state.Name) && plan.Description.Equal(state.Description) && plan.Order.Equal(state.Order) {
		return
	}

	order := int32(plan.Order.ValueInt64())
	body := sw.UpdateConfigRequest{
		Name:        *sw.NewNullableString(plan.Name.ValueStringPointer()),
		Description: *sw.NewNullableString(plan.Description.ValueStringPointer()),
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.UpdateConfig(plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+ConfigResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *configResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state configResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConfig(state.ID.ValueString())

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, ConfigResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+ConfigResourceName+", got error: %s", err))
		return
	}
}

func (r *configResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *configResourceModel) UpdateFromApiModel(model sw.ConfigModel) {
	modelOrder := int64(*model.Order)
	resourceModel.ID = types.StringPointerValue(model.ConfigId)
	resourceModel.ProductId = types.StringPointerValue(model.Product.ProductId)
	resourceModel.Name = types.StringPointerValue(model.Name.Get())
	resourceModel.Description = types.StringPointerValue(model.Description.Get())
	resourceModel.Order = types.Int64Value(modelOrder)
	resourceModel.EvaluationVersion = types.StringValue(string(*model.EvaluationVersion))
}
