package configcat

import (
	"context"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigCatSegment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSegmentCreate,
		ReadContext:   resourceSegmentRead,
		UpdateContext: resourceSegmentUpdate,
		DeleteContext: resourceSegmentDelete,
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

			SEGMENT_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			SEGMENT_DESCRIPTION: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SEGMENT_COMPARISON_ATTRIBUTE: {
				Type:     schema.TypeString,
				Required: true,
			},
			SEGMENT_COMPARATOR: {
				Type:     schema.TypeString,
				Required: true,
			},
			SEGMENT_COMPARISON_VALUE: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSegmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	productID := d.Get(PRODUCT_ID).(string)

	comparator, compErr := sw.NewRolloutRuleComparatorFromValue(d.Get(SEGMENT_COMPARATOR).(string))
	if compErr != nil {
		return diag.FromErr(compErr)
	}

	segmentDescription := d.Get(SEGMENT_DESCRIPTION).(string)
	body := sw.CreateSegmentModel{
		Name:                d.Get(SEGMENT_NAME).(string),
		Description:         *sw.NewNullableString(&segmentDescription),
		ComparisonAttribute: d.Get(SEGMENT_COMPARISON_ATTRIBUTE).(string),
		Comparator:          *comparator,
		ComparisonValue:     d.Get(SEGMENT_COMPARISON_VALUE).(string),
	}

	segment, err := c.CreateSegment(productID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*segment.SegmentId)

	return resourceSegmentRead(ctx, d, m)
}

func resourceSegmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	segment, err := c.GetSegment(d.Id())
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	d.Set(PRODUCT_ID, segment.Product.ProductId)
	d.Set(SEGMENT_NAME, segment.Name.Get())
	d.Set(SEGMENT_DESCRIPTION, segment.Description.Get())
	d.Set(SEGMENT_COMPARISON_ATTRIBUTE, segment.ComparisonAttribute.Get())
	d.Set(SEGMENT_COMPARATOR, segment.Comparator)
	d.Set(SEGMENT_COMPARISON_VALUE, segment.ComparisonValue.Get())

	return diags
}

func resourceSegmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.HasChanges(SEGMENT_NAME, SEGMENT_DESCRIPTION, SEGMENT_COMPARISON_ATTRIBUTE, SEGMENT_COMPARATOR, SEGMENT_COMPARISON_VALUE) {

		comparator, compErr := sw.NewRolloutRuleComparatorFromValue(d.Get(SEGMENT_COMPARATOR).(string))
		if compErr != nil {
			return diag.FromErr(compErr)
		}

		segmentName := d.Get(SEGMENT_NAME).(string)
		segmentDescription := d.Get(SEGMENT_DESCRIPTION).(string)
		segmentComparisonAttribute := d.Get(SEGMENT_COMPARISON_ATTRIBUTE).(string)
		segmentComparisonValue := d.Get(SEGMENT_COMPARISON_VALUE).(string)
		body := sw.UpdateSegmentModel{
			Name:                *sw.NewNullableString(&segmentName),
			Description:         *sw.NewNullableString(&segmentDescription),
			ComparisonAttribute: *sw.NewNullableString(&segmentComparisonAttribute),
			Comparator:          comparator,
			ComparisonValue:     *sw.NewNullableString(&segmentComparisonValue),
		}

		_, err := c.UpdateSegment(d.Id(), body)
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				d.SetId("")
				var diags diag.Diagnostics
				return diags
			}

			return diag.FromErr(err)
		}
	}

	return resourceSegmentRead(ctx, d, m)
}

func resourceSegmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*Client)

	err := c.DeleteSegment(d.Id())
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
