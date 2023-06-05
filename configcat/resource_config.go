package configcat

import (
	"context"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigCreate,
		ReadContext:   resourceConfigRead,
		UpdateContext: resourceConfigUpdate,
		DeleteContext: resourceConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
				ForceNew:     true,
			},

			CONFIG_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			CONFIG_DESCRIPTION: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	productID := d.Get(PRODUCT_ID).(string)

	body := sw.CreateConfigRequest{
		Name:        d.Get(CONFIG_NAME).(string),
		Description: *sw.NewNullableString(d.Get(CONFIG_DESCRIPTION).(*string)),
	}

	config, err := c.CreateConfig(productID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*config.ConfigId)

	return resourceConfigRead(ctx, d, m)
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	config, err := c.GetConfig(d.Id())
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(PRODUCT_ID, config.Product.ProductId)
	d.Set(CONFIG_NAME, config.Name)
	d.Set(CONFIG_DESCRIPTION, config.Description)

	return diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.HasChanges(CONFIG_NAME, CONFIG_DESCRIPTION) {
		body := sw.UpdateConfigRequest{
			Name:        *sw.NewNullableString(d.Get(CONFIG_NAME).(*string)),
			Description: *sw.NewNullableString(d.Get(CONFIG_DESCRIPTION).(*string)),
		}

		_, err := c.UpdateConfig(d.Id(), body)
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				d.SetId("")
				var diags diag.Diagnostics
				return diags
			}
			return diag.FromErr(err)
		}
	}

	return resourceConfigRead(ctx, d, m)
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.DeleteConfig(d.Id())
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
