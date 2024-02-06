package configcat

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sw "github.com/configcat/configcat-publicapi-go-client"
)

var _ resource.Resource = &settingValueResource{}
var _ resource.ResourceWithImportState = &settingValueResource{}

func NewSettingValueResource() resource.Resource {
	return &settingValueResource{}
}

type settingValueResource struct {
	client *client.Client
}

type rolloutRuleModel struct {
	ComparisonAttribute types.String `tfsdk:"comparison_attribute"`
	Comparator          types.String `tfsdk:"comparator"`
	ComparisonValue     types.String `tfsdk:"comparison_value"`

	SegmentComparator types.String `tfsdk:"segment_comparator"`
	SegmentId         types.String `tfsdk:"segment_id"`

	Value types.String `tfsdk:"value"`
}

type rolloutPercentageItemModel struct {
	Percentage types.String `tfsdk:"percentage"`
	Value      types.String `tfsdk:"value"`
}

type settingValueResourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	SettingId     types.String `tfsdk:"setting_id"`

	ID             types.String `tfsdk:"id"`
	InitOnly       types.Bool   `tfsdk:"init_only"`
	MandatoryNotes types.String `tfsdk:"mandatory_notes"`
	SettingType    types.String `tfsdk:"setting_type"`

	Value           types.String                 `tfsdk:"value"`
	RolloutRules    []rolloutRuleModel           `tfsdk:"rollout_rules"`
	PercentageItems []rolloutPercentageItemModel `tfsdk:"percentage_items"`
}

func (r *settingValueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_value"
}

