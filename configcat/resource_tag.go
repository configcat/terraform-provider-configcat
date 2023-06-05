package configcat

import (
	"context"
	"strconv"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagCreate,
		ReadContext:   resourceTagRead,
		UpdateContext: resourceTagUpdate,
		DeleteContext: resourceTagDelete,
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

			TAG_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},

			TAG_COLOR: {
				Type:         schema.TypeString,
				ValidateFunc: validateColorFunc,
				Optional:     true,
				Default:      "panther",
			},
		},
	}
}

func resourceTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	productID := d.Get(PRODUCT_ID).(string)

	color := d.Get(TAG_COLOR).(string)
	body := sw.CreateTagModel{
		Name:  d.Get(TAG_NAME).(string),
		Color: *sw.NewNullableString(&color),
	}

	tag, err := c.CreateTag(productID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(*tag.TagId, 10))

	return resourceTagRead(ctx, d, m)
}

func resourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	tagID, convErr := strconv.ParseInt(d.Id(), 10, 64)
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	tag, err := c.GetTag(tagID)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(PRODUCT_ID, tag.Product.ProductId)
	d.Set(TAG_NAME, tag.Name.Get())
	d.Set(TAG_COLOR, tag.Color.Get())

	return diags
}

func resourceTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	tagID, convErr := strconv.ParseInt(d.Id(), 10, 64)
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	if d.HasChanges(TAG_NAME, TAG_COLOR) {

		name := d.Get(TAG_NAME).(string)
		color := d.Get(TAG_COLOR).(string)
		body := sw.UpdateTagModel{
			Name:  *sw.NewNullableString(&name),
			Color: *sw.NewNullableString(&color),
		}

		_, err := c.UpdateTag(tagID, body)
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				d.SetId("")
				var diags diag.Diagnostics
				return diags
			}

			return diag.FromErr(err)
		}
	}

	return resourceTagRead(ctx, d, m)
}

func resourceTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	tagID, convErr := strconv.ParseInt(d.Id(), 10, 64)
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	err := c.DeleteTag(tagID)
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
