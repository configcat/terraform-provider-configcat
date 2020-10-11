package configcat

import (
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConfigCatSetting() *schema.Resource {
	return &schema.Resource{

		Create: settingCreate,
		Read:   settingRead,
		Update: settingUpdate,
		Delete: settingDelete,
		Exists: settingExists,

		Schema: map[string]*schema.Schema{

			"setting_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"config_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"setting_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Default:  "boolean",
			},

			"hint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func settingCreate(d *schema.ResourceData, meta interface{}) error {

}

func settingRead(d *schema.ResourceData, meta interface{}) error {
	setting, err := findSetting(d, meta)
	if err != nil {
		return err
	}

	updateSettingResourceData(d, setting)
	return nil
}

func findSetting(d *schema.ResourceData, meta interface{}) (*sw.SettingModel, error) {
	c := meta.(*Client)

	settingIdRaw := d.Get("setting_id")
	settingId, ok := settingIdRaw.(int32)
	if ok {
		setting, err := c.GetSetting(settingId)
		if err != nil {
			return nil, err
		}

		return &setting, nil
	}

	configId := d.Get("config_id").(string)
	if configId == "" {
		return nil, fmt.Errorf("setting_id or config_id+setting_key is required")
	}

	settingKey := d.Get("key").(string)
	if settingKey == "" {
		return nil, fmt.Errorf("setting_id or config_id+setting_key is required")
	}

	settings, err := c.GetSettings(configId)
	if err != nil {
		return nil, err
	}

	for i := range settings {
		if settings[i].Key == settingKey {
			return &settings[i], nil
		}
	}

	return nil, fmt.Errorf("could not find Setting")
}

func updateSettingResourceData(d *schema.ResourceData, m *sw.SettingModel) {
	d.SetId(fmt.Sprintf("%v", m.SettingId))
	d.Set("setting_id", m.SettingId)
	d.Set("config_id", m.ConfigId)
	d.Set("key", m.Key)
	d.Set("name", m.Name)
	d.Set("setting_type", m.SettingType)
	d.Set("hint", m.Hint)
}