func (r *settingValueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Initializes and updates **" + SettingResourceName + "** values. [Read more about the anatomy of a " + SettingResourceName + ".](https://configcat.com/docs/main-concepts) ",

		Attributes: map[string]schema.Attribute{
			ID: schema.StringAttribute{
				Description: "The unique ID of the " + SettingValueResourceName + ".",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			EnvironmentId: schema.StringAttribute{
				Description: "The ID of the " + EnvironmentResourceName + ".",
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			SettingId: schema.StringAttribute{
				Description: "The ID of the " + SettingResourceName + ".",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			InitOnly: schema.BoolAttribute{
				MarkdownDescription: "The main purpose of this resource to provide an initial value for the Feature Flag/Setting.  \n\n" +
					"The `init_only` argument's default value is `true`. Meaning that the Feature Flag or Setting's **value will be only be applied once** during resource creation. If someone modifies the value on the [ConfigCat Dashboard](https://app.configcat.com) those modifications will **not be overwritten** by the Terraform script.\n\n" +
					"If you want to fully manage the Feature Flag/Setting's value from Terraform, set `init_only` argument to `false`. After setting the`init_only` argument to `false` each terraform run will update the Feature Flag/Setting's value to the state provided in Terraform.  \n\n",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			MandatoryNotes: schema.StringAttribute{
				Description: "If the Product's \"Mandatory notes\" preference is turned on for the Environment the Mandatory note must be passed.",
				Optional:    true,
			},
			SettingType: schema.StringAttribute{
				Description: "The type of the " + SettingResourceName + ". Available values: `boolean`|`string`|`int`|`double`.",
				Computed:    true,
			},
			SettingValue: schema.StringAttribute{
				Description: "The " + SettingResourceName + "'s value. Type: `string`. It must be compatible with the `setting_type`.",
				Required:    true,
			},
		},
		Blocks: map[string]schema.Block{
			RolloutRules: schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						RolloutRuleComparisonAttribute: schema.StringAttribute{
							MarkdownDescription: "The [comparison attribute](https://configcat.com/docs/advanced/targeting/#comparison-attribute).",
							Optional:            true,
						},
						RolloutRuleComparator: schema.StringAttribute{
							MarkdownDescription: "The [comparator](https://configcat.com/docs/advanced/targeting/#comparator).",
							Optional:            true,
						},
						RolloutRuleComparisonValue: schema.StringAttribute{
							MarkdownDescription: "The [comparison value](https://configcat.com/docs/advanced/targeting/#comparison-value).",
							Optional:            true,
						},
						RolloutRuleSegmentComparator: schema.StringAttribute{
							Description: "The segment_comparator. Possible values: isIn, isNotIn.",
							Optional:    true,
						},
						RolloutRuleSegmentId: schema.StringAttribute{
							MarkdownDescription: "The [Segment's](https://configcat.com/docs/advanced/segments) unique identifier.",
							Optional:            true,
						},
						RolloutRuleValue: schema.StringAttribute{
							MarkdownDescription: "The exact [value](https://configcat.com/docs/advanced/targeting/#served-value) that will be served to the users who match the targeting rule. Type: `string`. It must be compatible with the `setting_type`.",
							Required:            true,
						},
					},
				},
			},
			RolloutPercentageItems: schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						RolloutPercentageItemPercentage: schema.StringAttribute{
							MarkdownDescription: "Any [number](https://configcat.com/docs/advanced/targeting/#-value) between 0 and 100 that represents a randomly allocated fraction of your users.",
							Required:            true,
						},
						RolloutPercentageItemValue: schema.StringAttribute{
							MarkdownDescription: "The exact [value](https://configcat.com/docs/advanced/targeting/#served-value-1) that will be served to the users that fall into that fraction. Type: `string`. It must be compatible with the `setting_type`.",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *settingValueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingValueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.createOrUpdate(ctx, &req.Plan, nil, &resp.State, &resp.Diagnostics)
}

func (r *settingValueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingValueResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.InitOnly.ValueBool() && !state.ID.IsNull() && !state.ID.IsUnknown() {
		return
	}

	settingID, convErr := strconv.ParseInt(state.SettingId.ValueString(), 10, 64)
	if convErr != nil {
		resp.Diagnostics.AddError("Could not parse Setting ID", convErr.Error())
		return
	}

	model, err := r.client.GetSettingValue(state.EnvironmentId.ValueString(), int32(settingID))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SettingValueResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingValueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.createOrUpdate(ctx, &req.Plan, &req.State, &resp.State, &resp.Diagnostics)
}

func (r *settingValueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *settingValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state settingValueResourceModel

	environmentID, settingID, err := resourceConfigCatSettingValueParseID(req.ID)

	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root(ID), "unexpected ID format", err.Error())
		return
	}

	state.EnvironmentId = types.StringValue(environmentID)
	state.SettingId = types.StringValue(settingID)
	state.InitOnly = types.BoolValue(false)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root(ID), req, resp)
}

func (r *settingValueResource) createOrUpdate(ctx context.Context, requestPlan *tfsdk.Plan, requestState *tfsdk.State, responseState *tfsdk.State, diag *diag.Diagnostics) {
	var plan settingValueResourceModel

	diag.Append(requestPlan.Get(ctx, &plan)...)

	if diag.HasError() {
		return
	}

	if plan.InitOnly.ValueBool() && !plan.ID.IsNull() && !plan.ID.IsUnknown() {
		diag.AddWarning("Changes will be only applied to the state.", "The init_only parameter is set to true so the changes won't be applied in ConfigCat. This mode is only for initializing a feature flag in ConfigCat.")

		if requestState != nil {
			// SettingType is a computed field, we have to set it.
			var state settingValueResourceModel
			diag.Append(requestState.Get(ctx, &state)...)
			if diag.HasError() {
				return
			}

			plan.SettingType = state.SettingType
		}

		diag.Append(responseState.Set(ctx, &plan)...)
		return
	}

	settingID, convErr := strconv.ParseInt(plan.SettingId.ValueString(), 10, 64)
	if convErr != nil {
		diag.AddError("Could not parse Setting ID", convErr.Error())
		return
	}

	settingTypeString := plan.SettingType.ValueString()
	if settingTypeString == "" {
		setting, err := r.client.GetSetting(int32(settingID))
		if err != nil {
			diag.AddAttributeError(path.Root(SettingType), "could not determine setting_type for "+SettingResourceName, err.Error())
			return
		}

		settingTypeString = fmt.Sprintf("%v", *setting.SettingType)
	}

	settingType, settingTypeConvertErr := sw.NewSettingTypeFromValue(settingTypeString)
	if settingTypeConvertErr != nil {
		diag.AddAttributeError(path.Root(SettingType), "could not determine setting_type for "+SettingResourceName, settingTypeConvertErr.Error())
	}

	settingValue, settingValueErr := getSettingValue(settingType, plan.Value.ValueString())

	if settingValueErr != nil {
		diag.AddAttributeError(path.Root(SettingValue), "could not determine value for "+SettingResourceName, settingValueErr.Error())
		return
	}

	rolloutRules, rolloutRulesErr := getRolloutRulesData(&plan.RolloutRules, *settingType)
	if rolloutRulesErr != nil {
		diag.AddAttributeError(path.Root(RolloutRules), "could not parse rollout_rules", rolloutRulesErr.Error())
		return
	}

	rolloutPercentageItems, rolloutPercentageItemsErr := getRolloutPercentageItemsData(&plan.PercentageItems, *settingType)
	if rolloutPercentageItemsErr != nil {
		diag.AddAttributeError(path.Root(RolloutPercentageItems), "could not parse rollout_rules", rolloutPercentageItemsErr.Error())
		return
	}

	body := sw.UpdateSettingValueModel{
		Value:                  settingValue,
		RolloutRules:           *rolloutRules,
		RolloutPercentageItems: *rolloutPercentageItems,
	}

	model, err := r.client.ReplaceSettingValue(plan.EnvironmentId.ValueString(), int32(settingID), body, plan.MandatoryNotes.ValueString())
	if err != nil {
		diag.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+SettingValueResourceName+", got error: %s", err))
		return
	}

	plan.UpdateFromApiModel(*model)

	diag.Append(responseState.Set(ctx, &plan)...)
}

