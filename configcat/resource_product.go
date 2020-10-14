package configcat

import (
	"context"
	"strconv"

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
				ForceNew: true,
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
	productID, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	product, err := c.GetProduct(int32(productID))
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(PRODUCT_KEY, product.Key)
	d.Set(PRODUCT_NAME, product.Name)
	d.Set(PRODUCT_HINT, product.Hint)
	d.Set(PRODUCT_TYPE, product.ProductType)
	d.Set(CONFIG_ID, product.ConfigId)

	return diags
}

func resourceProductUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(PRODUCT_NAME, PRODUCT_HINT) {
		body := []sw.Operation{}
		if d.HasChange(PRODUCT_NAME) {
			productName := d.Get(PRODUCT_NAME)
			body = append(body, sw.Operation{
				Op:    "replace",
				Path:  "/name",
				Value: &productName,
			})
		}

		if d.HasChange(PRODUCT_HINT) {
			productHint := d.Get(PRODUCT_HINT)
			body = append(body, sw.Operation{
				Op:    "replace",
				Path:  "/hint",
				Value: &productHint,
			})
		}

		_, err := c.UpdateProduct(int32(productID), body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceProductRead(ctx, d, m)
}

func resourceProductDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)
	productID, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteProduct(int32(productID))
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
