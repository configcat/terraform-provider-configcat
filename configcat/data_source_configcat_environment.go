package configcat

import (
	"context"
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatEnvironment() *schema.Resource {
	return &schema.Resource{

		ReadContext: environmentRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			ENVIRONMENT_NAME: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			ENVIRONMENT_ID: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func environmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(PRODUCT_ID).(string)
	environmentName := d.Get(CONFIG_NAME).(string)

	config, err := findEnvironment(c, productID, environmentName)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	updateEnvironmentResourceData(d, config, productID)
	var diags diag.Diagnostics
	return diags
}

func findEnvironment(c *Client, productID, environmentName string) (*sw.EnvironmentModel, error) {

	environments, err := c.GetEnvironments(productID)

	if err != nil {
		return nil, err
	}

	for i := range environments {
		if environments[i].Name == environmentName {
			return &environments[i], nil
		}
	}

	return nil, fmt.Errorf("could not find Environment. product_id: %s name: %s", productID, environmentName)
}

func updateEnvironmentResourceData(d *schema.ResourceData, m *sw.EnvironmentModel, productID string) {
	d.SetId(m.EnvironmentId)
	d.Set(PRODUCT_ID, productID)
	d.Set(ENVIRONMENT_ID, m.EnvironmentId)
	d.Set(ENVIRONMENT_NAME, m.Name)
}
