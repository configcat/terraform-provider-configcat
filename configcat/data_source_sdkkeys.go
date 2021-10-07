package configcat

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatSdkKeys() *schema.Resource {
	return &schema.Resource{

		ReadContext: sdkKeysRead,

		Schema: map[string]*schema.Schema{
			CONFIG_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			ENVIRONMENT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			PRIMARY_SDK_KEY: {
				Type:     schema.TypeString,
				Computed: true,
			},

			SECONDARY_SDK_KEY: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func sdkKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	configID := d.Get(CONFIG_ID).(string)
	environmentID := d.Get(ENVIRONMENT_ID).(string)

	sdkKeys, err := c.GetSdkKeys(configID, environmentID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(PRIMARY_SDK_KEY, sdkKeys.Primary)
	d.Set(SECONDARY_SDK_KEY, sdkKeys.Secondary)

	var diags diag.Diagnostics
	return diags
}
