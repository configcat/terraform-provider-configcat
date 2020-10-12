package configcat

import (
	"fmt"
	"log"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatSetting() *schema.Resource {
	return &schema.Resource{

		Create: settingCreate,
		Read:   settingRead,
		Update: settingUpdate,
		Delete: settingDelete,
		Exists: settingExists,

		Schema: map[string]*schema.Schema{

			"config_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	log.Printf("settingCreate")
	c := meta.(*Client)

	configId := d.Get("config_id").(string)
	if configId == "" {
		return fmt.Errorf("config_id is required")
	}

	settingType, err := getSettingType(d.Get("setting_type").(string))
	if err != nil {
		return err
	}

	body := sw.CreateSettingModel{
		Key:         d.Get("key").(string),
		Name:        d.Get("name").(string),
		SettingType: &settingType,
		Hint:        d.Get("hint").(string),
	}

	setting, err := c.CreateSetting(configId, body)
	if err != nil {
		return err
	}

	updateSettingResourceData(d, &setting)

	return nil
}

func getSettingType(settingType string) (sw.SettingType, error) {
	switch settingType {
	case "boolean":
		return sw.BOOLEAN_SettingType, nil
	case "string":
		return sw.STRING__SettingType, nil
	case "int":
		return sw.INT__SettingType, nil
	case "double":
		return sw.DOUBLE_SettingType, nil
	default:
		return sw.BOOLEAN_SettingType, fmt.Errorf("Could not parse setting_type: %s", settingType)
	}
}

func settingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("settingUpdate")
	//c := meta.(*Client)

	// Temporary workaround
	setting, err := findSetting(d, meta)
	if err != nil {
		return err
	}
	updateSettingResourceData(d, setting)

	/*settingId, ok := d.Id().(int32)
	if !ok {
		return fmt.Errorf("setting_id is required")
	}
	d.SetId(fmt.Sprintf("%v", settingId))
	d.SetId(fmt.Sprintf("%v", settingId))
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
	log.Printf("settingRead")
	setting, err := findSetting(d, meta)
	if err != nil {
		return err
	}

	updateSettingResourceData(d, setting)
	return nil
}

func settingExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("settingExists")
	_, err := findSetting(d, meta)
	if err != nil {
		return false, err
	}

	return true, nil
}

func settingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("settingDelete")
	c := meta.(*Client)

	settingId, convErr := strconv.Atoi(d.Id())
	if convErr != nil {
		return fmt.Errorf("could not parse setting_id")
	}

	err := c.DeleteSetting(int32(settingId))
	if err != nil {
		return err
	}

	return nil
}

func findSetting(d *schema.ResourceData, meta interface{}) (*sw.SettingModel, error) {
	c := meta.(*Client)

	settingId, convErr := strconv.Atoi(d.Id())
	if convErr == nil {
		setting, err := c.GetSetting(int32(settingId))
		if err != nil {
			return nil, err
		}

		return &setting, nil
	}

	configId := d.Get("config_id").(string)
	if configId == "" {
		return nil, fmt.Errorf("config_id+setting_key is required")
	}

	settingKey := d.Get("key").(string)
	if settingKey == "" {
		return nil, fmt.Errorf("config_id+setting_key is required")
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
	d.Set("config_id", m.ConfigId)
	d.Set("key", m.Key)
	d.Set("name", m.Name)
	d.Set("setting_type", m.SettingType)
	d.Set("hint", m.Hint)
}
