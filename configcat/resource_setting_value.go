package configcat

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				// Here we use a function to parse the import ID (like the example above) to simplify our logic
				environmentID, settingID, err := resourceConfigCatSettingValueParseID(d.Id())

				if err != nil {
					return nil, err
				}

				d.Set(SETTING_ID, settingID)
				d.Set(ENVIRONMENT_ID, environmentID)
				d.Set(INIT_ONLY, false)
				d.SetId(fmt.Sprintf("%s:%s", environmentID, settingID))

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			ENVIRONMENT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
				ForceNew:     true,
			},

			SETTING_ID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			SETTING_VALUE: {
				Type:     schema.TypeString,
				Required: true,
			},

			INIT_ONLY: {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			SETTING_TYPE: {
				Type:     schema.TypeString,
				Computed: true,
			},

			MANDATORY_NOTES: {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},

			ROLLOUT_RULES: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ROLLOUT_RULE_COMPARISON_ATTRIBUTE: {
							Type:     schema.TypeString,
							Optional: true,
						},

						ROLLOUT_RULE_COMPARATOR: {
							Type:     schema.TypeString,
							Optional: true,
						},

						ROLLOUT_RULE_COMPARISON_VALUE: {
							Type:     schema.TypeString,
							Optional: true,
						},

						ROLLOUT_RULE_SEGMENT_COMPARATOR: {
							Type:     schema.TypeString,
							Optional: true,
						},

						ROLLOUT_RULE_SEGMENT_ID: {
							Type:     schema.TypeString,
							Optional: true,
						},

						ROLLOUT_RULE_VALUE: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			ROLLOUT_PERCENTAGE_ITEMS: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE: {
							Type:     schema.TypeString,
							Required: true,
						},

						ROLLOUT_PERCENTAGE_ITEM_VALUE: {
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
	var diags diag.Diagnostics

	err := resourceSettingValueReadInternal(ctx, d, m, false)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	return diags
}

func resourceSettingValueCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	freezeAfterInit := d.Get(INIT_ONLY).(bool)

	var diags diag.Diagnostics
	if freezeAfterInit && id != "" {
		if d.HasChanges() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Changes will be only applied to the state.",
				Detail:   "The init_only parameter is set to true so the changes won't be applied in ConfigCat. This mode is only for initializing a feature flag in ConfigCat.",
			})
		}

		return diags
	}

	c := m.(*Client)

	environmentID := d.Get(ENVIRONMENT_ID).(string)
	settingID, convErr := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if convErr != nil {
		return diag.FromErr(convErr)
	}
	mandatoryNotes := d.Get(MANDATORY_NOTES).(string)

	// Read the settingtype first so we know about the settingTypes
	settingTypeString := d.Get(SETTING_TYPE).(string)
	if settingTypeString == "" {
		setting, err := c.GetSetting(int32(settingID))
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				d.SetId("")
				return diags
			}

			return diag.FromErr(err)
		}

		settingTypeString = fmt.Sprintf("%v", *setting.SettingType)
	}

	settingType, settingTypeConvertErr := sw.NewSettingTypeFromValue(settingTypeString)
	if settingTypeConvertErr != nil {
		return diag.FromErr(settingTypeConvertErr)
	}

	settingValue, settingValueErr := getSettingValue(settingType, d.Get(SETTING_VALUE).(string))

	if settingValueErr != nil {
		return diag.FromErr(settingValueErr)
	}

	rolloutRules, rolloutRulesErr := getRolloutRulesData(d.Get(ROLLOUT_RULES).([]interface{}), settingType)
	if rolloutRulesErr != nil {
		return diag.FromErr(rolloutRulesErr)
	}

	rolloutPercentageItems, rolloutPercentageItemsErr := getRolloutPercentageItemsData(d.Get(ROLLOUT_PERCENTAGE_ITEMS).([]interface{}), settingType)
	if rolloutPercentageItemsErr != nil {
		return diag.FromErr(rolloutPercentageItemsErr)
	}

	body := sw.UpdateSettingValueModel{
		Value:                  &settingValue,
		RolloutRules:           *rolloutRules,
		RolloutPercentageItems: *rolloutPercentageItems,
	}

	_, err := c.ReplaceSettingValue(environmentID, int32(settingID), body, mandatoryNotes)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", environmentID, settingID))

	readErr := resourceSettingValueReadInternal(ctx, d, m, true)
	if readErr != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(readErr)
	}

	return diags
}

func resourceSettingValueReadInternal(ctx context.Context, d *schema.ResourceData, m interface{}, forceRead bool) error {
	c := m.(*Client)

	id := d.Id()
	freezeAfterInit := d.Get(INIT_ONLY).(bool)

	if !forceRead && freezeAfterInit && id != "" {
		return nil
	}

	environmentID := d.Get(ENVIRONMENT_ID).(string)
	settingID, err := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if err != nil {
		return err
	}

	settingValue, err := c.GetSettingValue(environmentID, int32(settingID))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%d", environmentID, settingID))

	d.Set(SETTING_VALUE, fmt.Sprintf("%v", settingValue.Value))
	d.Set(SETTING_TYPE, settingValue.Setting.SettingType)
	d.Set(ROLLOUT_RULES, flattenRolloutRulesData(settingValue.RolloutRules))
	d.Set(ROLLOUT_PERCENTAGE_ITEMS, flattenRolloutPercentageItemsData(settingValue.RolloutPercentageItems))

	return nil
}

func resourceSettingValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}

func flattenRolloutRulesData(rolloutRules []sw.RolloutRuleModel) []interface{} {
	if rolloutRules != nil {
		elements := make([]interface{}, len(rolloutRules))

		for i, rolloutRule := range rolloutRules {
			element := make(map[string]interface{})

			element[ROLLOUT_RULE_COMPARISON_ATTRIBUTE] = rolloutRule.ComparisonAttribute.Get()
			element[ROLLOUT_RULE_COMPARATOR] = rolloutRule.Comparator
			element[ROLLOUT_RULE_COMPARISON_VALUE] = rolloutRule.ComparisonValue.Get()
			element[ROLLOUT_RULE_SEGMENT_COMPARATOR] = rolloutRule.SegmentComparator
			element[ROLLOUT_RULE_SEGMENT_ID] = rolloutRule.SegmentId.Get()
			element[ROLLOUT_RULE_VALUE] = fmt.Sprintf("%v", rolloutRule.Value)

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}

func flattenRolloutPercentageItemsData(rolloutPercentageItems []sw.RolloutPercentageItemModel) []interface{} {
	if rolloutPercentageItems != nil {
		elements := make([]interface{}, len(rolloutPercentageItems))

		for i, rolloutPercentageItem := range rolloutPercentageItems {
			element := make(map[string]interface{})

			element[ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE] = strconv.FormatInt(rolloutPercentageItem.Percentage, 10)
			element[ROLLOUT_PERCENTAGE_ITEM_VALUE] = fmt.Sprintf("%v", rolloutPercentageItem.Value)
			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}

func getRolloutRulesData(rolloutRules []interface{}, settingType *sw.SettingType) (*[]sw.RolloutRuleModel, error) {
	if rolloutRules != nil {
		elements := make([]sw.RolloutRuleModel, len(rolloutRules))

		for i, rolloutRule := range rolloutRules {
			item := rolloutRule.(map[string]interface{})

			value, err := getSettingValue(settingType, item[ROLLOUT_RULE_VALUE].(string))
			if err != nil {
				return nil, err
			}

			if len(item[ROLLOUT_RULE_COMPARATOR].(string)) > 0 {
				if len(item[ROLLOUT_RULE_COMPARISON_ATTRIBUTE].(string)) == 0 {
					return nil, fmt.Errorf("the %s field is required", ROLLOUT_RULE_COMPARISON_ATTRIBUTE)
				}
				if len(item[ROLLOUT_RULE_COMPARISON_VALUE].(string)) == 0 {
					return nil, fmt.Errorf("the %s field is required", ROLLOUT_RULE_COMPARISON_VALUE)
				}

				comparator, compErr := sw.NewRolloutRuleComparatorFromValue(item[ROLLOUT_RULE_COMPARATOR].(string))
				if compErr != nil {
					return nil, compErr
				}

				comparisonAttribute := item[ROLLOUT_RULE_COMPARISON_ATTRIBUTE].(string)
				comparisonValue := item[ROLLOUT_RULE_COMPARISON_VALUE].(string)
				element := sw.RolloutRuleModel{
					ComparisonAttribute: *sw.NewNullableString(&comparisonAttribute),
					Comparator:          comparator,
					ComparisonValue:     *sw.NewNullableString(&comparisonValue),
					Value:               &value,
				}

				elements[i] = element
			} else if len(item[ROLLOUT_RULE_SEGMENT_COMPARATOR].(string)) > 0 {
				if len(item[ROLLOUT_RULE_SEGMENT_ID].(string)) == 0 {
					return nil, fmt.Errorf("the %s field is required", ROLLOUT_RULE_SEGMENT_ID)
				}

				segmentComparator, compErr := sw.NewSegmentComparatorFromValue(item[ROLLOUT_RULE_SEGMENT_COMPARATOR].(string))
				if compErr != nil {
					return nil, compErr
				}

				segmentId := item[SEGMENT_ID].(string)
				element := sw.RolloutRuleModel{
					SegmentComparator: segmentComparator,
					SegmentId:         *sw.NewNullableString(&segmentId),
					Value:             &value,
				}

				elements[i] = element
			} else {
				return nil, fmt.Errorf("either the Comparator or the SegmentComparator should be set")
			}
		}

		return &elements, nil
	}
	empty := make([]sw.RolloutRuleModel, 0)
	return &empty, nil
}

func getRolloutPercentageItemsData(rolloutPercentageItems []interface{}, settingType *sw.SettingType) (*[]sw.RolloutPercentageItemModel, error) {
	if rolloutPercentageItems != nil {
		elements := make([]sw.RolloutPercentageItemModel, len(rolloutPercentageItems))

		for i, rolloutPercentageItem := range rolloutPercentageItems {
			item := rolloutPercentageItem.(map[string]interface{})

			value, err := getSettingValue(settingType, item[ROLLOUT_PERCENTAGE_ITEM_VALUE].(string))
			if err != nil {
				return nil, err
			}

			percentage, percErr := strconv.ParseInt(item[ROLLOUT_PERCENTAGE_ITEM_PERCENTAGE].(string), 10, 32)
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

	empty := make([]sw.RolloutPercentageItemModel, 0)
	return &empty, nil
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
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected environmentID.settingID", id)
	}

	_, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected environmentID.settingID. Error: %s", id, err)
	}

	return parts[0], parts[1], nil
}
