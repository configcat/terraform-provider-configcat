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
				Optional: true,
			},

			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"setting_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "boolean",
				ForceNew: true,
			},

			"hint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func settingCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	configId := d.Get("config_id").(string)
	if configId == "" {
		return fmt.Errorf("config_id is required")
	}

	body := sw.CreateSettingModel{
		Key:         d.Get("key").(string),
		Name:        d.Get("name").(string),
		SettingType: d.Get("setting_type").(*sw.SettingType),
		Hint:        d.Get("hint").(string),
	}

	setting, err := c.CreateSetting(configId, body)
	if err != nil {
		return err
	}

	updateSettingResourceData(d, &setting)

	return nil
}

func settingUpdate(d *schema.ResourceData, meta interface{}) error {
	//c := meta.(*Client)

	_, ok := d.Get("setting_id").(int32)
	if !ok {
		return fmt.Errorf("setting_id is required")
	}
	/*
		body := sw.UpdateSettingModel{
			Key:         d.Get("key").(string),
			Name:        d.Get("name").(string),
			SettingType: d.Get("setting_type").(*sw.SettingType),
			Hint:        d.Get("hint").(string),
		}

		setting, err := c.UpdateSetting(settingId, body)
		if err != nil {
			return err
		}

		updateSettingResourceData(d, &setting)
	*/
	return nil
}

func settingRead(d *schema.ResourceData, meta interface{}) error {
	setting, err := findSetting(d, meta)
	if err != nil {
		return err
	}

	updateSettingResourceData(d, setting)
	return nil
}

func settingExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	_, err := findSetting(d, meta)
	if err != nil {
		return false, err
	}

	return true, nil
}

func settingDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	settingId, ok := d.Get("setting_id").(int32)
	if !ok {
		return fmt.Errorf("setting_id is required")
	}

	err := c.DeleteSetting(settingId)
	if err != nil {
		return err
	}

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
