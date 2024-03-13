package configcat

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/configcat/terraform-provider-configcat/internal/configcat/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

func NewSettingValueV2Resource() resource.Resource {
	return &settingValueV2Resource{}
}

type settingValueV2Resource struct {
	client *client.Client
}

type comparisonValueListItemModel struct {
	Value types.String `tfsdk:"value"`
	Hint  types.String `tfsdk:"hint"`
}

type comparisonValueModel struct {
	StringValue types.String                   `tfsdk:"string_value"`
	DoubleValue types.Float64                  `tfsdk:"double_value"`
	ListValue   []comparisonValueListItemModel `tfsdk:"list_values"`
}

type settingValueModel struct {
	BoolValue   types.Bool    `tfsdk:"bool_value"`
	StringValue types.String  `tfsdk:"string_value"`
	IntValue    types.Int64   `tfsdk:"int_value"`
	DoubleValue types.Float64 `tfsdk:"double_value"`
}

type userConditionModel struct {
	ComparisonAttribute types.String          `tfsdk:"comparison_attribute"`
	Comparator          types.String          `tfsdk:"comparator"`
	ComparisonValue     *comparisonValueModel `tfsdk:"comparison_value"`
}

type segmentConditionModel struct {
	SegmentId  types.String `tfsdk:"segment_id"`
	Comparator types.String `tfsdk:"comparator"`
}

type prerequisiteFlagConditionModel struct {
	PrerequisiteSettingId types.String       `tfsdk:"prerequisite_setting_id"`
	Comparator            types.String       `tfsdk:"comparator"`
	ComparisonValue       *settingValueModel `tfsdk:"comparison_value"`
}

type conditionModel struct {
	UserCondition             *userConditionModel             `tfsdk:"user_condition"`
	SegmentCondition          *segmentConditionModel          `tfsdk:"segment_condition"`
	PrerequisiteFlagCondition *prerequisiteFlagConditionModel `tfsdk:"prerequisite_flag_condition"`
}

type percentageOptionModel struct {
	Percentage types.Int64        `tfsdk:"percentage"`
	Value      *settingValueModel `tfsdk:"value"`
}

type targetingRuleModel struct {
	Conditions []conditionModel `tfsdk:"conditions"`

	PercentageOptions []percentageOptionModel `tfsdk:"percentage_options"`
	Value             *settingValueModel      `tfsdk:"value"`
}

type settingValueV2ResourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	SettingId     types.String `tfsdk:"setting_id"`

	ID          types.String `tfsdk:"id"`
	SettingType types.String `tfsdk:"setting_type"`

	InitOnly                      types.Bool   `tfsdk:"init_only"`
	MandatoryNotes                types.String `tfsdk:"mandatory_notes"`
	PercentageEvaluationAttribute types.String `tfsdk:"percentage_evaluation_attribute"`

	DefaultValue   *settingValueModel   `tfsdk:"value"`
	TargetingRules []targetingRuleModel `tfsdk:"targeting_rules"`
}

func (r *settingValueV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting_value_v2"
}

