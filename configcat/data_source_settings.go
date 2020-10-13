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
			CONFIG_ID: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			SETTING_KEY_FILTER_REGEX: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			SETTINGS: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SETTING_ID: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_KEY: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_NAME: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_HINT: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						SETTING_TYPE: &schema.Schema{
							Type:     schema.TypeString,
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
			if regex.MatchString(settings[i].Key) {
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
		elements := make([]interface{}, len(*settings), len(*settings))

		for i, setting := range *settings {
			element := make(map[string]interface{})

			settingID := fmt.Sprintf("%d", setting.SettingId)

			element[SETTING_ID] = settingID
			element[SETTING_KEY] = setting.Key
			element[SETTING_NAME] = setting.Name
			element[SETTING_HINT] = setting.Hint
			element[SETTING_TYPE] = setting.SettingType

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
