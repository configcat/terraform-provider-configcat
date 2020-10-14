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

func dataSourceConfigCatTags() *schema.Resource {
	return &schema.Resource{

		ReadContext: tagRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			TAG_NAME_FILTER_REGEX: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			TAGS: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TAG_ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						TAG_NAME: {
							Type:     schema.TypeString,
							Computed: true,
						},
						TAG_COLOR: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func tagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(PRODUCT_ID).(string)
	tagNameFilterRegex := d.Get(TAG_NAME_FILTER_REGEX).(string)

	tags, err := c.GetTags(productID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredTags := []sw.TagModel{}
	if tagNameFilterRegex == "" {
		filteredTags = tags
	} else {
		regex := regexp.MustCompile(tagNameFilterRegex)
		for i := range tags {
			if regex.MatchString(tags[i].Name) {
				filteredTags = append(filteredTags, tags[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(TAGS, flattenTagsData(&filteredTags))

	var diags diag.Diagnostics
	return diags
}

func flattenTagsData(tags *[]sw.TagModel) []interface{} {
	if tags != nil {
		elements := make([]interface{}, len(*tags), len(*tags))

		for i, tag := range *tags {
			element := make(map[string]interface{})

			element[TAG_ID] = tag.TagId
			element[TAG_NAME] = tag.Name
			element[TAG_COLOR] = tag.Color

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