func createSettingValueSchema(required bool, validators []validator.Object) *schema.SingleNestedAttribute {
	return &schema.SingleNestedAttribute{
		Required:    required,
		Optional:    !required,
		Validators:  validators,
		Description: "Represents the value of a " + SettingResourceName + ".",
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
			IntValue: schema.Int64Attribute{
				Optional:    true,
				Description: "The whole number representation of the value.",
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(
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

func (r *settingValueV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	comparisonValueSchema := schema.SingleNestedAttribute{
		Required:    true,
		Description: "The value that the user object's attribute is compared to.",
		Attributes: map[string]schema.Attribute{
			StringValue: schema.StringAttribute{
				Optional:    true,
				Description: "The string representation of the comparison value.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName(DoubleValue),
						path.MatchRelative().AtParent().AtName(ListValues),
					),
				},
			},
			DoubleValue: schema.Float64Attribute{
				Optional:    true,
				Description: "The number representation of the comparison value.",
				Validators: []validator.Float64{
					float64validator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName(StringValue),
						path.MatchRelative().AtParent().AtName(ListValues),
					),
				},
			},
			ListValues: schema.ListNestedAttribute{
				Optional:    true,
				Description: "The list representation of the comparison value.",
				Validators: []validator.List{
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName(StringValue),
						path.MatchRelative().AtParent().AtName(DoubleValue),
					),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ListValueValue: schema.StringAttribute{
							Required:    true,
							Description: "The actual comparison value.",
						},
						ListValueHint: schema.StringAttribute{
							Optional:    true,
							Description: "An optional hint for the comparison value.",
						},
					},
				},
			},
		},
	}

	userConditionSchema := schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Describes a condition that is based on user attributes.",
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName(TargetingRuleSegmentCondition),
				path.MatchRelative().AtParent().AtName(TargetingRulePrerequisiteFlagCondition),
			),
		},
		Attributes: map[string]schema.Attribute{
			TargetingRuleUserConditionComparisonAttribute: schema.StringAttribute{
				Required:    true,
				Description: "The User Object attribute that the condition is based on.",
			},
			TargetingRuleUserConditionComparator: schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The comparison operator which defines the relation between the comparison attribute and the comparison value. For possible values check the [documentation](https://api.configcat.com/docs/index.html#tag/Feature-Flag-and-Setting-values-V2/operation/replace-setting-value-v2).",
			},
			TargetingRuleUserConditionComparisonValue: &comparisonValueSchema,
		},
	}

	segmentConditionSchema := schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Describes a condition that is based on a segment.",
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName(TargetingRuleUserCondition),
				path.MatchRelative().AtParent().AtName(TargetingRulePrerequisiteFlagCondition),
			),
		},
		Attributes: map[string]schema.Attribute{
			TargetingRuleSegmentConditionSegmentId: schema.StringAttribute{
				Required:    true,
				Validators:  []validator.String{IsGuid()},
				Description: "The segment's identifier.",
			},
			TargetingRuleSegmentConditionComparator: schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The segment comparison operator used during the evaluation process. Possible values: `isIn`,`isNotIn`",
			},
		},
	}

	prerequisiteFlagConditionSchema := schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Describes a condition that is based on a prerequisite flag.",
		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				path.MatchRelative().AtParent().AtName(TargetingRuleUserCondition),
				path.MatchRelative().AtParent().AtName(TargetingRuleSegmentCondition),
			),
		},
		Attributes: map[string]schema.Attribute{
			TargetingRulePrerequisiteFlagConditionSettingId: schema.StringAttribute{
				Description: "The prerequisite flag's identifier.",
				Required:    true,
			},
			TargetingRulePrerequisiteFlagConditionComparator: schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Prerequisite flag comparison operator used during the evaluation process. Possible values: `equals`,`doesNotEqual`",
			},
			TargetingRulePrerequisiteFlagConditionComparisonValue: createSettingValueSchema(true, nil),
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Initializes and updates **" + SettingResourceName + "** values for V2 configs. [Read more about the anatomy of a " + SettingResourceName + ".](https://configcat.com/docs/main-concepts) ",

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
			SettingType: schema.StringAttribute{
				Description:   "The type of the " + SettingResourceName + ". Available values: `boolean`|`string`|`int`|`double`.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
			PercentageEvaluationAttribute: schema.StringAttribute{
				Description: "The user attribute used for percentage evaluation. If not set, it defaults to the Identifier user object attribute.",
				Optional:    true,
			},
			DefaultValue: createSettingValueSchema(true, nil),

			TargetingRules: schema.ListNestedAttribute{
				Optional:    true,
				Description: "The targeting rules of the " + SettingResourceName,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						TargetingRuleConditions: schema.ListNestedAttribute{
							Description: "The conditions that are combined with the AND logical operator.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									TargetingRuleUserCondition:             userConditionSchema,
									TargetingRuleSegmentCondition:          segmentConditionSchema,
									TargetingRulePrerequisiteFlagCondition: prerequisiteFlagConditionSchema,
								},
							},
						},

						TargetingRuleValue: createSettingValueSchema(false, []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName(TargetingRulePercentageOptions),
							),
						}),
						TargetingRulePercentageOptions: schema.ListNestedAttribute{
							Optional:    true,
							Description: "The percentage options from where the evaluation process will choose a value based on the flag's percentage evaluation attribute.",
							Validators: []validator.List{
								listvalidator.ExactlyOneOf(
									path.MatchRelative().AtParent().AtName(TargetingRuleValue),
								),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									TargetingRulePercentageOptionPercentage: schema.Int64Attribute{
										Description: "A number between 0 and 100 that represents a randomly allocated fraction of the users.",
										Required:    true,
									},
									TargetingRulePercentageOptionValue: createSettingValueSchema(true, nil),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *settingValueV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *settingValueV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.createOrUpdate(ctx, &req.Plan, nil, &resp.State, &resp.Diagnostics)
}

func (r *settingValueV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingValueV2ResourceModel

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

	model, err := r.client.GetSettingValueV2(state.EnvironmentId.ValueString(), int32(settingID))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read "+SettingValueResourceName+", got error: %s", err))
		return
	}

	state.UpdateFromApiModel(*model)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingValueV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.createOrUpdate(ctx, &req.Plan, &req.State, &resp.State, &resp.Diagnostics)
}

