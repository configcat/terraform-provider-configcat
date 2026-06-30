package configcat

import (
	"context"
	"fmt"
	"strconv"

	"github.com/configcat/terraform-provider-configcat/v5/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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

type predefinedVariationValueModel struct {
	BoolValue   types.Bool    `tfsdk:"bool_value"`
	StringValue types.String  `tfsdk:"string_value"`
	IntValue    types.Int32   `tfsdk:"int_value"`
	DoubleValue types.Float64 `tfsdk:"double_value"`
}

type predefinedVariationModel struct {
	PredefinedVariationId types.String                   `tfsdk:"predefined_variation_id"`
	Value                 *predefinedVariationValueModel `tfsdk:"value"`
	Name                  types.String                   `tfsdk:"name"`
	Hint                  types.String                   `tfsdk:"hint"`
}

type settingResourceModel struct {
	ConfigId types.String `tfsdk:"config_id"`

	ID                   types.String                `tfsdk:"id"`
	Key                  types.String                `tfsdk:"key"`
	Name                 types.String                `tfsdk:"name"`
	Hint                 types.String                `tfsdk:"hint"`
	SettingType          types.String                `tfsdk:"setting_type"`
	IsJson               types.Bool                  `tfsdk:"is_json"`
	Order                types.Int64                 `tfsdk:"order"`
	PredefinedVariations *[]predefinedVariationModel `tfsdk:"predefined_variations"`
}

func createPredefinedVariationValueSchema() *schema.SingleNestedAttribute {
	return &schema.SingleNestedAttribute{
		Required:    true,
		Description: "Represents the value of a " + PredefinedVariationResourceName + ".",
		Attributes: map[string]schema.Attribute{
			BoolValue: schema.BoolAttribute{
				Optional:    true,
				Description: "The boolean representation of the value.",
				Validators: []validator.Bool{
					boolvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName(StringValue),
						path.MatchRelative().AtParent().AtName(IntValue),
						path.MatchRelative().AtParent().AtName(DoubleValue),
					),
				},
			},
			StringValue: schema.StringAttribute{
				Optional:    true,
				Description: "The string representation of the value.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName(BoolValue),
						path.MatchRelative().AtParent().AtName(IntValue),
						path.MatchRelative().AtParent().AtName(DoubleValue),
					),
				},
			},
			IntValue: schema.Int32Attribute{
				Optional:    true,
				Description: "The whole number representation of the value.",
				Validators: []validator.Int32{
					int32validator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName(BoolValue),
						path.MatchRelative().AtParent().AtName(StringValue),
						path.MatchRelative().AtParent().AtName(DoubleValue),
					),
				},
			},
			DoubleValue: schema.Float64Attribute{
				Optional:    true,
				Description: "The decimal number representation of the value.",
				Validators: []validator.Float64{
					float64validator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName(BoolValue),
						path.MatchRelative().AtParent().AtName(StringValue),
						path.MatchRelative().AtParent().AtName(IntValue),
					),
				},
			},
		},
	}
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
			SettingIsJson: schema.BoolAttribute{
				Description: "Whether this " + SettingResourceName + " should validate string values as JSON values.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			Order: schema.Int64Attribute{
				Description: "The order of the " + SettingResourceName + " within a " + ProductResourceName + " (zero-based). If multiple " + SettingsResourceName + " has the same order, they are displayed in alphabetical order.",
				Required:    true,
			},

			PredefinedVariations: schema.ListNestedAttribute{
				Optional:    true,
				Description: "The predefined variations of the " + SettingResourceName + ". The feature is currently in closed beta state and cannot be used. ",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						PredefinedVariationId: schema.StringAttribute{
							Description: "The unique ID of the " + PredefinedVariationResourceName + ".",
							Computed:    true,
						},
						PredefinedVariationValue: createPredefinedVariationValueSchema(),
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

	predefinedVariations, predefinedVariationsErr := getPredefinedVariationsForCreate(plan.PredefinedVariations)
	if predefinedVariationsErr != nil {
		resp.Diagnostics.AddAttributeError(path.Root(PredefinedVariations), "invalid predefined variations", predefinedVariationsErr.Error())
		return
	}

	body := sw.CreateSettingInitialValues{
		Key:                  plan.Key.ValueString(),
		Name:                 plan.Name.ValueString(),
		Hint:                 *sw.NewNullableString(plan.Hint.ValueStringPointer()),
		SettingType:          *settingType,
		IsJson:               *sw.NewNullableBool(plan.IsJson.ValueBoolPointer()),
		Order:                *sw.NewNullableInt32(&order),
		PredefinedVariations: predefinedVariations,
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

	settingID, convErr := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Setting ID", convErr.Error())
		return
	}

	if !plan.Name.Equal(state.Name) || !plan.Hint.Equal(state.Hint) || !plan.IsJson.Equal(state.IsJson) || !plan.Order.Equal(state.Order) {
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

		if !plan.IsJson.Equal(state.IsJson) {
			operations = append(operations, sw.JsonPatchOperation{
				Op:    sw.OPERATIONTYPE_REPLACE,
				Path:  "/isJson",
				Value: plan.IsJson.ValueBool(),
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

		updateSettingModel, updateSettingErr := r.client.UpdateSetting(int32(settingID), operations)
		if updateSettingErr != nil {
			resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+SettingResourceName+", got error: %s", updateSettingErr))
			return
		}

		updateSettingModelError := plan.UpdateFromApiModel(*updateSettingModel)
		if updateSettingModelError != nil {
			resp.Diagnostics.AddError("Unable to parse API response", fmt.Sprintf("Unable to parse API response for "+SettingResourceName+", got error: %s", updateSettingModelError))
			return
		}
	}

	if plan.PredefinedVariations != nil && state.PredefinedVariations != nil {

		settingTypeString := plan.SettingType.ValueString()
		settingType, err := sw.NewSettingTypeFromValue(settingTypeString)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root(SettingType), "invalid setting_type", err.Error())
			return
		}

		predefinedVariationsChanged := len(*plan.PredefinedVariations) != len(*state.PredefinedVariations)
		if !predefinedVariationsChanged {
			for index, planPredefinedVariation := range *plan.PredefinedVariations {
				statePredefinedVariation := (*state.PredefinedVariations)[index]
				predefinedVariationsChanged = !planPredefinedVariation.Name.Equal(statePredefinedVariation.Name) || !planPredefinedVariation.Hint.Equal(statePredefinedVariation.Hint)
				if predefinedVariationsChanged {
					break
				}
				switch *settingType {
				case sw.SETTINGTYPE_BOOLEAN:
					predefinedVariationsChanged = planPredefinedVariation.Value.BoolValue.Equal(statePredefinedVariation.Value.BoolValue)
					break
				case sw.SETTINGTYPE_STRING:
					predefinedVariationsChanged = planPredefinedVariation.Value.StringValue.Equal(statePredefinedVariation.Value.StringValue)
					break
				case sw.SETTINGTYPE_INT:
					predefinedVariationsChanged = planPredefinedVariation.Value.IntValue.Equal(statePredefinedVariation.Value.IntValue)
					break
				case sw.SETTINGTYPE_DOUBLE:
					predefinedVariationsChanged = planPredefinedVariation.Value.DoubleValue.Equal(statePredefinedVariation.Value.DoubleValue)
					break
				default:
					break
				}
				if predefinedVariationsChanged {
					break
				}
			}
		}

		if predefinedVariationsChanged {
			updatePredefinedVariations := make([]sw.UpdatePredefinedVariationModel, len(*plan.PredefinedVariations))

			for planIndex, planPredefinedVariation := range *plan.PredefinedVariations {
				var predefinedVariationId *string

				for _, statePredefinedVariation := range *state.PredefinedVariations {
					if (*settingType == sw.SETTINGTYPE_BOOLEAN && statePredefinedVariation.Value.BoolValue.Equal(planPredefinedVariation.Value.BoolValue)) ||
						(*settingType == sw.SETTINGTYPE_STRING && statePredefinedVariation.Value.StringValue.Equal(planPredefinedVariation.Value.StringValue)) ||
						(*settingType == sw.SETTINGTYPE_INT && statePredefinedVariation.Value.IntValue.Equal(planPredefinedVariation.Value.IntValue)) ||
						(*settingType == sw.SETTINGTYPE_DOUBLE && statePredefinedVariation.Value.DoubleValue.Equal(planPredefinedVariation.Value.DoubleValue)) {
						predefinedVariationId = statePredefinedVariation.PredefinedVariationId.ValueStringPointer()
						break
					}
				}

				updatePredefinedVariation := sw.UpdatePredefinedVariationModel{
					PredefinedVariationId: *sw.NewNullableString(predefinedVariationId),
					Value: sw.UpdatePredefinedVariationValueModel{
						BoolValue:   *sw.NewNullableBool(planPredefinedVariation.Value.BoolValue.ValueBoolPointer()),
						StringValue: *sw.NewNullableString(planPredefinedVariation.Value.StringValue.ValueStringPointer()),
						IntValue:    *sw.NewNullableInt32(planPredefinedVariation.Value.IntValue.ValueInt32Pointer()),
						DoubleValue: *sw.NewNullableFloat64(planPredefinedVariation.Value.DoubleValue.ValueFloat64Pointer()),
					},
					Name: *sw.NewNullableString(planPredefinedVariation.Name.ValueStringPointer()),
					Hint: *sw.NewNullableString(planPredefinedVariation.Hint.ValueStringPointer()),
				}
				updatePredefinedVariations[planIndex] = updatePredefinedVariation
			}

			updatePredefinedVariationsRequest := sw.UpdatePredefinedVariationsRequest{
				PredefinedVariations: updatePredefinedVariations,
			}

			updateVariationsModel, updateVariationsErr := r.client.UpdatePredefinedVariations(int32(settingID), updatePredefinedVariationsRequest)
			if updateVariationsErr != nil {
				resp.Diagnostics.AddError("Unable to Update Resource", fmt.Sprintf("Unable to update "+PredefinedVariations+", got error: %s", updateVariationsErr))
				return
			}

			updatePredefinedVariationsModelErr := plan.UpdatePredefinedVariationsFromApiModel(*settingType, updateVariationsModel.PredefinedVariations)
			if updatePredefinedVariationsModelErr != nil {
				resp.Diagnostics.AddError("Unable to parse API response", fmt.Sprintf("Unable to parse API response for "+PredefinedVariations+", got error: %s", updatePredefinedVariationsModelErr))
				return
			}
		}
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

func getPredefinedVariationValueModel(settingType sw.SettingType, value sw.PredefinedVariationValueModel) (*predefinedVariationValueModel, error) {

	result := predefinedVariationValueModel{}
	switch settingType {
	case sw.SETTINGTYPE_BOOLEAN:
		result.BoolValue = types.BoolPointerValue(value.BoolValue.Get())
		return &result, nil
	case sw.SETTINGTYPE_STRING:
		result.StringValue = types.StringPointerValue(value.StringValue.Get())
		return &result, nil
	case sw.SETTINGTYPE_INT:
		result.IntValue = types.Int32PointerValue(value.IntValue.Get())
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
	resourceModel.IsJson = types.BoolValue(model.IsJson)
	resourceModel.Order = types.Int64Value(modelOrder)

	predefinedVariationsErr := resourceModel.UpdatePredefinedVariationsFromApiModel(model.SettingType, model.PredefinedVariations)
	if predefinedVariationsErr != nil {
		return predefinedVariationsErr
	}

	return nil
}

func (resourceModel *settingResourceModel) UpdatePredefinedVariationsFromApiModel(settingType sw.SettingType, predefinedVariations []sw.PredefinedVariationModel) error {
	if len(predefinedVariations) == 0 {
		resourceModel.PredefinedVariations = nil
		return nil
	}

	predefinedVariationModels := make([]predefinedVariationModel, len(predefinedVariations))
	for index, predefinedVariation := range predefinedVariations {
		value, valueErr := getPredefinedVariationValueModel(settingType, predefinedVariation.Value)
		if valueErr != nil {
			return valueErr
		}
		predefinedVariationModels[index] = predefinedVariationModel{
			PredefinedVariationId: types.StringValue(predefinedVariation.PredefinedVariationId),
			Name:                  types.StringPointerValue(predefinedVariation.Name.Get()),
			Hint:                  types.StringPointerValue(predefinedVariation.Hint.Get()),
			Value:                 value,
		}
	}

	resourceModel.PredefinedVariations = &predefinedVariationModels
	return nil
}

func getPredefinedVariationsForCreate(predefinedVariations *[]predefinedVariationModel) ([]sw.CreatePredefinedVariationModel, error) {
	if predefinedVariations == nil {
		emptyElements := make([]sw.CreatePredefinedVariationModel, 0)
		return emptyElements, nil
	}

	elements := make([]sw.CreatePredefinedVariationModel, len(*predefinedVariations))

	for index, predefinedVariation := range *predefinedVariations {
		elements[index] = sw.CreatePredefinedVariationModel{
			Value: sw.CreatePredefinedVariationValueModel{
				BoolValue:   *sw.NewNullableBool(predefinedVariation.Value.BoolValue.ValueBoolPointer()),
				StringValue: *sw.NewNullableString(predefinedVariation.Value.StringValue.ValueStringPointer()),
				IntValue:    *sw.NewNullableInt32(predefinedVariation.Value.IntValue.ValueInt32Pointer()),
				DoubleValue: *sw.NewNullableFloat64(predefinedVariation.Value.DoubleValue.ValueFloat64Pointer()),
			},
			Name: *sw.NewNullableString(predefinedVariation.Name.ValueStringPointer()),
			Hint: *sw.NewNullableString(predefinedVariation.Hint.ValueStringPointer()),
		}
	}
	return elements, nil
}
