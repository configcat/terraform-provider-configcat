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

func resourceConfigCatSettingTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingTagCreate,
		ReadContext:   resourceSettingTagRead,
		DeleteContext: resourceSettingTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			SETTING_ID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			TAG_ID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSettingTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	settingID, sConvErr := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if sConvErr != nil {
		return diag.FromErr(sConvErr)
	}

	tagIDInterface := d.Get(TAG_ID)
	tagID, tConvErr := strconv.ParseInt(tagIDInterface.(string), 10, 32)
	if tConvErr != nil {
		return diag.FromErr(tConvErr)
	}

	body := []sw.Operation{sw.Operation{
		Op:    "add",
		Path:  "/tags/-",
		Value: &tagIDInterface,
	}}

	_, err := c.UpdateSetting(int32(settingID), body)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			var diags diag.Diagnostics
			return diags
		}

		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", settingID, tagID))

	return resourceSettingTagRead(ctx, d, m)
}

func resourceSettingTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	settingID, tagID, convErr := resourceConfigCatSettingTagParseID(d.Id())
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	setting, err := c.GetSetting(settingID)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	for _, tag := range setting.Tags {
		if tag.TagId == tagID {
			d.Set(SETTING_ID, fmt.Sprintf("%d", settingID))
			d.Set(TAG_ID, fmt.Sprintf("%d", tagID))

			return diags
		}
	}

	d.SetId("")
	return diags
}

func resourceSettingTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	settingID, sConvErr := strconv.ParseInt(d.Get(SETTING_ID).(string), 10, 32)
	if sConvErr != nil {
		return diag.FromErr(sConvErr)
	}

	tagIDInterface := d.Get(TAG_ID)

	body := []sw.Operation{sw.Operation{
		Op:    "remove",
		Path:  "/tags/-",
		Value: &tagIDInterface,
	}}

	_, err := c.UpdateSetting(int32(settingID), body)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			var diags diag.Diagnostics
			return diags
		}

		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceConfigCatSettingTagParseID(id string) (int32, int64, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected settingID.tagID", id)
	}

	settingID, sConvErr := strconv.ParseInt(parts[0], 10, 32)
	if sConvErr != nil {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected settingID.tagID. Error: %s", id, sConvErr)
	}

	tagID, tConvErr := strconv.ParseInt(parts[1], 10, 32)
	if tConvErr != nil {
		return 0, 0, fmt.Errorf("unexpected format of ID (%s), expected settingID.tagID. Error: %s", id, tConvErr)
	}

	return int32(settingID), tagID, nil
}