func (resourceModel *settingValueResourceModel) UpdateFromApiModel(model sw.SettingValueModel) {

	resourceModel.ID = types.StringValue(fmt.Sprintf("%s:%d", *model.Environment.EnvironmentId, *model.Setting.SettingId))
	resourceModel.Value = types.StringValue(fmt.Sprintf("%v", model.Value))
	resourceModel.SettingType = types.StringPointerValue((*string)(model.Setting.SettingType))

	resourceModel.RolloutRules = make([]rolloutRuleModel, len(model.RolloutRules))
	for i, rolloutRule := range model.RolloutRules {
		if rolloutRule.Comparator != nil {
			rolloutRuleModel := rolloutRuleModel{
				ComparisonAttribute: types.StringPointerValue(rolloutRule.ComparisonAttribute.Get()),
				Comparator:          types.StringPointerValue((*string)(rolloutRule.Comparator)),
				ComparisonValue:     types.StringPointerValue(rolloutRule.ComparisonValue.Get()),
				Value:               types.StringValue(fmt.Sprintf("%v", rolloutRule.Value)),
			}
			resourceModel.RolloutRules[i] = rolloutRuleModel
		} else if rolloutRule.SegmentComparator != nil {
			{
				rolloutRuleModel := rolloutRuleModel{
					SegmentId:         types.StringPointerValue(rolloutRule.SegmentId.Get()),
					SegmentComparator: types.StringPointerValue((*string)(rolloutRule.SegmentComparator)),
					Value:             types.StringValue(fmt.Sprintf("%v", rolloutRule.Value)),
				}
				resourceModel.RolloutRules[i] = rolloutRuleModel
			}
		}
	}

	resourceModel.PercentageItems = make([]rolloutPercentageItemModel, len(model.RolloutPercentageItems))
	for i, rolloutPercentageItem := range model.RolloutPercentageItems {
		rolloutPercentageItemModel := rolloutPercentageItemModel{
			Percentage: types.StringValue(strconv.FormatInt(rolloutPercentageItem.Percentage, 10)),
			Value:      types.StringValue(fmt.Sprintf("%v", rolloutPercentageItem.Value)),
		}
		resourceModel.PercentageItems[i] = rolloutPercentageItemModel
	}
}

