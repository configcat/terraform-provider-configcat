package configcat

import (
	"context"
	"regexp"
	"strconv"
	"time"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatSegments() *schema.Resource {
	return &schema.Resource{

		ReadContext: segmentRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			SEGMENT_NAME_FILTER_REGEX: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			SEGMENTS: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SEGMENT_ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SEGMENT_NAME: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SEGMENT_DESCRIPTION: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func segmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(PRODUCT_ID).(string)
	segmentNameFilterRegex := d.Get(SEGMENT_NAME_FILTER_REGEX).(string)

	segments, err := c.GetSegments(productID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredSegments := []sw.SegmentListModel{}
	if segmentNameFilterRegex == "" {
		filteredSegments = segments
	} else {
		regex := regexp.MustCompile(segmentNameFilterRegex)
		for i := range segments {
			if regex.MatchString(*segments[i].Name.Get()) {
				filteredSegments = append(filteredSegments, segments[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(SEGMENTS, flattenSegmentsData(&filteredSegments))

	var diags diag.Diagnostics
	return diags
}

func flattenSegmentsData(segments *[]sw.SegmentListModel) []interface{} {
	if segments != nil {
		elements := make([]interface{}, len(*segments))

		for i, segment := range *segments {
			element := make(map[string]interface{})

			element[SEGMENT_ID] = segment.SegmentId
			element[SEGMENT_NAME] = segment.Name.Get()
			element[SEGMENT_DESCRIPTION] = segment.Description.Get()

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
