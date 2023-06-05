package configcat

import (
	"context"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
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

			ENVIRONMENT_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			ENVIRONMENT_DESCRIPTION: {
				Type:     schema.TypeString,
				Optional: true,
			},
			ENVIRONMENT_COLOR: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	productID := d.Get(PRODUCT_ID).(string)

	environmentDescription := d.Get(ENVIRONMENT_DESCRIPTION).(string)
	environmentColor := d.Get(ENVIRONMENT_COLOR).(string)
	body := sw.CreateEnvironmentModel{
		Name:        d.Get(ENVIRONMENT_NAME).(string),
		Description: *sw.NewNullableString(&environmentDescription),
		Color:       *sw.NewNullableString(&environmentColor),
	}

	environment, err := c.CreateEnvironment(productID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*environment.EnvironmentId)

	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	environment, err := c.GetEnvironment(d.Id())
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(PRODUCT_ID, environment.Product.ProductId)
	d.Set(ENVIRONMENT_NAME, environment.Name.Get())
	d.Set(ENVIRONMENT_DESCRIPTION, environment.Description.Get())
	d.Set(ENVIRONMENT_COLOR, environment.Color.Get())

	return diags
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.HasChanges(ENVIRONMENT_NAME, ENVIRONMENT_DESCRIPTION, ENVIRONMENT_COLOR) {
		environmentName := d.Get(ENVIRONMENT_NAME).(string)
		environmentDescription := d.Get(ENVIRONMENT_DESCRIPTION).(string)
		environmentColor := d.Get(ENVIRONMENT_COLOR).(string)

		body := sw.UpdateEnvironmentModel{
			Name:        *sw.NewNullableString(&environmentName),
			Description: *sw.NewNullableString(&environmentDescription),
			Color:       *sw.NewNullableString(&environmentColor),
		}

		_, err := c.UpdateEnvironment(d.Id(), body)
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				d.SetId("")
				var diags diag.Diagnostics
				return diags
			}

			return diag.FromErr(err)
		}
	}

	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.DeleteEnvironment(d.Id())
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