func getSettingValue(settingType *sw.SettingType, value string) (interface{}, error) {

	switch *settingType {
	case sw.SETTINGTYPE_BOOLEAN:
		b, err := strconv.ParseBool(value)
		return b, err
	case sw.SETTINGTYPE_STRING:
		return value, nil
	case sw.SETTINGTYPE_INT:
		i, err := strconv.ParseInt(value, 10, 32)
		if err == nil {
			return int32(i), nil
		}
		return nil, err
	case sw.SETTINGTYPE_DOUBLE:
		f, err := strconv.ParseFloat(value, 64)
		return f, err
	default:
		return nil, fmt.Errorf("could not parse SettingType and Value: %s, %s", *settingType, value)
	}
}

func resourceConfigCatSettingValueParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected environmentID:settingID", id)
	}

	_, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected environmentID:settingID. Error: %s", id, err)
	}

	return parts[0], parts[1], nil
}

func getRolloutRulesData(rolloutRules *[]rolloutRuleModel, settingType sw.SettingType) (*[]sw.RolloutRuleModel, error) {
	if rolloutRules == nil {
		empty := make([]sw.RolloutRuleModel, 0)
		return &empty, nil
	}

	elements := make([]sw.RolloutRuleModel, len(*rolloutRules))
	for i, rolloutRule := range *rolloutRules {
		value, err := getSettingValue(&settingType, rolloutRule.Value.ValueString())
		if err != nil {
			return nil, err
		}

		if !rolloutRule.Comparator.IsUnknown() && !rolloutRule.Comparator.IsNull() {
			if rolloutRule.ComparisonAttribute.IsUnknown() || rolloutRule.ComparisonAttribute.IsNull() {
				return nil, fmt.Errorf("the %s field is required", RolloutRuleComparisonAttribute)
			}
			if rolloutRule.ComparisonValue.IsUnknown() || rolloutRule.ComparisonValue.IsNull() {
				return nil, fmt.Errorf("the %s field is required", RolloutRuleComparisonValue)
			}

			comparator, compErr := sw.NewRolloutRuleComparatorFromValue(rolloutRule.Comparator.ValueString())
			if compErr != nil {
				return nil, compErr
			}

			element := sw.RolloutRuleModel{
				ComparisonAttribute: *sw.NewNullableString(rolloutRule.ComparisonAttribute.ValueStringPointer()),
				Comparator:          comparator,
				ComparisonValue:     *sw.NewNullableString(rolloutRule.ComparisonValue.ValueStringPointer()),
				Value:               &value,
			}

			elements[i] = element
		} else if len(rolloutRule.SegmentComparator.ValueString()) > 0 {
			if rolloutRule.SegmentId.IsUnknown() || rolloutRule.SegmentId.IsNull() {
				return nil, fmt.Errorf("the %s field is required", RolloutRuleSegmentId)
			}

			segmentComparator, compErr := sw.NewSegmentComparatorFromValue(rolloutRule.SegmentComparator.ValueString())
			if compErr != nil {
				return nil, compErr
			}

			element := sw.RolloutRuleModel{
				SegmentComparator: segmentComparator,
				SegmentId:         *sw.NewNullableString(rolloutRule.SegmentId.ValueStringPointer()),
				Value:             &value,
			}

			elements[i] = element
		} else {
			return nil, fmt.Errorf("either the comparator or the segment_comparator should be set")
		}
	}
	return &elements, nil
}

func getRolloutPercentageItemsData(rolloutPercentageItems *[]rolloutPercentageItemModel, settingType sw.SettingType) (*[]sw.RolloutPercentageItemModel, error) {
	if rolloutPercentageItems == nil {
		empty := make([]sw.RolloutPercentageItemModel, 0)
		return &empty, nil
	}

	elements := make([]sw.RolloutPercentageItemModel, len(*rolloutPercentageItems))
	for i, rolloutPercentageItem := range *rolloutPercentageItems {
		value, err := getSettingValue(&settingType, rolloutPercentageItem.Value.ValueString())
		if err != nil {
			return nil, err
		}

		percentage, percErr := strconv.ParseInt(rolloutPercentageItem.Percentage.ValueString(), 10, 32)
		if percErr != nil {
			return nil, percErr
		}

		element := sw.RolloutPercentageItemModel{
			Percentage: percentage,
			Value:      &value,
		}
		elements[i] = element
	}
	return &elements, nil
}
