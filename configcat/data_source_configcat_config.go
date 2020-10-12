package configcat

import (
	"context"
	"fmt"

	sw "github.com/configcat/configcat-publicapi-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigCatConfig() *schema.Resource {
	return &schema.Resource{

		ReadContext: configRead,

		Schema: map[string]*schema.Schema{
			PRODUCT_ID: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateGUIDFunc,
			},

			CONFIG_NAME: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func configRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	productID := d.Get(PRODUCT_ID).(string)
	configName := d.Get(CONFIG_NAME).(string)

	config, err := findConfig(c, productID, configName)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(config.ConfigId)
	d.Set(PRODUCT_ID, productID)
	d.Set(CONFIG_NAME, config.Name)

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

func findConfig(c *Client, productID, configName string) (*sw.ConfigModel, error) {

	configs, err := getConfigs(c, productID)
	if err != nil {
		return nil, err
	}

	for i := range configs {
		if configs[i].Name == configName {
			return &configs[i], nil
		}
	}

	return nil, fmt.Errorf("could not find Config. product_id: %s name: %s", productID, configName)
}
