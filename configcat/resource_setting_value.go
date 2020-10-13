package configcat

import (
	"context"
	"fmt"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatSettingValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingValueCreateOrUpdate,
		ReadContext:   resourceSettingValueRead,
		UpdateContext: resourceSettingValueCreateOrUpdate,
		DeleteContext: resourceSettingValueDelete,

		Schema: map[string]*schema.Schema{
			ENVIRONMENT_ID: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
				ForceNew:     true,
			},

			SETTING_ID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			SETTING_TYPE: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			SETTING_VALUE: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			FREEZE_AFTER_INIT: &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			SETTING_VALUE_ID: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			ROLLOUT_RULES: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ROLLOUT_RULE_COMPARISON_ATTRIBUTE: &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						ROLLOUT_RULE_COMPARATOR: &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						ROLLOUT_RULE_COMPARISON_VALUE: &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						ROLLOUT_RULE_VALUE: &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			ROLLOUT_PERCENTAGE_ITEMS: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE: &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						ROLLOUT_PERCENTAGE_ITEM_VALUE: &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSettingValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := resourceSettingValueReadInternal(ctx, d, m, false)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	return diags
}

func resourceSettingValueCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	freezeAfterInit := d.Get(FREEZE_AFTER_INIT).(bool)

	var diags diag.Diagnostics
	if freezeAfterInit && id != "" {
		if d.HasChanges() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Changes will be only applied to the state.",
				Detail:   "The freeze_after_init parameter is set to true so the changes won't be applied in ConfigCat. This mode is only for initializing a feature flag in ConfigCat.",
			})
		}

		return diags
	}

	c := m.(*Client)

	environmentID := d.Get(ENVIRONMENT_ID).(string)
	settingID, convErr := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if convErr != nil {
		d.SetId("")
		return diag.FromErr(convErr)
	}

	// Read the settingtype first so we know about the settingTypes
	settingTypeString := d.Get(SETTING_TYPE).(string)

	settingValue, settingValueErr := getSettingValue(settingTypeString, d.Get(SETTING_VALUE).(string))

	if settingValueErr != nil {
		d.SetId("")
		return diag.FromErr(settingValueErr)
	}

	rolloutRules, rolloutRulesErr := getRolloutRulesData(d.Get(ROLLOUT_RULES).([]interface{}), settingTypeString)
	if rolloutRulesErr != nil {
		return diag.FromErr(rolloutRulesErr)
	}

	rolloutPercentageItems, rolloutPercentageItemsErr := getRolloutPercentageItemsData(d.Get(ROLLOUT_PERCENTAGE_ITEMS).([]interface{}), settingTypeString)
	if rolloutPercentageItemsErr != nil {
		return diag.FromErr(rolloutPercentageItemsErr)
	}

	body := sw.UpdateSettingValueModel{
		Value:                  &settingValue,
		RolloutRules:           *rolloutRules,
		RolloutPercentageItems: *rolloutPercentageItems,
	}

	_, err := c.ReplaceSettingValue(environmentID, int32(settingID), body)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s.%d", environmentID, settingID))

	readErr2 := resourceSettingValueReadInternal(ctx, d, m, true)
	if readErr2 != nil {
		d.SetId("")
		return diag.FromErr(readErr2)
	}

	return diags
}

func resourceSettingValueReadInternal(ctx context.Context, d *schema.ResourceData, m interface{}, forceRead bool) error {
	c := m.(*Client)

	id := d.Id()
	freezeAfterInit := d.Get(FREEZE_AFTER_INIT).(bool)

	if !forceRead && freezeAfterInit && id != "" {
		return nil
	}

	environmentID := d.Get(ENVIRONMENT_ID).(string)
	settingID, err := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if err != nil {
		d.SetId("")
		return err
	}

	settingValue, err := c.GetSettingValueSimple(environmentID, int32(settingID))
	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(fmt.Sprintf("%s.%d", environmentID, settingID))

	d.Set(SETTING_VALUE, fmt.Sprintf("%v", *settingValue.Value))
	d.Set(ROLLOUT_RULES, flattenRolloutRulesData(&settingValue.RolloutRules))
	d.Set(ROLLOUT_PERCENTAGE_ITEMS, flattenRolloutPercentageItemsData(&settingValue.RolloutPercentageItems))

	return nil
}

func resourceSettingValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}

