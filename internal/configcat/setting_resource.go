package configcat

import (
	"context"
	"fmt"
	"strconv"

	"github.com/configcat/terraform-provider-configcat/v5/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sw "github.com/configcat/configcat-publicapi-go-client/v3"
)

var _ resource.Resource = &settingResource{}
var _ resource.ResourceWithImportState = &settingResource{}

func NewSettingResource() resource.Resource {
	return &settingResource{}
}

type settingResource struct {
	client *client.Client
}

type predefinedVariationModel struct {
	PredefinedVariationId types.String       `tfsdk:"predefined_variation_id"`
	Value                 *settingValueModel `tfsdk:"value"`
	Name                  types.String       `tfsdk:"name"`
	Hint                  types.String       `tfsdk:"hint"`
}

type settingResourceModel struct {
	ConfigId types.String `tfsdk:"config_id"`

	ID                   types.String               `tfsdk:"id"`
	Key                  types.String               `tfsdk:"key"`
	Name                 types.String               `tfsdk:"name"`
	Hint                 types.String               `tfsdk:"hint"`
	SettingType          types.String               `tfsdk:"setting_type"`
	Order                types.Int64                `tfsdk:"order"`
	PredefinedVariations []predefinedVariationModel `tfsdk:"predefined_variations"`
}

func (r *settingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting"
}

func (r *settingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates and manages a **" + SettingResourceName + "**. [What is a " + SettingResourceName + " in ConfigCat?](https://configcat.com/docs/main-concepts)",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + SettingResourceName + ".",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			ConfigId: schema.StringAttribute{
				Description: "The ID of the " + ConfigResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			SettingKey: schema.StringAttribute{
				Description: "The key of the " + SettingResourceName + ".",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			Name: schema.StringAttribute{
				Description: "The name of the " + SettingResourceName + ".",
				Required:    true,
			},
			SettingHint: schema.StringAttribute{
				Description: "The hint of the " + SettingResourceName + ".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			SettingType: schema.StringAttribute{
				Description: "The type of the " + SettingResourceName + ". Available values: `boolean`|`string`|`int`|`double`. Default: `boolean`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("boolean"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			Order: schema.Int64Attribute{
				Description: "The order of the " + SettingResourceName + " within a " + ProductResourceName + " (zero-based). If multiple " + SettingsResourceName + " has the same order, they are displayed in alphabetical order.",
				Required:    true,
			},

			PredefinedVariations: schema.ListNestedAttribute{
				Optional:    true,
				Description: "The predefined variations of the " + SettingResourceName,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						PredefinedVariationId: schema.StringAttribute{
							Description: "The unique ID of the " + PredefinedVariationResourceName + ".",
							Computed:    true,
						},
						PredefinedVariationValue: createSettingValueSchema(true, nil),
						PredefinedVariationName: schema.StringAttribute{
							Description: "The name of the " + PredefinedVariationResourceName + ".",
							Optional:    true,
						},
						PredefinedVariationHint: schema.StringAttribute{
							Description: "The hint of the " + PredefinedVariationResourceName + ".",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *settingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	settingTypeString := plan.SettingType.ValueString()
	settingType, err := sw.NewSettingTypeFromValue(settingTypeString)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root(SettingType), "invalid setting_type", err.Error())
		return
	}

	order := int32(plan.Order.ValueInt64())
	body := sw.CreateSettingInitialValues{
		Key:         plan.Key.ValueString(),
		Name:        plan.Name.ValueString(),
		Hint:        *sw.NewNullableString(plan.Hint.ValueStringPointer()),
		SettingType: *settingType,
		Order:       *sw.NewNullableInt32(&order),
	}

	model, err := r.client.CreateSetting(plan.ConfigId.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+SettingResourceName+", got error: %s", err))
		return
	}

	createError := plan.UpdateFromApiModel(*model)
	if createError != nil {
		resp.Diagnostics.AddError("Unable to parse API response", fmt.Sprintf("Unable to parse API response for "+SettingResourceName+", got error: %s", createError))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	settingID, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Setting ID", convErr.Error())
		return
	}

	model, err := r.client.GetSetting(int32(settingID))
	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we have to remove it from the state.
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SettingResourceName+", got error: %s", err))
		return
	}

	readError := state.UpdateFromApiModel(*model)
	if readError != nil {
		resp.Diagnostics.AddError("Unable to parse API response", fmt.Sprintf("Unable to parse API response for "+SettingResourceName+", got error: %s", readError))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan settingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.Equal(state.Name) && plan.Hint.Equal(state.Hint) && plan.Order.Equal(state.Order) {
		return
	}

	settingID, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Setting ID", convErr.Error())
		return
	}

	operations := []sw.JsonPatchOperation{}
	if !plan.Name.Equal(state.Name) {
		operations = append(operations, sw.JsonPatchOperation{
			Op:    sw.OPERATIONTYPE_REPLACE,
			Path:  "/name",
			Value: plan.Name.ValueString(),
		})
	}

	if !plan.Hint.Equal(state.Hint) {
		operations = append(operations, sw.JsonPatchOperation{
			Op:    sw.OPERATIONTYPE_REPLACE,
			Path:  "/hint",
			Value: plan.Hint.ValueString(),
		})
	}

	if !plan.Order.Equal(state.Order) {
		order := int32(plan.Order.ValueInt64())
		operations = append(operations, sw.JsonPatchOperation{
			Op:    sw.OPERATIONTYPE_REPLACE,
			Path:  "/order",
			Value: order,
		})
	}

	model, err := r.client.UpdateSetting(int32(settingID), operations)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+SettingResourceName+", got error: %s", err))
		return
	}

	updateError := plan.UpdateFromApiModel(*model)
	if updateError != nil {
		diag.AddError("Unable to parse API response", fmt.Sprintf("Unable to parse API response for "+SettingResourceName+", got error: %s", updateError))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state settingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	settingID, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Setting ID", convErr.Error())
		return
	}

	err := r.client.DeleteSetting(int32(settingID))

	if err != nil {
		if _, ok := err.(client.NotFoundError); ok {
			// If the resource is already deleted, we can safely remove it from the state.
			tflog.Trace(ctx, SettingResourceName+" is already deleted in ConfigCat, removing from state.")
			return
		}
		resp.Diagnostics.AddError("Unable to Delete Resource", fmt.Sprintf("Unable to delete "+SettingResourceName+", got error: %s", err))
		return
	}
}

