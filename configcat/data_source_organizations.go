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

func dataSourceConfigCatOrganizations() *schema.Resource {
	return &schema.Resource{

		ReadContext: organizationRead,

		Schema: map[string]*schema.Schema{
			ORGANIZATION_NAME_FILTER_REGEX: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			ORGANIZATIONS: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ORGANIZATION_ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ORGANIZATION_NAME: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func organizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	organizationNameFilterRegex := d.Get(ORGANIZATION_NAME_FILTER_REGEX).(string)

	organizations, err := c.GetOrganizations()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredOrganizations := []sw.OrganizationModel{}
	if organizationNameFilterRegex == "" {
		filteredOrganizations = organizations
	} else {
		regex := regexp.MustCompile(organizationNameFilterRegex)
		for i := range organizations {
			if regex.MatchString(*organizations[i].Name.Get()) {
				filteredOrganizations = append(filteredOrganizations, organizations[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(ORGANIZATIONS, flattenOrganizationsData(&filteredOrganizations))

	var diags diag.Diagnostics
	return diags
}

func flattenOrganizationsData(organizations *[]sw.OrganizationModel) []interface{} {
	if organizations != nil {
		elements := make([]interface{}, len(*organizations))

		for i, organization := range *organizations {
			element := make(map[string]interface{})

			element[ORGANIZATION_ID] = organization.OrganizationId
			element[ORGANIZATION_NAME] = organization.Name.Get()

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
