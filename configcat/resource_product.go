package configcat

import (
	"context"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatProduct() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProductCreate,
		ReadContext:   resourceProductRead,
		UpdateContext: resourceProductUpdate,
		DeleteContext: resourceProductDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			ORGANIZATION_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
				ForceNew:     true,
			},

			PRODUCT_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceProductCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	organizationID := d.Get(ORGANIZATION_ID).(string)

	body := sw.CreateProductRequest{
		Name: d.Get(PRODUCT_NAME).(string),
	}

	product, err := c.CreateProduct(organizationID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(product.ProductId)

	return resourceProductRead(ctx, d, m)
}

func resourceProductRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	product, err := c.GetProduct(d.Id())
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(ORGANIZATION_ID, product.Organization.OrganizationId)
	d.Set(PRODUCT_NAME, product.Name)

	return diags
}

func resourceProductUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.HasChanges(PRODUCT_NAME) {
		body := sw.UpdateProductRequest{
			Name: d.Get(PRODUCT_NAME).(string),
		}

		_, err := c.UpdateProduct(d.Id(), body)
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				d.SetId("")
				var diags diag.Diagnostics
				return diags
			}
			return diag.FromErr(err)
		}
	}

	return resourceProductRead(ctx, d, m)
}

func resourceProductDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.DeleteProduct(d.Id())
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
