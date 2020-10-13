package configcat

import (
	"context"
	"fmt"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingCreate,
		ReadContext:   resourceSettingRead,
		UpdateContext: resourceSettingUpdate,
		DeleteContext: resourceSettingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			CONFIG_ID: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
				ForceNew:     true,
			},

			SETTING_KEY: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			SETTING_TYPE: &schema.Schema{
				Type:     schema.TypeString,
				Default:  "boolean",
				Optional: true,
				ForceNew: true,
			},

			SETTING_NAME: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			SETTING_HINT: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	configID := d.Get(CONFIG_ID).(string)

	settingTypeString := d.Get(SETTING_TYPE).(string)
	var settingType sw.SettingType
	switch d.Get(SETTING_TYPE).(string) {
	case "boolean":
		settingType = sw.BOOLEAN_SettingType
	case "string":
		settingType = sw.STRING__SettingType
	case "int":
		settingType = sw.INT__SettingType
	case "double":
		settingType = sw.DOUBLE_SettingType
	default:
		d.SetId("")
		return diag.FromErr(fmt.Errorf("setting_type parse failed: %s. Valid values: boolean/string/int/double", settingTypeString))
	}

	body := sw.CreateSettingModel{
		Key:         d.Get(SETTING_KEY).(string),
		Name:        d.Get(SETTING_NAME).(string),
		Hint:        d.Get(SETTING_HINT).(string),
		SettingType: &settingType,
	}

	setting, err := c.CreateSetting(configID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(setting.SettingId))

	return resourceSettingRead(ctx, d, m)
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)
	settingID, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	setting, err := c.GetSetting(int32(settingID))
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(SETTING_KEY, setting.Key)
	d.Set(SETTING_NAME, setting.Name)
	d.Set(SETTING_HINT, setting.Hint)
	d.Set(SETTING_TYPE, setting.SettingType)
	d.Set(CONFIG_ID, setting.ConfigId)

	return diags
}

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	settingID, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(SETTING_NAME, SETTING_HINT) {
		body := []sw.Operation{}
		if d.HasChange(SETTING_NAME) {
			settingName := d.Get(SETTING_NAME)
			body = append(body, sw.Operation{
				Op:    "replace",
				Path:  "/name",
				Value: &settingName,
			})
		}

		if d.HasChange(SETTING_HINT) {
			settingHint := d.Get(SETTING_HINT)
			body = append(body, sw.Operation{
				Op:    "replace",
				Path:  "/hint",
				Value: &settingHint,
			})
		}

		_, err := c.UpdateSetting(int32(settingID), body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSettingRead(ctx, d, m)
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)
	settingID, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteSetting(int32(settingID))
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
