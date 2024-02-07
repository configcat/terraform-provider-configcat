package configcat

import (
	"context"
	"fmt"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var _ resource.Resource = &productResource{}
var _ resource.ResourceWithImportState = &productResource{}

func NewProductResource() resource.Resource {
	return &productResource{}
}

type productResource struct {
	client *client.Client
}

type productResourceModel struct {
	OrganizationId types.String `tfsdk:"organization_id"`

	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Order       types.Int64  `tfsdk:"order"`
}

func (r *productResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

func (r *productResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages a **" + ProductResourceName + "**. [What is a " + ProductResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + ProductResourceName + ".",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			OrganizationId: schema.StringAttribute{
				Description: "The ID of the " + OrganizationResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			Name: schema.StringAttribute{
				Description: "The name of the " + ProductResourceName + ".",
				Required:    true,
			},
			Description: schema.StringAttribute{
				Description: "The description of the " + ProductResourceName + ".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			Order: schema.Int64Attribute{
				Description: "The order of the " + ProductResourceName + " within a " + OrganizationResourceName + " (zero-based). If multiple " + ProductResourceName + "s has the same order, they are displayed in alphabetical order.",
				Required:    true,
			},
		},
	}
}

func (r *productResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *productResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan productResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	order := int32(plan.Order.ValueInt64())
	body := sw.CreateProductRequest{
		Name:        plan.Name.ValueString(),
		Description: *sw.NewNullableString(plan.Description.ValueStringPointer()),
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.CreateProduct(plan.OrganizationId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+ProductResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state productResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetProduct(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+ProductResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *productResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan productResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.Equal(state.Name) && plan.Description.Equal(state.Description) && plan.Order.Equal(state.Order) {
		return
	}

	order := int32(plan.Order.ValueInt64())
	body := sw.UpdateProductRequest{
		Name:        *sw.NewNullableString(plan.Name.ValueStringPointer()),
		Description: *sw.NewNullableString(plan.Description.ValueStringPointer()),
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.UpdateProduct(plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+ProductResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *productResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state productResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProduct(state.ID.ValueString())

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, ProductResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+ProductResourceName+", got error: %s", err))
		return
	}
}

func (r *productResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *productResourceModel) UpdateFromApiModel(model sw.ProductModel) {
	modelOrder := int64(*model.Order)
	resourceModel.ID = types.StringPointerValue(model.ProductId)
	resourceModel.OrganizationId = types.StringPointerValue(model.Organization.OrganizationId)
	resourceModel.Name = types.StringPointerValue(model.Name.Get())
	resourceModel.Description = types.StringPointerValue(model.Description.Get())
	resourceModel.Order = types.Int64Value(modelOrder)
}