func (r *settingValueV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete operation should not do anything with the Feature flag's values.
}

func (r *settingValueV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	environmentID, settingID, err := resourceConfigCatSettingValueV2ParseID(req.ID)

	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root(ID), "unexpected ID format", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(EnvironmentId), types.StringValue(environmentID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(SettingId), types.StringValue(settingID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(InitOnly), types.BoolValue(true))...)
}

func (r *settingValueV2Resource) createOrUpdate(ctx context.Context, requestPlan *tfsdk.Plan, requestState *tfsdk.State, responseState *tfsdk.State, diag *diag.Diagnostics) {
	var plan settingValueV2ResourceModel
	diag.Append(requestPlan.Get(ctx, &plan)...)

	if diag.HasError() {
		return
	}

	if requestState != nil {
		var state settingValueV2ResourceModel
		diag.Append(requestState.Get(ctx, &state)...)
		if !hasChangesV2(&plan, &state) {
			return
		}
	}

	if plan.InitOnly.ValueBool() && !plan.ID.IsNull() && !plan.ID.IsUnknown() {
		diag.AddWarning("Changes will be only applied to the state.", "The init_only parameter is set to true so the changes won't be applied in ConfigCat. This mode is only for initializing a feature flag in ConfigCat.")
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

	settingValue, settingValueErr := getSettingValueV2(settingType, plan.DefaultValue)

	if settingValueErr != nil {
		diag.AddAttributeError(path.Root(SettingValue), "could not determine value for "+SettingResourceName, settingValueErr.Error())
		return
	}

	targetingRules, targetingRulesErr := getTargetingRulesData(plan.TargetingRules, *settingType)
	if targetingRulesErr != nil {
		diag.AddAttributeError(path.Root(RolloutRules), "could not parse targeting_rules", targetingRulesErr.Error())
		return
	}

	body := sw.UpdateEvaluationFormulaModel{
		DefaultValue:                  *settingValue,
		TargetingRules:                targetingRules,
		PercentageEvaluationAttribute: *sw.NewNullableString(plan.PercentageEvaluationAttribute.ValueStringPointer()),
	}

	model, err := r.client.ReplaceSettingValueV2(plan.EnvironmentId.ValueString(), int32(settingID), body, plan.MandatoryNotes.ValueString())
	if err != nil {
		diag.AddError("Unable to Create Resource", fmt.Sprintf("Unable to create "+SettingValueResourceName+", got error: %s", err))
		return
	}

	updateError := plan.UpdateFromApiModel(*model)
	if updateError != nil {
		diag.AddError("Unable to parse API response", fmt.Sprintf("Unable to parse API response for "+SettingValueResourceName+", got error: %s", updateError))
		return
	}

	diag.Append(responseState.Set(ctx, &plan)...)
}

func getTargetingRulesData(targetingRules []targetingRuleModel, settingType sw.SettingType) ([]sw.TargetingRuleModel, error) {
	if len(targetingRules) == 0 {
		return nil, nil
	}
	result := make([]sw.TargetingRuleModel, len(targetingRules))

	for targetingRuleIndex, targetingRule := range targetingRules {
		targetingRuleModel := &sw.TargetingRuleModel{}

		if len(targetingRule.Conditions) > 0 {
			conditions := make([]sw.ConditionModel, len(targetingRule.Conditions))

			for conditionIndex, condition := range targetingRule.Conditions {

				if condition.UserCondition != nil {
					comparator, comparatorErr := sw.NewUserComparatorFromValue(condition.UserCondition.Comparator.ValueString())
					if comparatorErr != nil {
						return nil, comparatorErr
					}
					comparisonValue, comparisonValueErr := getUserConditionComparisonValueData(condition.UserCondition.ComparisonValue)
					if comparisonValueErr != nil {
						return nil, comparisonValueErr
					}
					conditions[conditionIndex] = sw.ConditionModel{
						UserCondition: &sw.UserConditionModel{
							ComparisonAttribute: condition.UserCondition.ComparisonAttribute.ValueString(),
							Comparator:          *comparator,
							ComparisonValue:     *comparisonValue,
						},
					}
				} else if condition.SegmentCondition != nil {
					comparator, comparatorErr := sw.NewSegmentComparatorFromValue(condition.SegmentCondition.Comparator.ValueString())
					if comparatorErr != nil {
						return nil, comparatorErr
					}
					conditions[conditionIndex] = sw.ConditionModel{
						SegmentCondition: &sw.SegmentConditionModel{
							SegmentId:  condition.SegmentCondition.SegmentId.ValueString(),
							Comparator: *comparator,
						},
					}
				} else if condition.PrerequisiteFlagCondition != nil {
					comparator, comparatorErr := sw.NewPrerequisiteComparatorFromValue(condition.PrerequisiteFlagCondition.Comparator.ValueString())
					if comparatorErr != nil {
						return nil, comparatorErr
					}
					settingID, convErr := strconv.ParseInt(condition.PrerequisiteFlagCondition.PrerequisiteSettingId.ValueString(), 10, 32)
					if convErr != nil {
						return nil, convErr
					}
					prerequisiteComparisonValue, prerequisiteComparisonValueErr := getSettingValueV2WithoutSettingType(*condition.PrerequisiteFlagCondition.ComparisonValue)
					if prerequisiteComparisonValueErr != nil {
						return nil, prerequisiteComparisonValueErr
					}
					conditions[conditionIndex] = sw.ConditionModel{
						PrerequisiteFlagCondition: &sw.PrerequisiteFlagConditionModel{
							PrerequisiteSettingId:       int32(settingID),
							Comparator:                  *comparator,
							PrerequisiteComparisonValue: *prerequisiteComparisonValue,
						},
					}
				} else {
					return nil, fmt.Errorf("exactly one of the %s, %s or %s attributes is required", TargetingRuleUserCondition, TargetingRuleSegmentCondition, TargetingRulePrerequisiteFlagCondition)
				}
			}

			targetingRuleModel.Conditions = conditions
		}

		if len(targetingRule.PercentageOptions) == 0 && targetingRule.Value == nil {
			return nil, fmt.Errorf("the %s or the %s attributes must be provided", TargetingRuleValue, TargetingRulePercentageOptions)
		}

		if len(targetingRule.PercentageOptions) > 0 {
			percentageOptions := make([]sw.PercentageOptionModel, len(targetingRule.PercentageOptions))
			for percentageOptionIndex, percentageOption := range targetingRule.PercentageOptions {
				percentageOptionValue, percentageOptionValueErr := getSettingValueV2(&settingType, percentageOption.Value)
				if percentageOptionValueErr != nil {
					return nil, percentageOptionValueErr
				}

				percentageOptions[percentageOptionIndex] = sw.PercentageOptionModel{
					Percentage: int32(percentageOption.Percentage.ValueInt64()),
					Value:      *percentageOptionValue,
				}
			}

			targetingRuleModel.PercentageOptions = percentageOptions
		}

		if targetingRule.Value != nil {
			targetingRuleValue, targetingRuleValueErr := getSettingValueV2(&settingType, targetingRule.Value)
			if targetingRuleValueErr != nil {
				return nil, targetingRuleValueErr
			}
			targetingRuleModel.Value = targetingRuleValue
		}

		result[targetingRuleIndex] = *targetingRuleModel
	}

	return result, nil
}

func getUserConditionComparisonValueData(comparisonValue *comparisonValueModel) (*sw.ComparisonValueModel, error) {
	if !comparisonValue.StringValue.IsUnknown() && !comparisonValue.StringValue.IsNull() {
		return &sw.ComparisonValueModel{
			StringValue: *sw.NewNullableString(comparisonValue.StringValue.ValueStringPointer()),
		}, nil
	} else if !comparisonValue.DoubleValue.IsUnknown() && !comparisonValue.DoubleValue.IsNull() {
		return &sw.ComparisonValueModel{
			DoubleValue: *sw.NewNullableFloat64(comparisonValue.DoubleValue.ValueFloat64Pointer()),
		}, nil
	} else if len(comparisonValue.ListValue) > 0 {
		listValueItems := make([]sw.ComparisonValueListModel, len(comparisonValue.ListValue))

		for listValueItemIndex, listValueItem := range comparisonValue.ListValue {
			listValueItems[listValueItemIndex] = sw.ComparisonValueListModel{
				Value: listValueItem.Value.ValueString(),
				Hint:  *sw.NewNullableString(listValueItem.Hint.ValueStringPointer()),
			}
		}

		return &sw.ComparisonValueModel{
			ListValue: listValueItems,
		}, nil
	} else {
		return nil, fmt.Errorf("exactly one of the %s, %s or %s attributes is required", StringValue, DoubleValue, ListValues)
	}
}

func (resourceModel *settingValueV2ResourceModel) UpdateFromApiModel(model sw.SettingFormulaModel) error {

	resourceModel.ID = types.StringValue(fmt.Sprintf("%s:%d", *model.Environment.EnvironmentId, *model.Setting.SettingId))
	defaultValue, defaultValueErr := getSettingValueModelV2(model.Setting.SettingType, *model.DefaultValue)
	if defaultValueErr != nil {
		return defaultValueErr
	}

	resourceModel.DefaultValue = defaultValue
	resourceModel.SettingType = types.StringPointerValue((*string)(model.Setting.SettingType))
	resourceModel.PercentageEvaluationAttribute = types.StringPointerValue(model.PercentageEvaluationAttribute.Get())

	if len(model.TargetingRules) > 0 {
		targetingRules := make([]targetingRuleModel, len(model.TargetingRules))

		for targetingRuleIndex, targetingRule := range model.TargetingRules {
			targetingRuleModel := targetingRuleModel{}

			if len(targetingRule.Conditions) > 0 {
				conditions := make([]conditionModel, len(targetingRule.Conditions))

				for conditionIndex, condition := range targetingRule.Conditions {
					if condition.UserCondition != nil {

						comparisonValue := comparisonValueModel{}
						if condition.UserCondition.ComparisonValue.StringValue.IsSet() && condition.UserCondition.ComparisonValue.StringValue.Get() != nil {
							comparisonValue.StringValue = types.StringPointerValue(condition.UserCondition.ComparisonValue.StringValue.Get())
						} else if condition.UserCondition.ComparisonValue.DoubleValue.IsSet() && condition.UserCondition.ComparisonValue.DoubleValue.Get() != nil {
							comparisonValue.DoubleValue = types.Float64PointerValue(condition.UserCondition.ComparisonValue.DoubleValue.Get())
						} else if len(condition.UserCondition.ComparisonValue.ListValue) > 0 {
							listValues := make([]comparisonValueListItemModel, len(condition.UserCondition.ComparisonValue.ListValue))
							for listValueIndex, listValue := range condition.UserCondition.ComparisonValue.ListValue {
								listValues[listValueIndex] = comparisonValueListItemModel{
									Value: types.StringValue(listValue.Value),
									Hint:  types.StringPointerValue(listValue.Hint.Get()),
								}
							}
							comparisonValue.ListValue = listValues
						} else {
							return fmt.Errorf("invalid model. At least StringValue, DoubleValue or ListValue must be provided")
						}

						conditions[conditionIndex] = conditionModel{
							UserCondition: &userConditionModel{
								ComparisonAttribute: types.StringValue(condition.UserCondition.ComparisonAttribute),
								Comparator:          types.StringValue(string(condition.UserCondition.Comparator)),
								ComparisonValue:     &comparisonValue,
							},
						}
					} else if condition.SegmentCondition != nil {
						conditions[conditionIndex] = conditionModel{
							SegmentCondition: &segmentConditionModel{
								SegmentId:  types.StringValue(condition.SegmentCondition.SegmentId),
								Comparator: types.StringValue(string(condition.SegmentCondition.Comparator)),
							},
						}
					} else if condition.PrerequisiteFlagCondition != nil {
						prerequisiteFlagSettingValueModel, prerequisiteFlagSettingValueModelErr := getSettingValueModelV2WithoutSettingType(condition.PrerequisiteFlagCondition.PrerequisiteComparisonValue)
						if prerequisiteFlagSettingValueModelErr != nil {
							return prerequisiteFlagSettingValueModelErr
						}
						conditions[conditionIndex] = conditionModel{
							PrerequisiteFlagCondition: &prerequisiteFlagConditionModel{
								PrerequisiteSettingId: types.StringValue(strconv.FormatInt(int64(condition.PrerequisiteFlagCondition.PrerequisiteSettingId), 10)),
								Comparator:            types.StringValue(string(condition.PrerequisiteFlagCondition.Comparator)),
								ComparisonValue:       prerequisiteFlagSettingValueModel,
							},
						}
					} else {
						return fmt.Errorf("invalid model. At least UserCondition, SegmentCondition or PrerequisiteFlagCondition must be provided")
					}
				}

				targetingRuleModel.Conditions = conditions
			}

			if len(targetingRule.PercentageOptions) == 0 && targetingRule.Value == nil {
				return fmt.Errorf("invalid model. At least PercentageOptions or Value must be provided")
			}

			if len(targetingRule.PercentageOptions) > 0 {
				percentageOptions := make([]percentageOptionModel, len(targetingRule.PercentageOptions))

				for percentageOptionIndex, percentageOption := range targetingRule.PercentageOptions {
					percentageValue, percentageValueErr := getSettingValueModelV2(model.Setting.SettingType, percentageOption.Value)
					if percentageValueErr != nil {
						return percentageValueErr
					}

					percentageOptions[percentageOptionIndex] = percentageOptionModel{
						Percentage: types.Int64Value(int64(percentageOption.Percentage)),
						Value:      percentageValue,
					}
				}

				targetingRuleModel.PercentageOptions = percentageOptions
			}

			if targetingRule.Value != nil {
				targetingRuleValue, targetingRuleValueErr := getSettingValueModelV2(model.Setting.SettingType, *targetingRule.Value)
				if targetingRuleValueErr != nil {
					return targetingRuleValueErr
				}
				targetingRuleModel.Value = targetingRuleValue
			}

			targetingRules[targetingRuleIndex] = targetingRuleModel
		}

		resourceModel.TargetingRules = targetingRules
	}

	return nil
}

func getSettingValueModelV2(settingType *sw.SettingType, value sw.ValueModel) (*settingValueModel, error) {

	result := settingValueModel{}
	switch *settingType {
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
		return nil, fmt.Errorf("could not parse SettingType: %s", *settingType)
	}
}

func getSettingValueModelV2WithoutSettingType(value sw.ValueModel) (*settingValueModel, error) {

	if value.BoolValue.IsSet() && value.BoolValue.Get() != nil {
		return &settingValueModel{BoolValue: types.BoolPointerValue(value.BoolValue.Get())}, nil
	} else if value.StringValue.IsSet() && value.StringValue.Get() != nil {
		return &settingValueModel{StringValue: types.StringPointerValue(value.StringValue.Get())}, nil
	} else if value.IntValue.IsSet() && value.IntValue.Get() != nil {
		int64Value := int64(*value.IntValue.Get())
		return &settingValueModel{IntValue: types.Int64Value(int64Value)}, nil
	} else if value.DoubleValue.IsSet() && value.DoubleValue.Get() != nil {
		return &settingValueModel{DoubleValue: types.Float64PointerValue(value.DoubleValue.Get())}, nil
	} else {
		return nil, fmt.Errorf("invalid model")
	}
}

func getSettingValueV2WithoutSettingType(value settingValueModel) (*sw.ValueModel, error) {

	if !value.BoolValue.IsUnknown() && !value.BoolValue.IsNull() {
		return &sw.ValueModel{BoolValue: *sw.NewNullableBool(value.BoolValue.ValueBoolPointer())}, nil
	} else if !value.StringValue.IsUnknown() && !value.StringValue.IsNull() {
		return &sw.ValueModel{StringValue: *sw.NewNullableString(value.StringValue.ValueStringPointer())}, nil
	} else if !value.IntValue.IsUnknown() && !value.IntValue.IsNull() {
		int64Value := value.IntValue.ValueInt64()
		int32Value := int32(int64Value)
		return &sw.ValueModel{IntValue: *sw.NewNullableInt32(&int32Value)}, nil
	} else if !value.DoubleValue.IsUnknown() && !value.DoubleValue.IsNull() {
		return &sw.ValueModel{DoubleValue: *sw.NewNullableFloat64(value.DoubleValue.ValueFloat64Pointer())}, nil
	} else {
		return nil, fmt.Errorf("exactly one of the %s, %s, %s or %s attributes is required", BoolValue, StringValue, IntValue, DoubleValue)
	}
}

func getSettingValueV2(settingType *sw.SettingType, value *settingValueModel) (*sw.ValueModel, error) {

	result := sw.ValueModel{}
	switch *settingType {
	case sw.SETTINGTYPE_BOOLEAN:
		result.BoolValue = *sw.NewNullableBool(value.BoolValue.ValueBoolPointer())
		return &result, nil
	case sw.SETTINGTYPE_STRING:
		result.StringValue = *sw.NewNullableString(value.StringValue.ValueStringPointer())
		return &result, nil
	case sw.SETTINGTYPE_INT:
		int32Value := int32(value.IntValue.ValueInt64())
		result.IntValue = *sw.NewNullableInt32(&int32Value)
		return &result, nil
	case sw.SETTINGTYPE_DOUBLE:
		result.DoubleValue = *sw.NewNullableFloat64(value.DoubleValue.ValueFloat64Pointer())
		return &result, nil
	default:
		return nil, fmt.Errorf("could not parse SettingType and Value: %s, %s", *settingType, value)
	}
}

func resourceConfigCatSettingValueV2ParseID(id string) (string, string, error) {
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

func hasChangesV2(plan *settingValueV2ResourceModel, state *settingValueV2ResourceModel) bool {
	if !plan.EnvironmentId.Equal(state.EnvironmentId) ||
		!plan.SettingId.Equal(state.SettingId) ||
		!plan.InitOnly.Equal(state.InitOnly) ||
		!plan.MandatoryNotes.Equal(state.MandatoryNotes) ||
		!plan.PercentageEvaluationAttribute.Equal(state.PercentageEvaluationAttribute) ||
		hasSettingValueChanges(plan.DefaultValue, state.DefaultValue) ||

		len(plan.TargetingRules) != len(state.TargetingRules) {
		return true
	}

	for targetingRuleIndex, planTargetingRule := range plan.TargetingRules {
		stateTargetingRule := (state.TargetingRules)[targetingRuleIndex]
		if len(planTargetingRule.Conditions) != len(stateTargetingRule.Conditions) ||
			len(planTargetingRule.PercentageOptions) != len(stateTargetingRule.PercentageOptions) ||
			hasSettingValueChanges(planTargetingRule.Value, stateTargetingRule.Value) {
			return true
		}

		for conditionIndex, planCondition := range planTargetingRule.Conditions {
			stateCondition := (stateTargetingRule.Conditions)[conditionIndex]
			if (planCondition.UserCondition == nil) != (stateCondition.UserCondition == nil) ||
				(planCondition.SegmentCondition == nil) != (stateCondition.SegmentCondition == nil) ||
				(planCondition.PrerequisiteFlagCondition == nil) != (stateCondition.PrerequisiteFlagCondition == nil) {
				return true
			}

			if planCondition.UserCondition != nil {
				if !planCondition.UserCondition.Comparator.Equal(stateCondition.UserCondition.Comparator) ||
					!planCondition.UserCondition.ComparisonAttribute.Equal(stateCondition.UserCondition.ComparisonAttribute) ||
					!planCondition.UserCondition.ComparisonValue.StringValue.Equal(stateCondition.UserCondition.ComparisonValue.StringValue) ||
					!planCondition.UserCondition.ComparisonValue.DoubleValue.Equal(stateCondition.UserCondition.ComparisonValue.DoubleValue) ||
					len(planCondition.UserCondition.ComparisonValue.ListValue) != len(stateCondition.UserCondition.ComparisonValue.ListValue) {
					return true
				}

				for listValueIndex, planListValue := range planCondition.UserCondition.ComparisonValue.ListValue {
					stateListValue := planCondition.UserCondition.ComparisonValue.ListValue[listValueIndex]
					if !planListValue.Value.Equal(stateListValue.Value) ||
						!planListValue.Hint.Equal(stateListValue.Hint) {
						return true
					}
				}
			}

			if planCondition.SegmentCondition != nil {
				if !planCondition.SegmentCondition.SegmentId.Equal(stateCondition.SegmentCondition.SegmentId) ||
					!planCondition.SegmentCondition.Comparator.Equal(stateCondition.SegmentCondition.Comparator) {
					return true
				}
			}

			if planCondition.PrerequisiteFlagCondition != nil {
				if !planCondition.PrerequisiteFlagCondition.PrerequisiteSettingId.Equal(stateCondition.PrerequisiteFlagCondition.PrerequisiteSettingId) ||
					!planCondition.PrerequisiteFlagCondition.Comparator.Equal(stateCondition.PrerequisiteFlagCondition.Comparator) ||
					hasSettingValueChanges(planCondition.PrerequisiteFlagCondition.ComparisonValue, stateCondition.PrerequisiteFlagCondition.ComparisonValue) {
					return true
				}
			}
		}

		for percentageOptionIndex, planPercentageOption := range planTargetingRule.PercentageOptions {
			statePercentageOption := (stateTargetingRule.PercentageOptions)[percentageOptionIndex]
			if !planPercentageOption.Percentage.Equal(statePercentageOption.Percentage) ||
				hasSettingValueChanges(planPercentageOption.Value, statePercentageOption.Value) {
				return true
			}
		}
	}

	return false
}

func hasSettingValueChanges(plan *settingValueModel, state *settingValueModel) bool {
	if (plan == nil) != (state == nil) {
		return true
	}

	if plan == nil {
		return false
	}

	return !plan.BoolValue.Equal(state.BoolValue) ||
		!plan.StringValue.Equal(state.StringValue) ||
		!plan.IntValue.Equal(state.IntValue) ||
		!plan.DoubleValue.Equal(state.DoubleValue)
}
