package configcat

import (
	"context"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingCreate,
		ReadContext:   resourceSettingRead,
		UpdateContext: resourceSettingUpdate,
		DeleteContext: resourceSettingDelete,

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

			SETTING_NAME: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			SETTING_HINT: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			SETTING_TYPE: &schema.Schema{
				Type:    schema.TypeString,
				Default: "boolean",
			},

			SETTING_ID: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	settingID, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	settingKey := d.Get(SETTING_KEY).(string)

	setting, err := getSetting(c, int32(settingID))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	updateSettingResourceData(d, setting)
	var diags diag.Diagnostics
	return diags
}

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSettingRead(ctx, d, m)
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func getSetting(c *Client, settingID int32) (*sw.SettingModel, error) {

	setting, err := c.GetSetting(settingID)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}
