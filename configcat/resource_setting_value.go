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

			SETTING_VALUE: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			FREEZE_AFTER_INIT: &schema.Schema{
				Type:    schema.TypeBool,
				Default: true,
			},

			SETTING_VALUE_ID: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			SETTING_TYPE: &schema.Schema{
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
	return resourceSettingValueReadInternal(ctx, d, m, false)
}

func resourceSettingValueReadInternal(ctx context.Context, d *schema.ResourceData, m interface{}, forceRead bool) diag.Diagnostics {
	c := m.(*Client)

	id := d.Id()
	freezeAfterInit := d.Get(FREEZE_AFTER_INIT).(bool)

	if !forceRead && freezeAfterInit && id != "" {
		var diags diag.Diagnostics
		return diags
	}

	environmentID := d.Get(ENVIRONMENT_ID).(string)
	settingID, err := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	settingValue, err := c.GetSettingValue(environmentID, settingID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d.%d", settingValue.Environment.EnvironmentId, settingValue.Setting.SettingId))

	d.Set(SETTING_VALUE, settingValue.Value)
	d.Set(ROLLOUT_RULES, flattenRolloutRulesData(&settingValue.RolloutRules))
	d.Set(ROLLOUT_PERCENTAGE_ITEMS, flattenRolloutPercentageItemsData(&settingValue.RolloutPercentageItems))
}

func resourceSettingValueCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	environmentID := d.Get(ENVIRONMENT_ID).(string)
	settingID, err := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	body := sw.UpdateSettingValueModel{
		Value: d.Get(SETTING_VALUE).(string),
	}

	settingValue, err := c.ReplaceSettingValue(environmentID, settingID, body)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d.%d", settingValue.Environment.EnvironmentId, settingValue.Setting.SettingId))

	return resourceSettingValueReadInternal(ctx, d, m, true)
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
			element[ROLLOUT_RULE_VALUE] = rolloutRule.Value

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
			element[ROLLOUT_PERCENTAGE_ITEM_VALUE] = rolloutPercentageItem.Value

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
