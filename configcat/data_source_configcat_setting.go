package configcat

import (
	"context"
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatSetting() *schema.Resource {
	return &schema.Resource{

		ReadContext: settingRead,

		Schema: map[string]*schema.Schema{
			CONFIG_ID: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			SETTING_KEY: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			SETTING_ID: &schema.Schema{
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
	}
}

func settingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	configID := d.Get(CONFIG_ID).(string)
	settingKey := d.Get(SETTING_KEY).(string)

	setting, err := findSetting(c, configID, settingKey)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	updateSettingResourceData(d, setting, configID)
	var diags diag.Diagnostics
	return diags
}

func getSettings(c *Client, configID string) ([]sw.SettingModel, error) {

	settings, err := c.GetSettings(configID)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func findSetting(c *Client, configID, settingKey string) (*sw.SettingModel, error) {
	settings, err := getSettings(c, configID)
	if err != nil {
		return nil, err
	}

	for i := range settings {
		if settings[i].Key == settingKey {
			return &settings[i], nil
		}
	}

	return nil, fmt.Errorf("could not find Setting. config_id: %s key: %s", configID, settingKey)
}

func updateSettingResourceData(d *schema.ResourceData, m *sw.SettingModel, configID string) {
	settingID := fmt.Sprintf("%d", m.SettingId)

	d.SetId(settingID)
	d.Set(CONFIG_ID, configID)
	d.Set(SETTING_ID, settingID)
	d.Set(SETTING_KEY, m.Key)
	d.Set(SETTING_NAME, m.Name)
	d.Set(SETTING_HINT, m.Hint)
	d.Set(SETTING_TYPE, m.SettingType)
}