func (r *settingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func getPredefinedVariationValueModel(settingType sw.SettingType, value sw.PredefinedVariationValueModel) (*settingValueModel, error) {

	result := settingValueModel{}
	switch settingType {
	case sw.SETTINGTYPE_BOOLEAN:
		result.BoolValue = types.BoolPointerValue(value.BoolValue.Get())
		return &result, nil
	case sw.SETTINGTYPE_STRING:
		result.StringValue = types.StringPointerValue(value.StringValue.Get())
		return &result, nil
	case sw.SETTINGTYPE_INT:
		int64Value := int64(*value.IntValue.Get())
		result.IntValue = types.Int64PointerValue(&int64Value)
		return &result, nil
	case sw.SETTINGTYPE_DOUBLE:
		result.DoubleValue = types.Float64PointerValue(value.DoubleValue.Get())
		return &result, nil
	default:
		return nil, fmt.Errorf("could not parse SettingType: %s", settingType)
	}
}

func (resourceModel *settingResourceModel) UpdateFromApiModel(model sw.SettingModel) error {
	modelOrder := int64(model.Order)
	resourceModel.ID = types.StringValue(strconv.FormatInt(int64(model.SettingId), 10))
	resourceModel.ConfigId = types.StringValue(model.ConfigId)
	resourceModel.Key = types.StringValue(model.Key)
	resourceModel.Name = types.StringValue(model.Name)
	resourceModel.Hint = types.StringPointerValue(model.Hint.Get())
	resourceModel.SettingType = types.StringValue((string)(model.SettingType))
	resourceModel.Order = types.Int64Value(modelOrder)

	resourceModel.PredefinedVariations = make([]predefinedVariationModel, len(model.PredefinedVariations))
	for index, predefinedVariation := range model.PredefinedVariations {
		value, valueErr := getPredefinedVariationValueModel(model.SettingType, predefinedVariation.Value)
		if valueErr != nil {
			return valueErr
		}
		resourceModel.PredefinedVariations[index] = predefinedVariationModel{
			PredefinedVariationId: types.StringValue(predefinedVariation.PredefinedVariationId),
			Name:                  types.StringPointerValue(predefinedVariation.Name.Get()),
			Hint:                  types.StringPointerValue(predefinedVariation.Hint.Get()),
			Value:                 value,
		}
	}

	return nil
}
