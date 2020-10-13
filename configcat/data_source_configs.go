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

func dataSourceConfigCatConfigs() *schema.Resource {
	return &schema.Resource{

		ReadContext: configRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			CONFIG_NAME_FILTER_REGEX: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexFunc,
			},

			CONFIGS: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						CONFIG_ID: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						CONFIG_NAME: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func configRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(PRODUCT_ID).(string)
	configNameFilterRegex := d.Get(CONFIG_NAME_FILTER_REGEX).(string)

	configs, err := getConfigs(c, productID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	filteredConfigs := []sw.ConfigModel{}
	if configNameFilterRegex == "" {
		filteredConfigs = configs
	} else {
		regex := regexp.MustCompile(configNameFilterRegex)
		for i := range configs {
			if regex.MatchString(configs[i].Name) {
				filteredConfigs = append(filteredConfigs, configs[i])
			}
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set(CONFIGS, flattenConfigsData(&filteredConfigs))

	var diags diag.Diagnostics
	return diags
}

func getConfigs(c *Client, productID string) ([]sw.ConfigModel, error) {

	configs, err := c.GetConfigs(productID)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func flattenConfigsData(configs *[]sw.ConfigModel) []interface{} {
	if configs != nil {
		elements := make([]interface{}, len(*configs), len(*configs))

		for i, config := range *configs {
			element := make(map[string]interface{})

			element[CONFIG_ID] = config.ConfigId
			element[CONFIG_NAME] = config.Name

			elements[i] = element
		}

		return elements
	}

	return make([]interface{}, 0)
}
