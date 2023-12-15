package configcat

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatSettings() *schema.Resource {
	return &schema.Resource{

		ReadContext: settingRead,

		Schema: map[string]*schema.Schema{
			CONFIG_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			SETTING_KEY_FILTER_REGEX: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			SETTINGS: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SETTING_ID: {
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_KEY: {
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_NAME: {
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_HINT: {
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_TYPE: {
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_ORDER: {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func settingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	configID := d.Get(CONFIG_ID).(string)
	settingKeyFilterRegex := d.Get(SETTING_KEY_FILTER_REGEX).(string)

	settings, err := c.GetSettings(configID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredSettings := []sw.SettingModel{}
	if settingKeyFilterRegex == "" {
		filteredSettings = settings
	} else {
		regex := regexp.MustCompile(settingKeyFilterRegex)
		for i := range settings {
			if regex.MatchString(*settings[i].Key.Get()) {
				filteredSettings = append(filteredSettings, settings[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(SETTINGS, flattenSettingsData(&filteredSettings))

	var diags diag.Diagnostics
	return diags
}

func flattenSettingsData(settings *[]sw.SettingModel) []interface{} {
	if settings != nil {
		elements := make([]interface{}, len(*settings))

		for i, setting := range *settings {
			element := make(map[string]interface{})

			element[SETTING_ID] = fmt.Sprintf("%d", *setting.SettingId)
			element[SETTING_KEY] = setting.Key.Get()
			element[SETTING_NAME] = setting.Name.Get()
			element[SETTING_HINT] = setting.Hint.Get()
			element[SETTING_TYPE] = setting.SettingType
			element[SETTING_ORDER] = setting.Order

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
