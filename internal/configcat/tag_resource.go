package configcat

import (
	"context"
	"fmt"
	"strconv"

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

var _ resource.Resource = &tagResource{}
var _ resource.ResourceWithImportState = &tagResource{}

func NewTagResource() resource.Resource {
	return &tagResource{}
}

type tagResource struct {
	client *client.Client
}

type tagResourceModel struct {
	ProductId types.String `tfsdk:"product_id"`

	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Color types.String `tfsdk:"color"`
}

func (r *tagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *tagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages a **" + TagResourceName + "**.",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + TagResourceName + ".",
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
				Description: "The name of the " + TagResourceName + ".",
				Required:    true,
			},
			Color: schema.StringAttribute{
				MarkdownDescription: "The color of the " + TagResourceName + ". Default value. `panther`. Valid values: `panther`|`whale`|`salmon`|`lizard`|`canary`|`koala`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("panther"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *tagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := sw.CreateTagModel{
		Name:  plan.Name.ValueString(),
		Color: *sw.NewNullableString(plan.Color.ValueStringPointer()),
	}

	model, err := r.client.CreateTag(plan.ProductId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+TagResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagID, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Tag ID", convErr.Error())
		return
	}

	model, err := r.client.GetTag(tagID)
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+TagResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan tagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.Equal(state.Name) && plan.Color.Equal(state.Color) {
		return
	}

	body := sw.UpdateTagModel{
		Name:  *sw.NewNullableString(plan.Name.ValueStringPointer()),
		Color: *sw.NewNullableString(plan.Color.ValueStringPointer()),
	}

	tagID, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Tag ID", convErr.Error())
		return
	}

	model, err := r.client.UpdateTag(tagID, body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+TagResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tagResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagID, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Tag ID", convErr.Error())
		return
	}

	err := r.client.DeleteTag(tagID)

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, TagResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+TagResourceName+", got error: %s", err))
		return
	}
}

func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *tagResourceModel) UpdateFromApiModel(model sw.TagModel) {
	resourceModel.ID = types.StringValue(strconv.FormatInt(*model.TagId, 10))
	resourceModel.ProductId = types.StringPointerValue(model.Product.ProductId)
	resourceModel.Name = types.StringPointerValue(model.Name.Get())
	resourceModel.Color = types.StringPointerValue(model.Color.Get())
}
