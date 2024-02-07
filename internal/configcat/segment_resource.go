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

var _ resource.Resource = &segmentResource{}
var _ resource.ResourceWithImportState = &segmentResource{}

func NewSegmentResource() resource.Resource {
	return &segmentResource{}
}

type segmentResource struct {
	client *client.Client
}

type segmentResourceModel struct {
	ProductId types.String `tfsdk:"product_id"`

	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	ComparisonAttribute types.String `tfsdk:"comparison_attribute"`
	Comparator          types.String `tfsdk:"comparator"`
	ComparisonValue     types.String `tfsdk:"comparison_value"`
}

func (r *segmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (r *segmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages a **" + SegmentResourceName + "**. [What is a " + SegmentResourceName + " in ConfigCat?](https://configcat.com/docs/advanced/segments)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + SegmentResourceName + ".",
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
				Description: "The name of the " + SegmentResourceName + ".",
				Required:    true,
			},
			Description: schema.StringAttribute{
				Description: "The description of the " + SegmentResourceName + ".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			SegmentComparisonAttribute: schema.StringAttribute{
				MarkdownDescription: "The [comparison attribute](https://configcat.com/docs/advanced/targeting/#attribute) of the " + SegmentResourceName + ".",
				Required:            true,
			},
			SegmentComparator: schema.StringAttribute{
				MarkdownDescription: "The [comparator](https://configcat.com/docs/advanced/targeting/#comparator) of the " + SegmentResourceName + ".",
				Required:            true,
			},
			SegmentComparisonValue: schema.StringAttribute{
				MarkdownDescription: "The [comparison value](https://configcat.com/docs/advanced/targeting/#comparison-value) of the " + SegmentResourceName + ".",
				Required:            true,
			},
		},
	}
}

func (r *segmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *segmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan segmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	comparator, compErr := sw.NewRolloutRuleComparatorFromValue(plan.Comparator.ValueString())
	if compErr != nil {
		resp.Diagnostics.AddAttributeError(path.Root(SegmentComparator), "invalid comparator", compErr.Error())
		return
	}

	body := sw.CreateSegmentModel{
		Name:                plan.Name.ValueString(),
		Description:         *sw.NewNullableString(plan.Description.ValueStringPointer()),
		ComparisonAttribute: plan.ComparisonAttribute.ValueString(),
		Comparator:          *comparator,
		ComparisonValue:     plan.ComparisonValue.ValueString(),
	}

	model, err := r.client.CreateSegment(plan.ProductId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+SegmentResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *segmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state segmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.client.GetSegment(state.ID.ValueString())
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SegmentResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *segmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan segmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.Equal(state.Name) && plan.Description.Equal(state.Description) && plan.ComparisonAttribute.Equal(state.ComparisonAttribute) && plan.Comparator.Equal(state.Comparator) && plan.ComparisonValue.Equal(state.ComparisonValue) {
		return
	}

	comparator, compErr := sw.NewRolloutRuleComparatorFromValue(plan.Comparator.ValueString())
	if compErr != nil {
		resp.Diagnostics.AddAttributeError(path.Root(SegmentComparator), "invalid comparator", compErr.Error())
		return
	}

	body := sw.UpdateSegmentModel{
		Name:                *sw.NewNullableString(plan.Name.ValueStringPointer()),
		Description:         *sw.NewNullableString(plan.Description.ValueStringPointer()),
		ComparisonAttribute: *sw.NewNullableString(plan.ComparisonAttribute.ValueStringPointer()),
		Comparator:          comparator,
		ComparisonValue:     *sw.NewNullableString(plan.ComparisonValue.ValueStringPointer()),
	}

	model, err := r.client.UpdateSegment(plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+SegmentResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *segmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state segmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSegment(state.ID.ValueString())

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, SegmentResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+SegmentResourceName+", got error: %s", err))
		return
	}
}

func (r *segmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (resourceModel *segmentResourceModel) UpdateFromApiModel(model sw.SegmentModel) {
	resourceModel.ID = types.StringPointerValue(model.SegmentId)
	resourceModel.ProductId = types.StringPointerValue(model.Product.ProductId)
	resourceModel.Name = types.StringPointerValue(model.Name.Get())
	resourceModel.Description = types.StringPointerValue(model.Description.Get())
	resourceModel.ComparisonAttribute = types.StringPointerValue(model.ComparisonAttribute.Get())
	resourceModel.Comparator = types.StringPointerValue((*string)(model.Comparator))
	resourceModel.ComparisonValue = types.StringPointerValue(model.ComparisonValue.Get())
}
