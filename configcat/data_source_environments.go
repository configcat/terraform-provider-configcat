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

func dataSourceConfigCatEnvironments() *schema.Resource {
	return &schema.Resource{

		ReadContext: environmentRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			ENVIRONMENT_NAME_FILTER_REGEX: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			ENVIRONMENTS: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ENVIRONMENT_ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ENVIRONMENT_NAME: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ENVIRONMENT_DESCRIPTION: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ENVIRONMENT_COLOR: {
							Type:     schema.TypeString,
							Computed: true,
						},
						ENVIRONMENT_ORDER: {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func environmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(PRODUCT_ID).(string)
	environmentNameFilterRegex := d.Get(ENVIRONMENT_NAME_FILTER_REGEX).(string)

	environments, err := c.GetEnvironments(productID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredEnvironments := []sw.EnvironmentModel{}
	if environmentNameFilterRegex == "" {
		filteredEnvironments = environments
	} else {
		regex := regexp.MustCompile(environmentNameFilterRegex)
		for i := range environments {
			if regex.MatchString(*environments[i].Name.Get()) {
				filteredEnvironments = append(filteredEnvironments, environments[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(ENVIRONMENTS, flattenEnvironmentsData(&filteredEnvironments))

	var diags diag.Diagnostics
	return diags
}

func flattenEnvironmentsData(environments *[]sw.EnvironmentModel) []interface{} {
	if environments != nil {
		elements := make([]interface{}, len(*environments))

		for i, environment := range *environments {
			element := make(map[string]interface{})

			element[ENVIRONMENT_ID] = environment.EnvironmentId
			element[ENVIRONMENT_NAME] = environment.Name.Get()
			element[ENVIRONMENT_DESCRIPTION] = environment.Description.Get()
			element[ENVIRONMENT_COLOR] = environment.Color.Get()
			element[ENVIRONMENT_ORDER] = environment.Order

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