func flattenRolloutRulesData(rolloutRules *[]sw.RolloutRuleModel) []interface{} {
	if rolloutRules != nil {
		elements := make([]interface{}, len(*rolloutRules), len(*rolloutRules))

		for i, rolloutRule := range *rolloutRules {
			element := make(map[string]interface{})

			element[ROLLOUT_RULE_COMPARISON_ATTRIBUTE] = rolloutRule.ComparisonAttribute
			element[ROLLOUT_RULE_COMPARATOR] = rolloutRule.Comparator
			element[ROLLOUT_RULE_COMPARISON_VALUE] = rolloutRule.ComparisonValue
			element[ROLLOUT_RULE_VALUE] = fmt.Sprintf("%v", rolloutRule.Value)

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}

func flattenRolloutPercentageItemsData(rolloutPercentageItems *[]sw.RolloutPercentageItemModel) []interface{} {
	if rolloutPercentageItems != nil {
		elements := make([]interface{}, len(*rolloutPercentageItems), len(*rolloutPercentageItems))

		for i, rolloutPercentageItem := range *rolloutPercentageItems {
			element := make(map[string]interface{})

			element[ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE] = rolloutPercentageItem.Percentage
			element[ROLLOUT_PERCENTAGE_ITEM_VALUE] = fmt.Sprintf("%v", rolloutPercentageItem.Value)
			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}

func getRolloutRulesData(rolloutRules []interface{}, settingType string) (*[]sw.RolloutRuleModel, error) {
	if rolloutRules != nil {
		elements := make([]sw.RolloutRuleModel, len(rolloutRules), len(rolloutRules))

		for i, rolloutRule := range rolloutRules {
			item := rolloutRule.(map[string]interface{})

			value, err := getSettingValue(settingType, item[ROLLOUT_RULE_VALUE].(string))
			if err != nil {
				return nil, err
			}

			comparator, compErr := getComparator(item[ROLLOUT_RULE_COMPARATOR].(string))
			if compErr != nil {
				return nil, compErr
			}

			element := sw.RolloutRuleModel{
				ComparisonAttribute: item[ROLLOUT_RULE_COMPARISON_ATTRIBUTE].(string),
				Comparator:          comparator,
				ComparisonValue:     item[ROLLOUT_RULE_COMPARISON_VALUE].(string),
				Value:               &value,
			}

			elements[i] = element
		}

		return &elements, nil
	}
	empty := make([]sw.RolloutRuleModel, 0)
	return &empty, nil
}

func getRolloutPercentageItemsData(rolloutPercentageItems []interface{}, settingType string) (*[]sw.RolloutPercentageItemModel, error) {
	if rolloutPercentageItems != nil {
		elements := make([]sw.RolloutPercentageItemModel, len(rolloutPercentageItems), len(rolloutPercentageItems))

		for i, rolloutPercentageItem := range rolloutPercentageItems {
			item := rolloutPercentageItem.(map[string]interface{})

			value, err := getSettingValue(settingType, item[ROLLOUT_PERCENTAGE_ITEM_VALUE].(string))
			if err != nil {
				return nil, err
			}

			element := sw.RolloutPercentageItemModel{
				Percentage: item[ROLLOUT_PERCENTAGE_ITEM_VALUE].(int64),
				Value:      &value,
			}

			elements[i] = element
		}

		return &elements, nil
	}

	empty := make([]sw.RolloutPercentageItemModel, 0)
	return &empty, nil
}

func getSettingValue(settingType, value string) (interface{}, error) {

	switch settingType {
	case "boolean":
		b, err := strconv.ParseBool(value)
		return b, err
	case "string":
		return value, nil
	case "int":
		i, err := strconv.ParseInt(value, 10, 32)
		if err == nil {
			return int32(i), nil
		}
		return nil, err
	case "double":
		f, err := strconv.ParseFloat(value, 64)
		return f, err
	default:
		return nil, fmt.Errorf("Could not parse SettingType and Value: %s, %s", settingType, value)
	}
}

func getComparator(comparator string) (*sw.RolloutRuleComparator, error) {
	switch comparator {
	case "isOneOf":
		comparator := sw.IS_ONE_OF_RolloutRuleComparator
		return &comparator, nil
	case "isNotOneOf":
		comparator := sw.IS_NOT_ONE_OF_RolloutRuleComparator
		return &comparator, nil
	case "contains":
		comparator := sw.CONTAINS_RolloutRuleComparator
		return &comparator, nil
	case "doesNotContain":
		comparator := sw.DOES_NOT_CONTAIN_RolloutRuleComparator
		return &comparator, nil
	case "semVerIsOneOf":
		comparator := sw.SEM_VER_IS_ONE_OF_RolloutRuleComparator
		return &comparator, nil
	case "semVerIsNotOneOf":
		comparator := sw.SEM_VER_IS_NOT_ONE_OF_RolloutRuleComparator
		return &comparator, nil
	case "semVerLess":
		comparator := sw.SEM_VER_LESS_RolloutRuleComparator
		return &comparator, nil
	case "semVerLessOrEquals":
		comparator := sw.SEM_VER_LESS_OR_EQUALS_RolloutRuleComparator
		return &comparator, nil
	case "semVerGreater":
		comparator := sw.SEM_VER_GREATER_RolloutRuleComparator
		return &comparator, nil
	case "semVerGreaterOrEquals":
		comparator := sw.SEM_VER_GREATER_OR_EQUALS_RolloutRuleComparator
		return &comparator, nil
	case "numberEquals":
		comparator := sw.NUMBER_EQUALS_RolloutRuleComparator
		return &comparator, nil
	case "numberDoesNotEqual":
		comparator := sw.NUMBER_DOES_NOT_EQUAL_RolloutRuleComparator
		return &comparator, nil
	case "numberLess":
		comparator := sw.NUMBER_LESS_RolloutRuleComparator
		return &comparator, nil
	case "numberLessOrEquals":
		comparator := sw.NUMBER_LESS_OR_EQUALS_RolloutRuleComparator
		return &comparator, nil
	case "numberGreater":
		comparator := sw.NUMBER_GREATER_RolloutRuleComparator
		return &comparator, nil
	case "numberGreaterOrEquals":
		comparator := sw.NUMBER_GREATER_OR_EQUALS_RolloutRuleComparator
		return &comparator, nil
	case "sensitiveIsOneOf":
		comparator := sw.SENSITIVE_IS_ONE_OF_RolloutRuleComparator
		return &comparator, nil
	case "sensitiveIsNotOneOf":
		comparator := sw.SENSITIVE_IS_NOT_ONE_OF_RolloutRuleComparator
		return &comparator, nil
	}

	return nil, fmt.Errorf("could not parse Comparator: %s", comparator)
}
