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

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var _ resource.Resource = &environmentResource{}
var _ resource.ResourceWithImportState = &environmentResource{}

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

type environmentResource struct {
	client *client.Client
}

type environmentResourceModel struct {
	ProductId types.String `tfsdk:"product_id"`

	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	Order       types.Int64  `tfsdk:"order"`
}

func (r *environmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *environmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + EnvironmentResourceName + ".",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			ProductId: schema.StringAttribute{
				Description: "The ID of the Product.",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			Name: schema.StringAttribute{
				Description: "The name of the " + EnvironmentResourceName + ".",
				Required:    true,
			},
			Description: schema.StringAttribute{
				Description: "The description of the " + EnvironmentResourceName + ".",
				Optional:    true,
			},
			Color: schema.StringAttribute{
				Description: "The color of the " + EnvironmentResourceName + ".",
				Optional:    true,
			},
			Order: schema.Int64Attribute{
				Description: "The order of the " + EnvironmentResourceName + " within a Product (zero-based).",
				Required:    true,
			},
		},
	}
}

func (r *environmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data environmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	order := int32(data.Order.ValueInt64())
	body := sw.CreateEnvironmentModel{
		Name:        data.Name.ValueString(),
		Description: *sw.NewNullableString(data.Description.ValueStringPointer()),
		Color:       *sw.NewNullableString(data.Color.ValueStringPointer()),
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.CreateEnvironment(data.ProductId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+EnvironmentResourceName+", got error: %s", err))
		return
	}

	data.UpdateFromApiModel(*model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data environmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetEnvironment(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+EnvironmentResourceName+", got error: %s", err))
		return
	}

	data.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data environmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	order := int32(data.Order.ValueInt64())
	body := sw.UpdateEnvironmentModel{
		Name:        *sw.NewNullableString(data.Name.ValueStringPointer()),
		Description: *sw.NewNullableString(data.Description.ValueStringPointer()),
		Color:       *sw.NewNullableString(data.Color.ValueStringPointer()),
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.UpdateEnvironment(data.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+EnvironmentResourceName+", got error: %s", err))
		return
	}

	data.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data environmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEnvironment(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+EnvironmentResourceName+", got error: %s", err))
		return
	}
}

func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (state *environmentResourceModel) UpdateFromApiModel(model sw.EnvironmentModel) {
	modelOrder := int64(*model.Order)
	state.ID = types.StringPointerValue(model.EnvironmentId)
	state.ProductId = types.StringPointerValue(model.Product.ProductId)
	state.Name = types.StringPointerValue(model.Name.Get())
	state.Description = types.StringPointerValue(model.Description.Get())
	state.Color = types.StringPointerValue(model.Color.Get())
	state.Order = types.Int64Value(modelOrder)
}
